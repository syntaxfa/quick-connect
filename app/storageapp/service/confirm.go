package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) ConfirmFile(ctx context.Context, fileID types.ID) error {
	const op = "service.confirm.ConfirmFile"

	exists, exErr := s.repo.IsExistByID(ctx, fileID)
	if exErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(exErr).WithKind(richerror.KindUnexpected), s.logger)
	}
	if !exists {
		return richerror.New(op).WithMessage(servermsg.MsgFileNotFound).WithKind(richerror.KindNotFound)
	}

	file, gErr := s.repo.GetByID(ctx, fileID)
	if gErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if file.IsConfirmed {
		return richerror.New(op).WithMessage(servermsg.MsgFileAlreadyConfirmed).WithKind(richerror.KindConflict)
	}

	if cErr := s.repo.ConfirmFile(ctx, fileID); cErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return nil
}
