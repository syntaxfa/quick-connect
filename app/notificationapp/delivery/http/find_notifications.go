package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// FindNotifications docs
// @Summary find user notifications
// @Description This API endpoint find an userID notifications.
// @Tags NotificationClient
// @Accept json
// @Produce json
// @Param Request body service.ListNotificationRequest true "find user notifications"
// @Success 200 {object} service.ListNotificationResponse
// @Failure 400 {string} string Bad Request
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong
// @Router /v1/notifications/list [POST].
func (h Handler) findNotifications(c echo.Context) error {
	var req service.ListNotificationRequest

	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	userID, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user id is not valid")
	}

	req.ExternalUserID = userID

	resp, sErr := h.svc.FindNotificationByUserID(c.Request().Context(), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
