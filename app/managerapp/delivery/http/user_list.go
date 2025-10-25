package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// UserList docs
// @Summary get user list
// @Description This API endpoint get user list.
// @Tags User
// @Accept json
// @Produce json
// @Param Request body userservice.ListUserRequest true "list of users"
// @Success 200 {object} userservice.ListUserResponse
// @Failure 400 {string} string Bad Request
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong
// @Security JWT
// @Router /users/list [POST].
func (h Handler) UserList(c echo.Context) error {
	var req userservice.ListUserRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.userSvc.UserList(c.Request().Context(), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
