package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) GetFileInfo(ctx context.Context, fileID types.ID) (File, error) {
	const op = "service.get_file.GetFileInfo"

	exists, exErr := s.repo.IsExistByID(ctx, fileID)
	if exErr != nil {
		return File{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(exErr).WithKind(richerror.KindUnexpected), s.logger)
	}
	if !exists {
		return File{}, richerror.New(op).WithMessage(servermsg.MsgFileNotFound).WithKind(richerror.KindNotFound)
	}

	file, gErr := s.repo.GetByID(ctx, fileID)
	if gErr != nil {
		return File{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return file, nil
}
