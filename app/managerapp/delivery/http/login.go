package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// UserLogin docs
// @Summary UserLogin
// @Description user log in and generate pair token(access and refresh)
// @Tags User
// @Accept json
// @Produce json
// @Param Request body userservice.UserLoginRequest true "check token validation"
// @Success 200 {object} userservice.UserLoginResponse
// @Failure 400 {string} string Bad Request
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong
// @Router /users/login [POST].
func (h Handler) UserLogin(c echo.Context) error {
	var req userservice.UserLoginRequest

	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, lErr := h.userSvc.Login(c.Request().Context(), req)
	if lErr != nil {
		return servermsg.HTTPMsg(c, lErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
