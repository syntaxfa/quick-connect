package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

func (s Service) AddTemplate(ctx context.Context, req AddTemplateRequest) (AddTemplateResponse, error) {
	const op = "service.template.AddTemplate"

	if vErr := s.vld.ValidateAddTemplateRequest(req); vErr != nil {
		return AddTemplateResponse{}, vErr
	}

	exists, eErr := s.repo.IsExistTemplateName(ctx, req.Name)
	if eErr != nil {
		return AddTemplateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(eErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if exists {
		return AddTemplateResponse{}, richerror.New(op).WithMessage(servermsg.MsgConflictTemplate).
			WithKind(richerror.KindConflict)
	}

	template, cErr := s.repo.CreateTemplate(ctx, req)
	if cErr != nil {
		return AddTemplateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(cErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return template, nil
}

func (s Service) UpdateTemplate(ctx context.Context, req AddTemplateRequest) (AddTemplateResponse, error) {
	const op = "service.template.UpdateTemplate"

	if vErr := s.vld.ValidateAddTemplateRequest(req); vErr != nil {
		return AddTemplateResponse{}, vErr
	}

	exists, eErr := s.repo.IsExistTemplateName(ctx, req.Name)
	if eErr != nil {
		return AddTemplateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(eErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if !exists {
		return AddTemplateResponse{}, richerror.New(op).WithMessage(servermsg.MsgTemplateNotFound).
			WithKind(richerror.KindNotFound)
	}

	template, uErr := s.repo.UpdateTemplate(ctx, req)
	if uErr != nil {
		return AddTemplateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(uErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return template, nil
}
