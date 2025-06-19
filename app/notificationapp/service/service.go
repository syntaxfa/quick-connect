package service

import (
	"context"
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	"github.com/syntaxfa/quick-connect/types"
)

type Repository interface {
	Save(ctx context.Context, req SendNotificationRequest) (Notification, error)
	FindNotificationByUserID(ctx context.Context, userID types.ID, paginated paginate.RequestBase, isRead *bool) (ListNotificationResponse, error)
	MarkAsRead(ctx context.Context, notificationID, userID types.ID) error
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
	hub    *Hub
}

func New(cfg Config, vld Validate, cache *cachemanager.CacheManager, repo Repository, logger *slog.Logger, hub *Hub) Service {
	go hub.Run()

	return Service{
		cfg:    cfg,
		vld:    vld,
		cache:  cache,
		repo:   repo,
		logger: logger,
		hub:    hub,
	}
}

func (s Service) JoinClient(ctx context.Context, conn Connection, externalUserID string) {
	client := s.NewClient(ctx, conn, externalUserID)

	s.hub.register <- client
	go client.WritePump()
}
