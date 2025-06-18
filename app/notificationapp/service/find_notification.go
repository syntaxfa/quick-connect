package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

func (s Service) FindNotificationByUserID(ctx context.Context, req ListNotificationRequest) (ListNotificationResponse, error) {
	const op = "service.find_notification.FindNotificationByUserID"

	if vErr := s.vld.ValidateListNotificationRequest(req); vErr != nil {
		return ListNotificationResponse{}, vErr
	}

	if bErr := req.Paginated.BasicValidation(); bErr != nil {
		return ListNotificationResponse{}, richerror.New(op).WithKind(richerror.KindBadRequest)
	}

	userID, gErr := s.getUserIDFromExternalUserID(ctx, req.ExternalUserID)
	if gErr != nil {
		return ListNotificationResponse{}, errlog.ErrLog(gErr, s.logger)
	}

	notifications, fErr := s.repo.FindNotificationByUserID(ctx, userID, req.Paginated, req.IsRead)
	if fErr != nil {
		return ListNotificationResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(fErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return notifications, nil
}
