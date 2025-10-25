package userservice

import (
	"context"
	"fmt"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) UserDelete(ctx context.Context, userID types.ID) error {
	const op = "service.user_delete.UserDelete"

	if exist, eErr := s.repo.IsExistUserByID(ctx, userID); eErr != nil {
		return errlog.ErrLog(richerror.New(op).WithKind(richerror.KindUnexpected).WithWrapError(eErr), s.logger)
	} else if !exist {
		return richerror.New(op).WithKind(richerror.KindNotFound).WithMessage(servermsg.MsgRecordNotFound)
	}

	if dErr := s.repo.DeleteUser(ctx, userID); dErr != nil {
		return errlog.ErrLog(richerror.New(op).WithWrapError(dErr).WithKind(richerror.KindUnexpected).WithMessage(fmt.Sprintf("can't delete user id: %s", userID)), s.logger)
	}

	return nil
}
