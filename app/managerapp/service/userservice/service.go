package userservice

import (
	"context"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/types"
)

type TokenSvc interface {
	GenerateTokenPair(userID types.ID, role types.Role) (*tokenservice.TokenGenerateResponse, error)
}

type Repository interface {
	IsExistUserByUsername(ctx context.Context, username string) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
}

type Service struct {
	tokenSvc TokenSvc
	vld      Validate
	repo     Repository
}

func New(tokenSvc TokenSvc, vld Validate, repo Repository) Service {
	return Service{
		tokenSvc: tokenSvc,
		vld:      vld,
		repo:     repo,
	}
}
