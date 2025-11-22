package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"github.com/syntaxfa/quick-connect/pkg/auth"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
)

// GetChatHistory docs
// @Summary GetChatHistory
// @Description get conversation chat history
// @Tags Chat
// @Accept json
// @Produce json
// @Param Request body service.ChatHistoryRequest true "get conversation chat history"
// @Success 200 {object} service.ChatHistoryResponse
// @Success 400 {string} bad request
// @Failure 403 {string} you not participant in this conversation
// @Failure 500 {string} something went wrong
// @Security JWT
// @Router /chats [POST].
func (h Handler) GetChatHistory(c echo.Context) error {
	var req service.ChatHistoryRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	claims, cErr := auth.GetUserClaimFormContext(c)
	if cErr != nil {
		return echo.NewHTTPError(http.StatusUnauthorized, cErr.Error())
	}

	req.UserID = claims.UserID
	req.UserRoles = claims.Roles

	resp, sErr := h.svc.ChatHistory(c.Request().Context(), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
