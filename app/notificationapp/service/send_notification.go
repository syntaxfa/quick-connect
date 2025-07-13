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

func (s Service) SendNotification(ctx context.Context, req SendNotificationRequest) (Notification, error) {
	const op = "service.send_notification.SendNotification"

	if vErr := s.vld.ValidateSendNotificationRequest(req); vErr != nil {
		return Notification{}, vErr
	}

	userID, gErr := s.getUserIDFromExternalUserID(ctx, req.ExternalUserID)
	if gErr != nil {
		return Notification{}, errlog.ErrLog(richerror.New(op).WithWrapError(gErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}
	req.UserID = userID

	req.ID = types.ID(ulid.Make().String())

	for _, channel := range req.ChannelDeliveries {
		if channel.Channel == ChannelTypeInApp {
			req.IsInApp = true
		}
	}

	notification, sErr := s.repo.Save(ctx, req)
	if sErr != nil {
		return Notification{}, errlog.ErrLog(richerror.New(op).
			WithMessage("can't save notification").WithWrapError(sErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	go s.publishNotification(s.cfg.PublishTimeout, notification) //nolint:contextcheck // This function run asynchronously

	return notification, nil
}

func (s Service) publishNotification(ctxTimeout time.Duration, notification Notification) {
	const op = "service.send_notification.publishNotification"

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	// TODO: render template
	notificationMsg := &NotificationMessage{
		NotificationID: notification.ID,
		UserID:         notification.UserID,
		Type:           notification.Type,
		Title:          "test title",
		Body:           "test body",
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
