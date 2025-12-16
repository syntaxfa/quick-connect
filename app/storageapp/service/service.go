package service

import (
	"context"
	"io"
)

type Storage interface {
	Upload(ctx context.Context, file io.Reader, size int64, key string, contentType string, isPublic bool) (string, error)
	Delete(ctx context.Context, key string) error
	GetURL(ctx context.Context, key string) (string, error)
	GetPresignedURL(ctx context.Context, key string) (string, error)
	Exists(ctx context.Context, key string) (bool, error)
}

type Service struct{}

func New() Service {
	return Service{}
}
