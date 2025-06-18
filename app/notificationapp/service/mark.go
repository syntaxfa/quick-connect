package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) MarkNotificationAsRead(ctx context.Context, notificationID types.ID) error {
	const op = "service.mark.MarkNotificationAsRead"

	if mErr := s.repo.MarkAsRead(ctx, notificationID); mErr != nil {
		return errlog.ErrLog(richerror.New(op).WithWrapError(mErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return nil
}

func (s Service) MarkAllNotificationAsRead(ctx context.Context, externalUserID string) error {
	const op = "service.mark.MarkAllNotificationAsRead"

	userID, gErr := s.getUserIDFromExternalUserID(ctx, externalUserID)
	if gErr != nil {
		return errlog.ErrLog(richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if mErr := s.repo.MarkAllAsReadByUserID(ctx, userID); mErr != nil {
		return errlog.ErrLog(richerror.New(op).WithWrapError(mErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return nil
}
