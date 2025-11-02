package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ShowSettingsPage renders the settings page
func (h Handler) ShowSettingsPage(c echo.Context) error {
	user, _ := getUserFromContext(c)

	data := map[string]interface{}{
		"Title":        "Settings Management",
		"TemplateName": "settings_page",
		"User":         user,
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "settings_page", data)
	}

	return c.Render(http.StatusOK, "main_layout", data)
}
