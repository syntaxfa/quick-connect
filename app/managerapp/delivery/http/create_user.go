package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// CreateUser docs
// @Summary CreateUser
// @Description create a new user by superuser
// @Tags User
// @Accept json
// @Produce json
// @Param Request body userservice.UserCreateRequest true "check token validation"
// @Success 201 {object} userservice.UserCreateResponse
// @Failure 404 {string} This user already exists
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong
// @Security JWT
// @Router /users [POST].
func (h Handler) CreateUser(c echo.Context) error {
	var req userservice.UserCreateRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	user, sErr := h.userSvc.CreateUser(c.Request().Context(), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusCreated, user)
}
