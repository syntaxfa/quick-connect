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

	isInApp := true
	notifications, paginateResp, fErr := s.repo.FindNotificationByUserID(ctx, userID, req.Paginated, req.IsRead, &isInApp)
	if fErr != nil {
		return ListNotificationResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(fErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	userSetting, usErr := s.GetUserSetting(ctx, req.ExternalUserID)
	if usErr != nil {
		return ListNotificationResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(usErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	var accessNotifications []Notification
	for _, notification := range notifications {
		if s.CheckNotificationAccessToSend(notification, userSetting, ChannelTypeInApp) {
			accessNotifications = append(accessNotifications, notification)
		}
	}

	notificationMsgs, rErr := s.RenderNotificationTemplates(ctx, ChannelTypeInApp, userSetting.Lang, accessNotifications...)
	if rErr != nil {
		return ListNotificationResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(rErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return ListNotificationResponse{
		Results:  notificationMsgs,
		Paginate: paginateResp,
	}, nil
}
