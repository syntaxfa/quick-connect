package userservice

import (
	"context"
	"fmt"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) ChangePassword(ctx context.Context, userID types.ID, req ChangePasswordRequest) error {
	const op = "service.change_password.ChangePassword"

	if vErr := s.vld.ChangePasswordRequest(req); vErr != nil {
		return vErr
	}

	if exist, existErr := s.repo.IsExistUserByID(ctx, userID); existErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(existErr).WithKind(richerror.KindUnexpected), s.logger)
	} else if !exist {
		return errlog.ErrContext(ctx, richerror.New(op).WithMessage(fmt.Sprintf("user id %s does not exists", userID)).
			WithKind(richerror.KindNotFound), s.logger)
	}

	user, guErr := s.repo.GetUserByID(ctx, userID)
	if guErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithKind(richerror.KindUnexpected).WithWrapError(guErr), s.logger)
	}

	if user.Username == DemoUsername {
		return richerror.New(op).WithMessage(servermsg.MsgQuickConnectReservedUsername).WithKind(richerror.KindForbidden)
	}

	userHashedPass, gErr := s.repo.GetUserHashedPassword(ctx, userID)
	if gErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if !VerifyPassword(userHashedPass, req.OldPassword) {
		return errlog.ErrContext(ctx, richerror.New(op).WithMessage(servermsg.MsgPasswordIsNotCorrect).WithKind(richerror.KindForbidden).
			WithMeta(map[string]interface{}{"user_id": userID}), s.logger)
	}

	newHashPass, nhErr := HashPassword(req.NewPassword)
	if nhErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(nhErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if changeErr := s.repo.ChangePassword(ctx, userID, newHashPass); changeErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(changeErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return nil
}
