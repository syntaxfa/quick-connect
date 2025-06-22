package service

import (
	"context"
	"encoding/json"
	"time"

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
		ID:                types.ID(ulid.Make().String()),
		UserID:            req.UserID,
		Type:              req.Type,
		Title:             req.Title,
		Body:              req.Body,
		Data:              req.Data,
		ChannelDeliveries: req.ChannelDeliveries,
	})
	if sErr != nil {
		return SendNotificationResponse{}, errlog.ErrLog(richerror.New(op).
			WithMessage("can't save notification").WithWrapError(sErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	go s.publishNotification(s.cfg.PublishTimeout, notification) //nolint:contextcheck // This function run asynchronously

	return SendNotificationResponse{Notification: notification}, nil
}

func (s Service) publishNotification(ctxTimeout time.Duration, notification Notification) {
	const op = "service.send_notification.publishNotification"

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	notificationMsg := &NotificationMessage{
		NotificationID: notification.ID,
		UserID:         notification.UserID,
		Type:           notification.Type,
		Title:          notification.Title,
		Body:           notification.Body,
		Data:           notification.Data,
		Timestamp:      notification.CreatedAt.Unix(),
	}

	jsonData, mErr := json.Marshal(notificationMsg)
	if mErr != nil {
		errlog.WithoutErr(richerror.New(op).WithMessage("can't marshalling notification message").WithWrapError(mErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if pErr := s.publisher.Publish(ctx, s.cfg.ChannelName, jsonData); pErr != nil {
		errlog.WithoutErr(richerror.New(op).WithMessage("can't publish notification message").WithWrapError(mErr).WithKind(richerror.KindUnexpected), s.logger)
	}
}
