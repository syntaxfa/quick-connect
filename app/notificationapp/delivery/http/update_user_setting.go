package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// updateUserSettingAdmin docs
// @Router /v1/settings/{externalUserID} [POST]
// @Summary update user setting
// @Description This API endpoint updates user notification settings
// @Tags NotificationAdmin
// @Accept json
// @Produce json
// @Param externalUserID path string true "ID of the template to update"
// @Param Request body service.UpdateUserSettingRequest true "user settings"
// @Success 200 {object} service.UserSetting
// @Failure 400 {string} string Bad Request
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong.
func (h Handler) updateUserSettingAdmin(c echo.Context) error {
	var req service.UpdateUserSettingRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.svc.UpdateUserSetting(c.Request().Context(), c.Param("externalUserID"), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}

// updateUserSettingClient docs
// @Router /v1/settings [POST]
// @Summary update user setting
// @Description This API endpoint updates user notification settings
// @Tags NotificationClient
// @Accept json
// @Produce json
// @Param Request body service.UpdateUserSettingRequest true "user settings"
// @Success 200 {object} service.UserSetting
// @Failure 400 {string} string Bad Request
// @Failure 401 {string} unauthorized
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong.
func (h Handler) updateUserSettingClient(c echo.Context) error {
	externalUserID, ok := c.Get("user_id").(string)
	if !ok {
		return echo.NewHTTPError(http.StatusUnauthorized, "user id is not valid")
	}

	var req service.UpdateUserSettingRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.svc.UpdateUserSetting(c.Request().Context(), externalUserID, req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
