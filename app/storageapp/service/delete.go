package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) Delete(ctx context.Context, fileID types.ID) error {
	const op = "service.delete.Delete"

	exists, exErr := s.repo.IsExistByID(ctx, fileID)
	if exErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(exErr).WithKind(richerror.KindUnexpected), s.logger)
	}
	if !exists {
		return richerror.New(op).WithMessage(servermsg.MsgFileNotFound).WithKind(richerror.KindNotFound)
	}

	if dErr := s.repo.DeleteByID(ctx, fileID); dErr != nil {
		return errlog.ErrContext(ctx, richerror.New(op).WithWrapError(dErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return nil
}
