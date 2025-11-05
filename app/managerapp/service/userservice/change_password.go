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

	oldHashedPass, ohErr := HashPassword(req.OldPassword)
	if ohErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(ohErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if isCorrect, cErr := s.repo.PasswordIsCorrect(ctx, userID, oldHashedPass); cErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), s.logger)
	} else if !isCorrect {
		return errlog.ErrContext(ctx, richerror.New(op).WithMessage(servermsg.MsgPasswordIsNotCorrect).
			WithKind(richerror.KindForbidden), s.logger)
	}

	newHashedPass, nhErr := HashPassword(req.NewPassword)
	if nhErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(nhErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if changeErr := s.repo.ChangePassword(ctx, userID, newHashedPass); changeErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(changeErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return nil
}
