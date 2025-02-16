package handler

import (
	"context"

	// "grpc-vertica-crud/internal/database"

	dbs "github.com/prk327/grpc-rest-crudapi/server/database"
	// pb "grpc-vertica-crud/pkg/proto"

	crudv1 "github.com/prk327/grpc-rest-crudapi/proto/crud/v1"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type CRUDService struct {
	crudv1.UnimplementedCrudServiceServer
	db *dbs.Database
}

func NewCRUDService(db *dbs.Database) *CRUDService {
	return &CRUDService{db: db}
}

func (s *CRUDService) CreateItem(ctx context.Context, req *crudv1.CreateItemRequest) (*crudv1.ItemResponse, error) {
	var id string
	err := s.db.Conn.QueryRowContext(ctx,
		"INSERT INTO items (name, description) VALUES ($1, $2) RETURNING id",
		req.Name, req.Description).Scan(&id)

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &crudv1.ItemResponse{
		Id:          id,
		Name:        req.Name,
		Description: req.Description,
	}, nil
}

// Implement other CRUD methods following the same pattern
