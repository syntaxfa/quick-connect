package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

func (s Service) UpdateUserSetting(ctx context.Context, externalUserID string, req UpdateUserSettingRequest) (UserSetting, error) {
	const op = "service.user_notification_setting.UpdateUserNotificationSetting"

	if vErr := s.vld.ValidateUpdateUserSettingsRequest(req); vErr != nil {
		return UserSetting{}, vErr
	}

	userID, geErr := s.getUserIDFromExternalUserID(ctx, externalUserID)
	if geErr != nil {
		return UserSetting{}, errlog.ErrLog(richerror.New(op).WithWrapError(geErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	exists, eErr := s.repo.IsExistUserSetting(ctx, userID)
	if eErr != nil {
		return UserSetting{}, errlog.ErrLog(richerror.New(op).WithWrapError(eErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if !exists {
		userSetting, cErr := s.repo.CreateUserSetting(ctx, userID, req)
		if cErr != nil {
			return UserSetting{}, errlog.ErrLog(richerror.New(op).WithWrapError(cErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}

		return userSetting, nil
	}

	userSetting, gErr := s.repo.GetUserSetting(ctx, userID)
	if gErr != nil {
		return UserSetting{}, errlog.ErrLog(richerror.New(op).WithWrapError(gErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if uErr := s.repo.UpdateUserSetting(ctx, userID, req); uErr != nil {
		return UserSetting{}, errlog.ErrLog(richerror.New(op).WithWrapError(uErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	userSetting.Lang = req.Lang
	userSetting.IgnoreChannels = req.IgnoreChannels

	return userSetting, nil
}
