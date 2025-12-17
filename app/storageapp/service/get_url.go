package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) GetPublicLink(ctx context.Context, fileID types.ID) (string, error) {
	const op = "service.get_url.GetPublicLink"

	exists, exErr := s.repo.IsExistByID(ctx, fileID)
	if exErr != nil {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(exErr).WithKind(richerror.KindUnexpected), s.logger)
	}
	if !exists {
		return "", richerror.New(op).WithMessage(servermsg.MsgFileNotFound).WithKind(richerror.KindNotFound)
	}

	file, gfErr := s.repo.GetByID(ctx, fileID)
	if gfErr != nil {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gfErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if !file.IsPublic {
		return "", richerror.New(op).WithMessage(servermsg.MsgFileInNotPublic)
	}

	url, gErr := s.storage.GetURL(ctx, file.Key)
	if gErr != nil {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return url, nil
}

func (s Service) GetLink(ctx context.Context, fileID types.ID) (string, error) {
	const op = "service.get_url.GetLink"

	exists, exErr := s.repo.IsExistByID(ctx, fileID)
	if exErr != nil {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(exErr).WithKind(richerror.KindUnexpected), s.logger)
	}
	if !exists {
		return "", richerror.New(op).WithMessage(servermsg.MsgFileNotFound).WithKind(richerror.KindNotFound)
	}

	file, gfErr := s.repo.GetByID(ctx, fileID)
	if gfErr != nil {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gfErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	if file.IsPublic {
		url, gErr := s.storage.GetURL(ctx, file.Key)
		if gErr != nil {
			return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
		}

		return url, nil
	}

	url, gErr := s.storage.GetPresignedURL(ctx, file.Key)
	if gErr != nil {
		return "", errlog.ErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return url, nil
}
