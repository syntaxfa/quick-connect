package service

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, file io.Reader, size int64, key, contentType string, isPublic bool) (string, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) (string, error)
	GetPresignedURL(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type Repository interface {
}

type Service struct {
	storage Storage
	repo    Repository
}

func New(storage Storage, repo Repository) Service {
	return Service{
		storage: storage,
		repo:    repo,
	}
}
