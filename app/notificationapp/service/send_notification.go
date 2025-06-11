package service

import (
	"context"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) SendNotification(ctx context.Context, req SendNotificationRequest) (SendNotificationResponse, error) {
	const op = "service.send_notification.SendNotification"

	if vErr := s.vld.ValidateSendNotificationRequest(req); vErr != nil {
		return SendNotificationResponse{}, vErr
	}

	userID, gErr := s.getUserIDFromExternalUserID(ctx, req.ExternalUserID)
	if gErr != nil {
		return SendNotificationResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(gErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}
	req.UserID = userID

	notification, sErr := s.repo.Save(ctx, SendNotificationRequest{
		ID:     types.ID(ulid.Make().String()),
		UserID: req.UserID,
		Type:   req.Type,
		Title:  req.Title,
		Body:   req.Body,
		Data:   req.Data,
	})
	if sErr != nil {
		return SendNotificationResponse{}, errlog.ErrLog(richerror.New(op).
			WithMessage("can't save notification").WithWrapError(sErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	// TODO: check notification type if is critical, send it and return response

	return SendNotificationResponse{Notification: notification}, nil
}
