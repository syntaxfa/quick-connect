package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// ChangePassword change password.
// @Router /users/change-password [POST]
// @Security JWT
// @Summary change password
// @Description This API change user password
// @Tags User
// @Accept json
// @Produce json
// @Param Request body userservice.ChangePasswordRequest true "change password"
// @Success 200 {string} string user password changed
// @Failure 404 {string} string user not found
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong.
func (h Handler) ChangePassword(c echo.Context) error {
	var req userservice.ChangePasswordRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	userClaims, gErr := auth.GetUserClaimFormContext(c)
	if gErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, gErr.Error())
	}

	if sErr := h.userSvc.ChangePassword(c.Request().Context(), userClaims.UserID, req); sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, nil)
}
