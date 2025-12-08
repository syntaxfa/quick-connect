package userservice

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) UserUpdateFromSuperuser(ctx context.Context, userID types.ID,
	req UserUpdateFromSuperuserRequest) (UserUpdateResponse, error) {
	const op = "service.user_update.UserUpdateFromSuperuser"

	if vErr := s.vld.UserUpdateFromSuperuserRequest(req); vErr != nil {
		return UserUpdateResponse{}, vErr
	}

	exist, existErr := s.repo.IsExistUserByID(ctx, userID)
	if existErr != nil {
		return UserUpdateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(existErr).WithKind(richerror.KindUnexpected), s.logger)
	} else if !exist {
		return UserUpdateResponse{}, richerror.New(op).WithMessage(servermsg.MsgRecordNotFound).WithKind(richerror.KindNotFound)
	}

	user, guErr := s.repo.GetUserByID(ctx, userID)
	if guErr != nil {
		return UserUpdateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(guErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if req.Username != user.Username {
		userExist, usExistErr := s.repo.IsExistUserByUsername(ctx, req.Username)
		if usExistErr != nil {
			return UserUpdateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(usExistErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}

		if userExist {
			return UserUpdateResponse{}, errlog.ErrLog(richerror.New(op).WithKind(richerror.KindConflict).
				WithMessage(servermsg.MsgConflictUsername), s.logger)
		}
	}

	upErr := s.repo.UpdateUser(ctx, userID, req)
	if upErr != nil {
		return UserUpdateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(upErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	updatedUser, guErr := s.repo.GetUserByID(ctx, userID)
	if guErr != nil {
		return UserUpdateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(guErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return UserUpdateResponse{
		User: updatedUser,
	}, nil
}

func (s Service) UserUpdateFromOwn(ctx context.Context, userID types.ID, req UserUpdateFromOwnRequest) (UserUpdateResponse, error) {
	const op = "service.user_update.UpdateUserFromOwn"

	if req.Username == DemoUsername {
		return UserUpdateResponse{}, richerror.New(op).WithMessage(servermsg.MsgQuickConnectReservedUsername).
			WithKind(richerror.KindForbidden)
	}

	user, guErr := s.repo.GetUserByID(ctx, userID)
	if guErr != nil {
		return UserUpdateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(guErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if user.Username == DemoUsername {
		return UserUpdateResponse{}, richerror.New(op).WithMessage(servermsg.MsgQuickConnectReservedUsername).
			WithKind(richerror.KindForbidden)
	}

	return s.UserUpdateFromSuperuser(ctx, userID, UserUpdateFromSuperuserRequest{
		Username:    req.Username,
		Fullname:    req.Fullname,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
		Roles:       user.Roles,
	})
}
