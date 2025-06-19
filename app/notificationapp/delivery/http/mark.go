package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

// MarkNotificationAsRead docs
// @Summary mark notification as read
// @Description mark notification as read.
// @Tags NotificationClient
// @Accept json
// @Produce json
// @Param notificationID path string true "notification id"
// @Success 200 {string} string "marked as read"
// @Failure 500 {string} something went wrong
// @Router /v1/notifications/{notificationID}/mark-as-read [GET].
func (h Handler) markNotificationAsRead(c echo.Context) error {
	if sErr := h.svc.MarkNotificationAsRead(c.Request().Context(), types.ID(c.Param("notificationID"))); sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, "")
}

// MarkNotificationAsRead docs
// @Summary mark all notification as read
// @Description mark all  notification as read.
// @Tags NotificationClient
// @Accept json
// @Produce json
// @Param externalUserID path string true "external user id"
// @Success 200 {string} string "marked all as read"
// @Failure 500 {string} something went wrong
// @Router /v1/notifications/{externalUserID}/mark-all-as-read [GET].
func (h Handler) markAllNotificationAsRead(c echo.Context) error {
	if sErr := h.svc.MarkAllNotificationAsRead(c.Request().Context(), c.Param("externalUserID")); sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, "")
}
