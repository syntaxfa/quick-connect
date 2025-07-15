package service

import (
	"context"
	"fmt"

	"github.com/oklog/ulid/v2"
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

	req.ID = types.ID(ulid.Make().String())

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
	template.Contents = req.Contents

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

func (s Service) TemplateList(ctx context.Context, req ListTemplateRequest) (ListTemplateResponse, error) {
	const op = "service.template.TemplateList"

	if bErr := req.Paginated.BasicValidation(); bErr != nil {
		return ListTemplateResponse{}, richerror.New(op).WithKind(richerror.KindBadRequest)
	}

	templates, gErr := s.repo.GetTemplates(ctx, req)
	if gErr != nil {
		return ListTemplateResponse{}, errlog.ErrLog(richerror.New(op).WithWrapError(gErr).
			WithKind(richerror.KindUnexpected), s.logger)
	}

	return templates, nil
}

func (s Service) RenderNotificationTemplates(ctx context.Context, channel ChannelType, lang string, notifications ...Notification) ([]NotificationMessage, error) {
	const op = "service.template.RenderNotificationTemplates"

	templateNames := make([]string, 0)
	for _, notification := range notifications {
		templateNames = append(templateNames, notification.TemplateName)
	}

	redisKeys := make([]string, 0)
	for _, name := range templateNames {
		redisKeys = append(redisKeys, "template:"+name)
	}

	// Fetch all of them in one go using MGET
	// The result will be a alice of interfaces
	results, mErr := s.cache.MGet(ctx, redisKeys...)
	if mErr != nil {
		return nil, errlog.ErrLog(richerror.New(op).WithWrapError(mErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	cachedTemplates := make(map[string]Template)
	var missedNames []string

	for i, result := range results {
		originalName := templateNames[i]

		if result == nil {
			missedNames = append(missedNames, originalName)
		} else {
			var ok bool
			cachedTemplates[originalName], ok = result.(Template)
			if !ok {
				return nil, errlog.ErrLog(richerror.New(op).WithMessage("cache template type assertion is not ok").
					WithMeta(map[string]interface{}{"value": result}).WithKind(richerror.KindUnexpected), s.logger)
			}
		}
	}

	var dbTemplates []Template
	if len(missedNames) > 0 {
		var gtErr error
		dbTemplates, gtErr = s.repo.GetTemplatesByNames(ctx, missedNames...)
		if gtErr != nil {
			return nil, errlog.ErrLog(richerror.New(op).WithWrapError(gtErr).WithKind(richerror.KindUnexpected), s.logger)
		}
	}

	finalTemplates := make(map[string]Template)

	// Add templates fetch from DB
	for _, t := range dbTemplates {
		finalTemplates[t.Name] = t

		// Update cache for the next request
		redisKey := "template:" + t.Name
		if sErr := s.cache.Set(ctx, redisKey, t, s.cfg.TemplateCacheExpiration); sErr != nil {
			return nil, errlog.ErrLog(richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected), s.logger)
		}
	}

	// Add templates fetch from cache
	for name, t := range cachedTemplates {
		finalTemplates[name] = t
	}

	notificationMessages := make([]NotificationMessage, 0)
	for _, n := range notifications {
		res, rErr := s.renderSvc.RenderTemplate(finalTemplates[n.TemplateName], channel, lang, n.DynamicTitleData, n.DynamicBodyData)
		if rErr != nil {
			return nil, errlog.ErrLog(richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected).
				WithMessage(fmt.Sprintf("can't render notification %s", n.ID)), s.logger)
		}

		notificationMessages = append(notificationMessages, NotificationMessage{
			ID:        n.ID,
			UserID:    n.UserID,
			Type:      n.Type,
			Data:      n.Data,
			Title:     res.Title,
			Body:      res.Body,
			IsRead:    n.IsRead,
			Timestamp: n.CreatedAt.Unix(),
		})
	}

	return notificationMessages, nil
}
