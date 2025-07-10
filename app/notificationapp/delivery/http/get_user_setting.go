package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// getUserSettingAdmin docs
// @Router /v1/settings/{externalUserID} [GET]
// @Summary retrieve user setting
// @Description retrieve user settings
// @Tags NotificationAdmin
// @Accept json
// @Produce json
// @Param externalUserID path string true "ID of the template to update"
// @Success 200 {object} service.UserSetting
// @Failure 500 {string} something went wrong.
func (h Handler) getUserSettingAdmin(c echo.Context) error {
	resp, sErr := h.svc.GetUserSetting(c.Request().Context(), c.Param("externalUserID"))
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}

// getUserSetting docs
// @Router /v1/settings [GET]
// @Summary retrieve user setting
// @Description retrieve user settings
// @Tags NotificationClient
// @Accept json
// @Produce json
// @Success 200 {string} string "marked as read"
// @Success 401 {string} unauthorized
// @Failure 500 {string} something went wrong.
func (h Handler) getUserSettingClient(c echo.Context) error {
	externalUserID, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user id is not valid")
	}

	resp, sErr := h.svc.GetUserSetting(c.Request().Context(), externalUserID)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
