package service

import (
	"context"
	"fmt"
	"strings"

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

func (s Service) getTemplates(ctx context.Context, names []string) (map[string]Template, error) {
	const op = "service.template.getTemplates"

	if len(names) == 0 {
		return make(map[string]Template), nil
	}

	templatesToProcess := make(map[string]*Template)
	destMap := make(map[string]any, len(names))
	cacheKeys := make([]string, len(names))

	for i, name := range names {
		key := "template:" + name
		cacheKeys[i] = key
		t := &Template{}
		templatesToProcess[name] = t
		destMap[key] = t
	}

	missedKeys, mgErr := s.cache.MGet(ctx, destMap, cacheKeys...)
	if mgErr != nil {
		return nil, richerror.New(op).WithWrapError(mgErr).WithKind(richerror.KindUnexpected)
	}

	if len(missedKeys) > 0 {
		missedNames := make([]string, len(missedKeys))
		for i, key := range missedKeys {
			missedNames[i] = strings.TrimPrefix(key, "template:")
		}

		dbTemplates, dbErr := s.repo.GetTemplatesByNames(ctx, missedNames...)
		if dbErr != nil {
			return nil, richerror.New(op).WithWrapError(dbErr).WithKind(richerror.KindUnexpected)
		}

		for _, dbTemp := range dbTemplates {
			if tempPtr, ok := templatesToProcess[dbTemp.Name]; ok {
				*tempPtr = dbTemp
			}

			cacheKey := "template:" + dbTemp.Name
			if sErr := s.cache.Set(ctx, cacheKey, dbTemp, s.cfg.TemplateCacheExpiration); sErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected), s.logger)
			}
		}
	}

	finalTemplates := make(map[string]Template)
	for name, tempPtr := range templatesToProcess {
		if tempPtr != nil && tempPtr.ID != "" {
			finalTemplates[name] = *tempPtr
		}
	}

	return finalTemplates, nil
}

func (s Service) RenderNotificationTemplates(ctx context.Context, channel ChannelType, lang string, notifications ...Notification) ([]NotificationMessage, error) {
	const op = "service.template.RenderNotificationTemplates"

	templateNames := make([]string, len(notifications))
	for i, notification := range notifications {
		templateNames[i] = notification.TemplateName
	}

	templates, tErr := s.getTemplates(ctx, templateNames)
	if tErr != nil {
		return nil, errlog.ErrLog(richerror.New(op).WithWrapError(tErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	notificationMessages := make([]NotificationMessage, 0)
	for _, n := range notifications {
		res, rErr := s.renderSvc.RenderTemplate(templates[n.TemplateName], channel, lang, n.DynamicTitleData, n.DynamicBodyData)
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
