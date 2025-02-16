package handler

import (
	"context"

	"google.golang.org/grpc/credentials/insecure"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	crudv1 "github.com/prk327/grpc-rest-crudapi/proto/crud/v1"
	"google.golang.org/grpc"
)

func RegisterHTTPHandlers(ctx context.Context, mux *runtime.ServeMux, grpcAddr string) error {
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := crudv1.RegisterCrudServiceHandlerFromEndpoint(ctx, mux, grpcAddr, opts)
	if err != nil {
		return err
	}

	return nil
}
