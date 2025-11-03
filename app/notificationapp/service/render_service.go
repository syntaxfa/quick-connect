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

func (r *RenderService) RenderTemplate(template Template, channel ChannelType, lang string, titleData,
	bodyData map[string]string) (RenderTemplate, error) {
	const op = "service.template_render.RenderTemplate"

	if !IsValidChannelType(channel) {
		return RenderTemplate{}, richerror.New(op).WithMessage(fmt.Sprintf("channel %s is not valid", channel)).
			WithKind(richerror.KindUnexpected)
	}

	content := r.findContentByChannel(template, channel)
	if content == nil {
		return RenderTemplate{}, richerror.New(op).WithKind(richerror.KindUnexpected).
			WithMessage(fmt.Sprintf("%s channel is not in %s template", channel, template.Name))
	}

	tempLang := r.selectLanguage(content.Bodies, lang)
	contentBody := r.findContentBody(content.Bodies, tempLang)
	templateName := fmt.Sprintf("%s:%s:%s", template.Name, channel, tempLang)
	templateType := r.getTemplateType(channel)

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

func (r *RenderService) findContentByChannel(template Template, channel ChannelType) *TemplateContent {
	for _, c := range template.Contents {
		if c.Channel == channel {
			return &c
		}
	}
	return nil
}

func (r *RenderService) selectLanguage(bodies []ContentBody, requestedLang string) string {
	if len(bodies) == 0 {
		return r.defaultLang
	}

	tempLang := bodies[0].Lang

	for _, b := range bodies {
		if b.Lang == r.defaultLang {
			tempLang = r.defaultLang
			break
		}
	}

	for _, b := range bodies {
		if b.Lang == requestedLang {
			return requestedLang
		}
	}

	return tempLang
}

func (r *RenderService) findContentBody(bodies []ContentBody, lang string) *ContentBody {
	for _, b := range bodies {
		if b.Lang == lang {
			return &b
		}
	}
	return nil
}

func (r *RenderService) getTemplateType(channel ChannelType) TemplateType {
	if channel == ChannelTypeEmail {
		return TemplateTypeHTML
	}
	return TemplateTypeText
}

func (r *RenderService) renderTemplate(templateName, template string, templateType TemplateType, data map[string]string) (string, error) {
	const op = "service.template_render.RenderTemplate"

	var executor Executor
	var err error

	if templateType == TemplateTypeText {
		executor, err = r.getOrCreateTextTemplate(templateName, template)
	} else {
		executor, err = r.getOrCreateHTMLTemplate(templateName, template)
	}

	if err != nil {
		return "", err
	}

	var buf bytes.Buffer
	if eErr := executor.Execute(&buf, data); eErr != nil {
		return "", richerror.New(op).WithMessage("can't execute template").WithWrapError(eErr).
			WithKind(richerror.KindUnexpected)
	}

	return buf.String(), nil
}

func (r *RenderService) getOrCreateTextTemplate(templateName, template string) (Executor, error) {
	const op = "service.template_render.getOrCreateTextTemplate"

	r.mu.RLock()
	temp, exists := r.textTemps[templateName]
	r.mu.RUnlock()

	if exists {
		return temp, nil
	}

	newTemp, pErr := textTemp.New(templateName).Parse(template)
	if pErr != nil {
		return nil, richerror.New(op).WithMessage("can't parse text template").WithWrapError(pErr).
			WithKind(richerror.KindUnexpected)
	}

	r.mu.Lock()
	r.textTemps[templateName] = newTemp
	r.mu.Unlock()

	return newTemp, nil
}

func (r *RenderService) getOrCreateHTMLTemplate(templateName, template string) (Executor, error) {
	const op = "service.template_render.getOrCreateHTMLTemplate"

	r.mu.RLock()
	temp, exists := r.htmlTemps[templateName]
	r.mu.RUnlock()

	if exists {
		return temp, nil
	}

	newTemp, pErr := htmlTemp.New(templateName).Parse(template)
	if pErr != nil {
		return nil, richerror.New(op).WithMessage(fmt.Sprintf("can't parse html template, %s", pErr)).WithWrapError(pErr).
			WithKind(richerror.KindUnexpected)
	}

	r.mu.Lock()
	r.htmlTemps[templateName] = newTemp
	r.mu.Unlock()

	return newTemp, nil
}
