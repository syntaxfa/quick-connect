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
	mu        sync.RWMutex
	textTemps map[string]*textTemp.Template
	htmlTemps map[string]*htmlTemp.Template
}

func NewRenderService() *RenderService {
	return &RenderService{
		textTemps: make(map[string]*textTemp.Template),
		htmlTemps: make(map[string]*htmlTemp.Template),
	}
}

const (
	TemplateTypeText TemplateType = "text"
	TemplateTypeHTML TemplateType = "html"
)

func (r *RenderService) RenderTemplate(templateName, template string, templateType TemplateType, data map[string]string) (string, error) {
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
