package service

import (
	"context"
	"io"
	"log/slog"

	"github.com/syntaxfa/quick-connect/types"
)

type Storage interface {
	Upload(ctx context.Context, file io.Reader, size int64, key, contentType string, isPublic bool) (string, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) (string, error)
	GetPresignedURL(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type Repository interface {
	Save(ctx context.Context, file File) error
	IsExistByID(ctx context.Context, fileID types.ID) (bool, error)
	GetByID(ctx context.Context, fileID types.ID) (File, error)
	DeleteByID(ctx context.Context, fileID types.ID) error
}

type Service struct {
	cfg     Config
	storage Storage
	repo    Repository
	logger  *slog.Logger
}

func New(cfg Config, storage Storage, repo Repository, logger *slog.Logger) Service {
	return Service{
		cfg:     cfg,
		storage: storage,
		repo:    repo,
		logger:  logger,
	}
}
