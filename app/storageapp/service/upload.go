package service

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) Upload(ctx context.Context, req UploadRequest) (File, error) {
	const op = "service.upload.Upload"
	ext := filepath.Ext(req.Filename)

	newID := types.ID(ulid.Make().String())

	key := fmt.Sprintf("uploads/%s%s", newID, ext)

	uploadKey, uErr := s.storage.Upload(ctx, req.File, req.Size, key, req.ContentType, req.IsPublic)
	if uErr != nil {
		return File{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	file := File{
		ID:          newID,
		UploaderID:  req.UploaderID,
		Name:        req.Filename,
		Key:         uploadKey,
		MimeType:    req.ContentType,
		Size:        req.Size,
		Driver:      s.cfg.Driver,
		Bucket:      s.cfg.Bucket,
		IsPublic:    req.IsPublic,
		IsConfirmed: false,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DeletedAt:   nil,
	}

	if sErr := s.repo.Save(ctx, file); sErr != nil {
		return File{}, errlog.ErrContext(ctx, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	return file, nil
}
