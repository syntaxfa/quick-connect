package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/types"
)

// UserDelete deletes a user by ID.
// @Router /users/{userID} [DELETE]
// @Security JWT
// @Summary delete user
// @Description This API delete user
// @Tags User
// @Accept json
// @Produce json
// @Param userID path string true "ID of the user to delete"
// @Success 204 {string} string user deleted
// @Failure 404 {string} string user not found
// @Failure 500 {string} something went wrong
func (h Handler) UserDelete(c echo.Context) error {
	if sErr := h.userSvc.UserDelete(c.Request().Context(), types.ID(c.Param("userID"))); sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusNotFound, "")
}
