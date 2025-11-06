package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// RegisterGuestUser docs
// @Summary RegisterGuestUser
// @Description register guest user and generate QCToken (long expire time)
// @Tags User
// @Accept json
// @Produce json
// @Param Request body userservice.RegisterGuestUserRequest true "check token validation"
// @Success 201 {object} userservice.RegisterGuestUserResponse
// @Failure 400 {string} string Bad Request
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong
// @Router /users/register-guest [POST].
func (h Handler) RegisterGuestUser(c echo.Context) error {
	var req userservice.RegisterGuestUserRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.userSvc.RegisterGuestUser(c.Request().Context(), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusCreated, resp)
}
