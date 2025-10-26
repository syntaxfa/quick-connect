package http

import (
	"html/template"
	"io"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom html/template renderer for Echo framework
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new TemplateRenderer
func NewTemplateRenderer(path string) *TemplateRenderer {
	return &TemplateRenderer{
		templates: template.Must(template.ParseGlob(path)),
	}
}

// Render returns a template document
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}
