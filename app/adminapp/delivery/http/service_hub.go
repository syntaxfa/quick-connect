package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ShowDashboard renders the main dashboard page with service hub
func (h Handler) ShowDashboard(c echo.Context) error {
	user, _ := getUserFromContext(c)

	data := map[string]interface{}{
		"TemplateName": "dashboard_page",
		"User":         user,
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "dashboard_page", data)
	}

	return c.Render(http.StatusOK, "main_layout", data)
}

// ShowSupportService renders the support service page
func (h Handler) ShowSupportService(c echo.Context) error {
	user, _ := getUserFromContext(c)

	// TODO: Fetch support tickets data
	data := map[string]interface{}{
		"Title": "Support Management",
		"Stats": map[string]interface{}{
			"Active":  89,
			"Pending": 12,
			"Closed":  456,
		},
		"TemplateName": "support_page",
		"User":         user,
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "support_page", data)
	}

	return c.Render(http.StatusOK, "main_layout", data)
}

// ShowNotificationService renders the notification service page
func (h Handler) ShowNotificationService(c echo.Context) error {
	user, _ := getUserFromContext(c)

	// TODO: Fetch notification data
	data := map[string]interface{}{
		"Title": "Notification Management",
		"Stats": map[string]interface{}{
			"Sent":      1234,
			"Delivered": 98,
			"Failed":    2,
		},
		"TemplateName": "notification_page",
		"User":         user,
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "notification_page", data)
	}

	return c.Render(http.StatusOK, "main_layout", data)
}

// ShowStoryService renders the story service page
func (h Handler) ShowStoryService(c echo.Context) error {
	user, _ := getUserFromContext(c)

	// TODO: Fetch story data
	data := map[string]interface{}{
		"Title": "Story Management",
		"Stats": map[string]interface{}{
			"Active": 456,
			"Views":  2300000,
		},
		"TemplateName": "story_page",
		"User":         user,
	}

	if isHTMX(c) {
		return c.Render(http.StatusOK, "story_page", data)
	}

	return c.Render(http.StatusOK, "main_layout", data)
}
