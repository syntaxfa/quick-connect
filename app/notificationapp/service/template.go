package service

import (
	"context"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

func (s Service) AddTemplate(ctx context.Context, req AddTemplateRequest) (Template, error) {
	const op = "service.template.AddTemplate"

	if vErr := s.vld.ValidateAddTemplateRequest(req); vErr != nil {
		return Template{}, vErr
	}

	exists, eErr := s.repo.IsExistTemplateByName(ctx, req.Name)
	if eErr != nil {
		return Template{}, errlog.ErrLog(richerror.New(op).WithWrapError(eErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if exists {
		return Template{}, richerror.New(op).WithMessage(servermsg.MsgConflictTemplate).
			WithKind(richerror.KindConflict)
	}

	template, cErr := s.repo.CreateTemplate(ctx, req)
	if cErr != nil {
		return Template{}, errlog.ErrLog(richerror.New(op).WithWrapError(cErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return template, nil
}

func (s Service) UpdateTemplate(ctx context.Context, templateID types.ID, req AddTemplateRequest) (Template, error) {
	const op = "service.template.UpdateTemplate"

	if vErr := s.vld.ValidateAddTemplateRequest(req); vErr != nil {
		return Template{}, vErr
	}

	exists, eErr := s.repo.IsExistTemplateByID(ctx, templateID)
	if eErr != nil {
		return Template{}, errlog.ErrLog(richerror.New(op).WithWrapError(eErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}
	if !exists {
		return Template{}, richerror.New(op).WithMessage(servermsg.MsgTemplateNotFound).WithKind(richerror.KindNotFound)
	}

	template, gIDErr := s.repo.GetTemplateByID(ctx, templateID)
	if gIDErr != nil {
		return Template{}, errlog.ErrLog(richerror.New(op).WithWrapError(gIDErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if template.Name != req.Name {
		exists, eErr = s.repo.IsExistTemplateByName(ctx, req.Name)
		if eErr != nil {
			return Template{}, errlog.ErrLog(richerror.New(op).WithWrapError(eErr).
				WithKind(richerror.KindUnexpected), s.logger)
		}

		if exists {
			return Template{}, richerror.New(op).WithMessage(servermsg.MsgConflictTemplate).
				WithKind(richerror.KindConflict)
		}
	}

	uErr := s.repo.UpdateTemplate(ctx, template.ID, req)
	if uErr != nil {
		return Template{}, errlog.ErrLog(richerror.New(op).WithWrapError(uErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	template.Name = req.Name
	template.Bodies = req.Bodies

	return template, nil
}

func (s Service) GetTemplate(ctx context.Context, templateID types.ID) (Template, error) {
	const op = "service.template.GetTemplate"

	exists, eErr := s.repo.IsExistTemplateByID(ctx, templateID)
	if eErr != nil {
		return Template{}, errlog.ErrLog(richerror.New(op).WithWrapError(eErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	if !exists {
		return Template{}, richerror.New(op).WithMessage(servermsg.MsgTemplateNotFound).WithKind(richerror.KindNotFound)
	}

	template, gErr := s.repo.GetTemplateByID(ctx, templateID)
	if gErr != nil {
		return Template{}, errlog.ErrLog(richerror.New(op).WithWrapError(gErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return template, nil
}
