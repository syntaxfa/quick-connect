package service

import (
	"context"
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	"github.com/syntaxfa/quick-connect/pkg/pubsub"
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
	IsExistTemplateByName(ctx context.Context, name string) (bool, error)
	IsExistTemplateByID(ctx context.Context, id types.ID) (bool, error)
	CreateTemplate(ctx context.Context, req AddTemplateRequest) (Template, error)
	UpdateTemplate(ctx context.Context, id types.ID, req AddTemplateRequest) error
	GetTemplateByName(ctx context.Context, name string) (Template, error)
	GetTemplateByID(ctx context.Context, id types.ID) (Template, error)
}

type Service struct {
	cfg       Config
	vld       Validate
	cache     *cachemanager.CacheManager
	repo      Repository
	logger    *slog.Logger
	hub       *Hub
	publisher pubsub.Publisher
}

func New(cfg Config, vld Validate, cache *cachemanager.CacheManager, repo Repository, logger *slog.Logger, hub *Hub, publisher pubsub.Publisher) Service {
	go hub.Run(context.Background())

	return Service{
		cfg:       cfg,
		vld:       vld,
		cache:     cache,
		repo:      repo,
		logger:    logger,
		hub:       hub,
		publisher: publisher,
	}
}

func (s Service) JoinClient(ctx context.Context, conn Connection, externalUserID string) {
	client := s.NewClient(ctx, conn, externalUserID)

	s.hub.register <- client
	go client.WritePump()
}
