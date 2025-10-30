package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

// UserUpdateFormSuperuser docs
// @Router /users/{userID} [PUT]
// @Security JWT
// @Summary update user by superuser
// @Description update user by superuser
// @Tags User
// @Accept json
// @Produce json
// @Param userID path string true "ID of the user to update"
// @Param Request body userservice.UserUpdateFromSuperuserRequest true "check token validation"
// @Success 200 {object} userservice.UserUpdateResponse
// @Failure 404 {string} user not found
// @Failure 409 {string} This username already exists
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong.
func (h Handler) UserUpdateFormSuperuser(c echo.Context) error {
	var req userservice.UserUpdateFromSuperuserRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.userSvc.UserUpdateFromSuperuser(c.Request().Context(), types.ID(c.Param("userID")), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
