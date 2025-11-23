package userservice

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) UserProfile(ctx context.Context, userID types.ID) (UserProfileResponse, error) {
	const op = "service.profile.UserProfile"

	if exist, eErr := s.repo.IsExistUserByID(ctx, userID); eErr != nil {
		return UserProfileResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(eErr).WithKind(richerror.KindUnexpected), s.logger)
	} else if !exist {
		return UserProfileResponse{}, richerror.New(op).WithMessage(servermsg.MsgRecordNotFound).WithKind(richerror.KindNotFound)
	}

	user, gErr := s.repo.GetUserByID(ctx, userID)
	if gErr != nil {
		return UserProfileResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return UserProfileResponse{user}, nil
}

// UserInfo retrieves User info for client & support
// We don't check user exists.
func (s Service) UserInfo(ctx context.Context, userID types.ID) (UserInfoResponse, error) {
	const op = "service.profile.UserInfo"

	userInfo, gErr := s.repo.GetUserInfoByID(ctx, userID)
	if gErr != nil {
		return UserInfoResponse{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return userInfo, nil
}
