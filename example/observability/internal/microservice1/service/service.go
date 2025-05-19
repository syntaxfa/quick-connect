package service

import (
	"context"
	"fmt"
	"github.com/syntaxfa/quick-connect/adapter/observability/traceotela"
	"github.com/syntaxfa/quick-connect/example/observability/internal/adapter/micro2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type Repository interface {
	GetUserByID(ctx context.Context, userID int)
}

type Service struct {
	repo        Repository
	microClient *micro2.Client
}

func New(repo Repository, micro2Client *micro2.Client) Service {
	return Service{
		repo:        repo,
		microClient: micro2Client,
	}
}

func (s Service) GetUser(ctx context.Context) error {
	fmt.Println("GetUserService")

	trCtx, span := traceotela.Tracer().Start(ctx, "GetUserService")

	s.repo.GetUserByID(trCtx, 2)

	span.AddEvent("user", trace.WithAttributes(attribute.Int("userID", 2)))
	span.End()

	comment, err := s.microClient.GetComment(trCtx, 2)
	if err != nil {
		return err
	}

	fmt.Printf("comment: %+v", comment)

	return nil
}
