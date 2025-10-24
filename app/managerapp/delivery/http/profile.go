package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// UserProfile docs
// @Summary UserProfile
// @Description get user profile
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} userservice.UserProfileResponse
// @Failure 500 {string} something went wrong
// @Security JWT
// @Router /users/profile [GET].
func (h Handler) UserProfile(c echo.Context) error {
	userClaims, gErr := auth.GetUserClaimFormContext(c)
	if gErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, gErr.Error())
	}

	user, sErr := h.userSvc.UserProfile(c.Request().Context(), userClaims.UserID)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, user)
}
