package service

import (
	"context"
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
	"github.com/syntaxfa/quick-connect/types"
)

type Repository interface {
	Save(ctx context.Context, req SendNotificationRequest) (Notification, error)
	FindByUserID(ctx context.Context, userID types.ID) ([]Notification, error)
	MarkAsRead(ctx context.Context, notificationID types.ID) error
	MarkAllAsReadByUserID(ctx context.Context, userID types.ID) error
	IsExistUserIDFromExternalUserID(ctx context.Context, externalUserID string) (bool, error)
	GetUserIDFromExternalUserID(ctx context.Context, externalUserID string) (types.ID, error)
	CreateUserIDFromExternalUserID(ctx context.Context, externalUserID string, userID types.ID) error
}

type Service struct {
	cfg    Config
	vld    Validate
	cache  *cachemanager.CacheManager
	repo   Repository
	logger *slog.Logger
}

func New(cfg Config, vld Validate, cache *cachemanager.CacheManager, repo Repository, logger *slog.Logger) Service {
	return Service{
		cfg:    cfg,
		vld:    vld,
		cache:  cache,
		repo:   repo,
		logger: logger,
	}
}
