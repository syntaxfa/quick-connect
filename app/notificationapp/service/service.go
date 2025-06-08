package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/types"
)

type Repository interface {
	Save(ctx context.Context, notification Notification) error
	FindByUserID(ctx context.Context, userID types.ID) ([]Notification, error)
	MarkAsRead(ctx context.Context, notificationID types.ID) error
	MarkAllAsReadByUserID(ctx context.Context, userID types.ID) error
}
