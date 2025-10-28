package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// ShowDashboard renders the main dashboard page with service hub
func (h Handler) ShowDashboard(c echo.Context) error {
	// TODO: Get user info from context/session
	// user := getUserFromContext(c)

	return c.Render(http.StatusOK, "main_layout", nil)
}

// ShowSupportService renders the support service page
func (h Handler) ShowSupportService(c echo.Context) error {
	// TODO: Fetch support tickets data
	data := map[string]interface{}{
		"Title": "Support Management",
		"Stats": map[string]interface{}{
			"Active":  89,
			"Pending": 12,
			"Closed":  456,
		},
	}

	return c.Render(http.StatusOK, "support_page", data)
}

// ShowNotificationService renders the notification service page
func (h Handler) ShowNotificationService(c echo.Context) error {
	// TODO: Fetch notification data
	data := map[string]interface{}{
		"Title": "Notification Management",
		"Stats": map[string]interface{}{
			"Sent":      1234,
			"Delivered": 98,
			"Failed":    2,
		},
	}

	return c.Render(http.StatusOK, "notification_page", data)
}

// ShowStoryService renders the story service page
func (h Handler) ShowStoryService(c echo.Context) error {
	// TODO: Fetch story data
	data := map[string]interface{}{
		"Title": "Story Management",
		"Stats": map[string]interface{}{
			"Active": 456,
			"Views":  2300000,
		},
	}

	return c.Render(http.StatusOK, "story_page", data)
}
