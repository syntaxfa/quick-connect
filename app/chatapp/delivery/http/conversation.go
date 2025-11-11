package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// GetActiveConversation docs
// @Summary GetActiveConversation
// @Description get user active conversation
// @Tags Chat
// @Accept json
// @Produce json
// @Success 200 {object} service.Conversation
// @Failure 500 {string} something went wrong
// @Security JWT
// @Router /conversations/active [GET].
func (h Handler) GetActiveConversation(c echo.Context) error {
	claims, cErr := auth.GetUserClaimFormContext(c)
	if cErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, cErr.Error())
	}

	resp, sErr := h.svc.GetUserActiveConversation(c.Request().Context(), claims.UserID)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
