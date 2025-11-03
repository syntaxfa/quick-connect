package service

import (
	"context"
	"github.com/syntaxfa/quick-connect/adapter/observability/traceotela"
)

type Service struct {
	repo Repository
}

type Repository interface {
	GetCommentByID(ctx context.Context, commentID uint64) (GetCommentResponse, error)
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) GetCommentByID(ctx context.Context, commentID uint64) (GetCommentResponse, error) {
	cCtx, span := traceotela.Tracer().Start(ctx, "Get comment By ID Service")
	defer span.End()

	return s.repo.GetCommentByID(cCtx, commentID)
}
