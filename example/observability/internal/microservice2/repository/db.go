package repository

import (
	"context"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice2/service"
)

type Database struct {
	db *postgres.Database
}

func New(db *postgres.Database) Database {
	return Database{db: db}
}

func (d Database) GetCommentByID(ctx context.Context, commentID uint64) (service.GetCommentResponse, error) {
	d.db.Conn().Ping(ctx)

	return service.GetCommentResponse{
		ID:   commentID,
		Body: "Hello",
	}, nil
}
