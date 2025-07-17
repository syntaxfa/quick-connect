package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) SendNotification(ctx context.Context, req SendNotificationRequest) (Notification, error) {
	const op = "service.send_notification.SendNotification"

	if vErr := s.vld.ValidateSendNotificationRequest(req); vErr != nil {
		return Notification{}, vErr
	}

	exists, eErr := s.repo.IsExistTemplateByName(ctx, req.TemplateName)
	if eErr != nil {
		return Notification{}, errlog.ErrLog(richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if !exists {
		return Notification{}, richerror.New(op).WithMessage(servermsg.MsgTemplateNotFound).WithKind(richerror.KindNotFound)
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

	for _, ch := range notification.ChannelDeliveries {
		if ch.Channel == ChannelTypeInApp {
			go s.publishNotification(s.cfg.PublishTimeout, notification) //nolint:contextcheck // This function run asynchronously
		}
	}

	return notification, nil
}

func (s Service) publishNotification(ctxTimeout time.Duration, notification Notification) {
	const op = "service.send_notification.publishNotification"

	ctx, cancel := context.WithTimeout(context.Background(), ctxTimeout)
	defer cancel()

	userSetting, usErr := s.GetUserSetting(ctx, string(notification.UserID))
	if usErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(usErr).WithMessage(fmt.Sprintf("can't get user setting for user id: %s", notification.UserID)), s.logger)
	}
	fmt.Println(userSetting)

	if !CheckNotificationAccessToSend(notification, userSetting, ChannelTypeInApp) {
		return
	}
	notificationMsgs, rErr := s.RenderNotificationTemplates(ctx, ChannelTypeInApp, userSetting.Lang, notification)
	if rErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected).
			WithMessage("can't render notification"), s.logger)

		return
	}

	jsonData, mErr := json.Marshal(notificationMsgs[0])
	if mErr != nil {
		errlog.WithoutErr(richerror.New(op).WithMessage("can't marshalling notification message").WithWrapError(mErr).WithKind(richerror.KindUnexpected), s.logger)

		return
	}

	if pErr := s.publisher.Publish(ctx, s.cfg.ChannelName, jsonData); pErr != nil {
		errlog.WithoutErr(richerror.New(op).WithMessage("can't publish notification message").WithWrapError(mErr).WithKind(richerror.KindUnexpected), s.logger)

		return
	}
}

// CheckNotificationAccessToSend if notification type is critical, notification send and doesn't check user ignore channels.
func CheckNotificationAccessToSend(notification Notification, userSetting UserSetting, channel ChannelType) bool {
	if notification.Type == NotificationTypeCritical {
		return true
	}

	for _, ignore := range userSetting.IgnoreChannels {
		if ignore.Channel == channel {
			for _, notificationType := range ignore.NotificationTypes {
				if notificationType == notification.Type {
					return false
				}
			}
		}
	}

	return true
}
