package http

import (
	"fmt"
	"html/template"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// TemplateRenderer is a custom html/template renderer for Echo framework.
type TemplateRenderer struct {
	templates *template.Template
}

// NewTemplateRenderer creates a new TemplateRenderer.
func NewTemplateRenderer(basePath string) *TemplateRenderer {
	// Define custom functions
	funcMap := template.FuncMap{
		"initials": getInitials,
		"hasRole":  HasRole,
	}

	// Check if base path exists
	if _, err := os.Stat(basePath); os.IsNotExist(err) {
		panic(fmt.Sprintf("template base path does not exist: %s", basePath))
	}

	fmt.Printf("Loading templates from: %s\n", basePath)

	// Collect all template files
	var allFiles []string

	// Get layout files
	layoutPattern := filepath.Join(basePath, "layout", "*.html")
	layoutFiles, err := filepath.Glob(layoutPattern)
	if err != nil {
		panic(fmt.Sprintf("failed to glob layout templates: %v", err))
	}
	allFiles = append(allFiles, layoutFiles...)
	fmt.Printf("  ✓ Found %d layout template(s)\n", len(layoutFiles))

	// Get page files
	pagesPattern := filepath.Join(basePath, "pages", "*.html")
	pageFiles, err := filepath.Glob(pagesPattern)
	if err != nil {
		panic(fmt.Sprintf("failed to glob page templates: %v", err))
	}
	allFiles = append(allFiles, pageFiles...)
	fmt.Printf("  ✓ Found %d page template(s)\n", len(pageFiles))

	// Get partial files (optional)
	partialsPattern := filepath.Join(basePath, "partials", "*.html")
	partialFiles, _ := filepath.Glob(partialsPattern)
	allFiles = append(allFiles, partialFiles...)
	if len(partialFiles) > 0 {
		fmt.Printf("  ✓ Found %d partial template(s)\n", len(partialFiles))
	}

	// Parse all templates together so they can reference each other
	if len(allFiles) == 0 {
		panic("no template files found")
	}

	tmpl := template.New("").Funcs(funcMap)
	tmpl, err = tmpl.ParseFiles(allFiles...)
	if err != nil {
		panic(fmt.Sprintf("failed to parse templates: %v", err))
	}

	// List all loaded templates
	fmt.Println("\nAvailable templates:")
	for _, t := range tmpl.Templates() {
		if t.Name() != "" {
			fmt.Printf("  - %s\n", t.Name())
		}
	}

	return &TemplateRenderer{
		templates: tmpl,
	}
}

// Render returns a template document.
func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, _ echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

// getInitials is a helper function to get initials from a full name.
func getInitials(fullname string) string {
	fullname = strings.TrimSpace(fullname)
	if fullname == "" {
		return "?"
	}

	parts := strings.Fields(fullname)
	if len(parts) == 1 {
		return strings.ToUpper(string([]rune(parts[0])[0]))
	}

	if len(parts) > 1 {
		first := []rune(parts[0])[0]
		last := []rune(parts[len(parts)-1])[0]
		return strings.ToUpper(string(first) + string(last))
	}

	return "?"
}
