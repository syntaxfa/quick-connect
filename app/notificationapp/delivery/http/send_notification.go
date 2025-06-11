package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// sendNotification docs
// @Summary send notification
// @Description This API endpoint send a new notification.
// @Tags Notification
// @Accept json
// @Produce json
// @Param Request body service.SendNotificationRequest true "generate pair(refresh & access) tokens"
// @Success 200 {object} service.SendNotificationResponse
// @Failure 400 {string} string Bad Request
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong
// @Router /notifications [POST].
func (h Handler) sendNotification(c echo.Context) error {
	var req service.SendNotificationRequest

	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.svc.SendNotification(c.Request().Context(), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusCreated, resp)
}
