package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

// UserDetail docs
// @Router /users/{userID} [GET]
// @Security JWT
// @Summary get user detail by superuser
// @Description get user detail by superuser
// @Tags User
// @Accept json
// @Produce json
// @Param userID path string true "ID of the user to update"
// @Success 200 {object} userservice.UserProfileResponse
// @Failure 404 {string} user not found
// @Failure 500 {string} something went wrong.
func (h Handler) UserDetail(c echo.Context) error {
	resp, sErr := h.userSvc.UserProfile(c.Request().Context(), types.ID(c.Param("userID")))
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
