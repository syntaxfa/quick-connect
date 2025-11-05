package userservice

import (
	"context"
	"log/slog"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	"github.com/syntaxfa/quick-connect/types"
)

type TokenSvc interface {
	GenerateTokenPair(userID types.ID, roles []types.Role) (*tokenservice.TokenGenerateResponse, error)
}

type Repository interface {
	IsExistUserByUsername(ctx context.Context, username string) (bool, error)
	IsExistUserByID(ctx context.Context, userID types.ID) (bool, error)
	GetUserByUsername(ctx context.Context, username string) (User, error)
	GetUserByID(ctx context.Context, userID types.ID) (User, error)
	CreateUser(ctx context.Context, req UserCreateRequest) (User, error)
	GetUserList(ctx context.Context, paginated paginate.RequestBase, username string) ([]User, paginate.ResponseBase, error)
	DeleteUser(ctx context.Context, userID types.ID) error
	UpdateUser(ctx context.Context, userID types.ID, req UserUpdateFromSuperuserRequest) error
	PasswordIsCorrect(ctx context.Context, userID types.ID, hashedPassword string) (bool, error)
	ChangePassword(ctx context.Context, userID types.ID, hashedPassword string) error
}

type ExternalUserRepository interface {
	IsExistUserIDFromExternalUserID(ctx context.Context, externalUserID string) (bool, error)
	GetUserIDFromExternalUserID(ctx context.Context, externalUserID string) (types.ID, error)
	CreateUserIDFromExternalUserID(ctx context.Context, externalUserID string, userID types.ID) error
}

type Service struct {
	cfg              Config
	tokenSvc         TokenSvc
	vld              Validate
	repo             Repository
	externalUserRepo ExternalUserRepository
	logger           *slog.Logger
	cache            *cachemanager.CacheManager
}

func New(cfg Config, tokenSvc TokenSvc, vld Validate, repo Repository, externalUserRepo ExternalUserRepository, logger *slog.Logger,
	cache *cachemanager.CacheManager) Service {
	return Service{
		cfg:              cfg,
		tokenSvc:         tokenSvc,
		vld:              vld,
		repo:             repo,
		externalUserRepo: externalUserRepo,
		logger:           logger,
		cache:            cache,
	}
}
