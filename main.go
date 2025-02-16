package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/prk327/grpc-rest-crudapi/insecure"
	crudv1 "github.com/prk327/grpc-rest-crudapi/proto/crud/v1"
	usersv1 "github.com/prk327/grpc-rest-crudapi/proto/users/v1"
	"github.com/prk327/grpc-rest-crudapi/server"
	"github.com/prk327/grpc-rest-crudapi/server/config"
	"github.com/prk327/grpc-rest-crudapi/server/database"
	"github.com/prk327/grpc-rest-crudapi/server/handler"
	"github.com/rs/cors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/grpclog"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	// Adds gRPC internal logs. This is quite verbose, so adjust as desired!
	log := grpclog.NewLoggerV2(os.Stdout, ioutil.Discard, ioutil.Discard)
	grpclog.SetLoggerV2(log)

	// Load configuration
	dbConfig := config.LoadDatabaseConfig()
	serverConfig := config.LoadServerConfig()

	// Initialize database
	db, err := database.New(dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	if err := db.ValidateSchema(ctx, dbConfig.Schema); err != nil {
		log.Fatalf("Schema validation error: %v", err)
	}

	// Run migrations
	if err := db.Migrate(ctx); err != nil {
		log.Fatalf("Database migration failed: %v", err)
	}

	// Create services
	crudService := handler.NewCRUDService(db)

	// Create server group
	g, ctx := errgroup.WithContext(ctx)

	// Start gRPC server
	g.Go(func() error {
		lis, err := net.Listen("tcp", ":"+serverConfig.GRPCPort)
		if err != nil {
			return fmt.Errorf("failed to listen: %w", err)
		}

		s := grpc.NewServer(
			// TODO: Replace with your own certificate!
			grpc.Creds(credentials.NewServerTLSFromCert(&insecure.Cert)),
		)

		usersv1.RegisterUserServiceServer(s, server.New())
		crudv1.RegisterCrudServiceServer(s, crudService)

		log.Info("gRPC server listening on :%s", serverConfig.GRPCPort)
		return s.Serve(lis)
	})

	// Start HTTP gateway
	g.Go(func() error {
		mux := runtime.NewServeMux()
		if err := handler.RegisterHTTPHandlers(ctx, mux, ":"+serverConfig.GRPCPort); err != nil {
			return err
		}

		handler := cors.Default().Handler(mux)
		httpServer := &http.Server{
			Addr:    ":" + serverConfig.HTTPPort,
			Handler: handler,
		}

		log.Info("HTTP gateway listening on :%s", serverConfig.HTTPPort)
		return httpServer.ListenAndServe()
	})

	// Wait for servers to exit
	if err := g.Wait(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
