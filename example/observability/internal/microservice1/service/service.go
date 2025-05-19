package service

import (
	"context"
	"fmt"
	"github.com/syntaxfa/quick-connect/adapter/observability/traceotela"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Repository interface {
	GetUserByID(ctx context.Context, userID int)
}

type Service struct {
	repo Repository
}

func New(repo Repository) Service {
	return Service{
		repo: repo,
	}
}

func (s Service) GetUser(ctx context.Context) error {
	fmt.Println("GetUserService")

	trCtx, span := traceotela.Tracer().Start(ctx, "GetUserService")

	s.repo.GetUserByID(trCtx, 2)

	span.AddEvent("user", trace.WithAttributes(attribute.Int("userID", 2)))
	span.End()

	return nil
}
