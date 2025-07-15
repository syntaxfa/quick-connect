package service

import (
	"bytes"
	"fmt"
	htmlTemp "html/template"
	"io"
	"sync"
	textTemp "text/template"

	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

type TemplateType string

// Executor an interface for execute templates.
type Executor interface {
	Execute(wr io.Writer, data any) error
}

type RenderService struct {
	mu          sync.RWMutex
	textTemps   map[string]*textTemp.Template
	htmlTemps   map[string]*htmlTemp.Template
	defaultLang string
}

func NewRenderService(defaultLang string) *RenderService {
	return &RenderService{
		textTemps:   make(map[string]*textTemp.Template),
		htmlTemps:   make(map[string]*htmlTemp.Template),
		defaultLang: defaultLang,
	}
}

const (
	TemplateTypeText TemplateType = "text"
	TemplateTypeHTML TemplateType = "html"
)

type RenderTemplate struct {
	Name  string `json:"name"`
	Lang  string `json:"lang"`
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (r *RenderService) RenderTemplate(template Template, channel ChannelType, lang string, titleData, bodyData map[string]string) (RenderTemplate, error) {
	const op = "service.template_render.RenderTemplate"

	if !IsValidChannelType(channel) {
		return RenderTemplate{}, richerror.New(op).WithMessage(fmt.Sprintf("channel %s is not valid", channel)).
			WithKind(richerror.KindUnexpected)
	}
	var tempLang string
	var templateName string
	var templateType TemplateType
	var content *TemplateContent
	var contentBody *ContentBody

	for _, c := range template.Contents {
		if c.Channel == channel {
			content = &c
		}
	}
	if content == nil {
		return RenderTemplate{}, richerror.New(op).WithKind(richerror.KindUnexpected).
			WithMessage(fmt.Sprintf("%s channel is not in %s template", channel, template.Name))
	}

	tempLang = content.Bodies[0].Lang

	for _, b := range content.Bodies {
		if b.Lang == r.defaultLang {
			tempLang = r.defaultLang
		}
	}

	for _, b := range content.Bodies {
		if b.Lang == lang {
			tempLang = lang
		}
	}

	for _, b := range content.Bodies {
		if b.Lang == tempLang {
			contentBody = &b
		}
	}

	templateName = fmt.Sprintf("%s:%s:%s", template.Name, channel, tempLang)

	switch channel {
	case ChannelTypeEmail:
		templateType = TemplateTypeHTML
	default:
		templateType = TemplateTypeText
	}

	title, rtErr := r.renderTemplate(fmt.Sprintf("%s:%s", templateName, "title"), contentBody.Title, templateType, titleData)
	if rtErr != nil {
		return RenderTemplate{}, richerror.New(op).WithMessage("can't render template content title").
			WithKind(richerror.KindUnexpected).WithMeta(map[string]interface{}{"contentTitle": contentBody.Title, "titleData": titleData})
	}

	body, rbErr := r.renderTemplate(fmt.Sprintf("%s:%s", templateName, "body"), contentBody.Body, templateType, bodyData)
	if rbErr != nil {
		return RenderTemplate{}, richerror.New(op).WithMessage("can't render template content body").
			WithKind(richerror.KindUnexpected).WithMeta(map[string]interface{}{"contentBody": contentBody.Body, "bodyData": bodyData})
	}

	return RenderTemplate{
		Name:  template.Name,
		Lang:  tempLang,
		Title: title,
		Body:  body,
	}, nil
}

func (r *RenderService) renderTemplate(templateName, template string, templateType TemplateType, data map[string]string) (string, error) {
	const op = "service.template_render.RenderTemplate"

	var buf bytes.Buffer
	var executor Executor

	if templateType == TemplateTypeText {
		r.mu.RLock()
		temp, exists := r.textTemps[templateName]
		r.mu.RUnlock()

		if !exists {
			var pErr error
			temp, pErr = textTemp.New(templateName).Parse(template)
			if pErr != nil {
				return "", richerror.New(op).WithMessage("can't parse text template").WithWrapError(pErr).
					WithKind(richerror.KindUnexpected)
			}

			r.mu.Lock()
			r.textTemps[templateName] = temp
			r.mu.Unlock()
		}
		executor = temp

	} else {
		r.mu.RLock()
		temp, exists := r.htmlTemps[templateName]
		r.mu.RUnlock()

		if !exists {
			var pErr error
			temp, pErr = htmlTemp.New(templateName).Parse(template)
			if pErr != nil {
				return "", richerror.New(op).WithMessage(fmt.Sprintf("can't parse html template, %s", pErr)).WithWrapError(pErr).
					WithKind(richerror.KindUnexpected)
			}

			r.mu.Lock()
			r.htmlTemps[templateName] = temp
			r.mu.Unlock()
		}
		executor = temp
	}

	if eErr := executor.Execute(&buf, data); eErr != nil {
		return "", richerror.New(op).WithMessage("can't execute template").WithWrapError(eErr).
			WithKind(richerror.KindUnexpected)
	}

	return buf.String(), nil
}
