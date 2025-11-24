package http

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
)

const messageLimitNumber = 50

// GetConversationModal renders the chat modal with details and history.
func (h Handler) GetConversationModal(c echo.Context) error {
	ctx := grpcContext(c)
	user, _ := getUserFromContext(c)
	convID := c.Param("id")

	// 1. Get Conversation Details
	detailReq := &conversationpb.ConversationDetailRequest{
		ConversationId: convID,
	}
	detailResp, err := h.conversationAd.ConversationDetail(ctx, detailReq)
	if err != nil {
		return h.renderGRPCError(c, "GetConversationDetail", err)
	}

	// 2. Get Chat History
	historyReq := &conversationpb.ChatHistoryRequest{
		ConversationId: convID,
		Limit:          messageLimitNumber,
	}
	historyResp, err := h.conversationAd.ChatHistory(ctx, historyReq)
	if err != nil {
		h.logger.Error("Failed to fetch chat history", "error", err)
	}

	conv := convertConversationPbToConversation(detailResp.GetConversation())

	// Logic to determine state
	// If assigned ID is empty or "0", it's unassigned.
	isUnassigned := conv.AssignedSupportID == "" || conv.AssignedSupportID == "0"

	// Is this chat assigned to the current logged-in user?
	isMyChat := conv.AssignedSupportID == user.ID

	data := map[string]interface{}{
		"Conversation": conv,
		"ClientInfo":   detailResp.GetClientInfo(),
		"SupportInfo":  detailResp.GetSupportInfo(),
		"Messages":     historyResp.GetResults(),
		"IsUnassigned": isUnassigned,
		"IsMyChat":     isMyChat,
		"CurrentUser":  user,
	}

	return c.Render(http.StatusOK, "support_chat_modal", data)
}

// JoinConversation handles the "Start" button logic.
func (h Handler) JoinConversation(c echo.Context) error {
	ctx := grpcContext(c)
	convID := c.Param("id")

	req := &conversationpb.OpenConversationRequest{
		ConversationId: convID,
	}

	_, err := h.conversationAd.OpenConversation(ctx, req)
	if err != nil {
		return h.renderGRPCError(c, "JoinConversation", err)
	}

	// Re-render the modal to show the update (buttons should change).
	return h.GetConversationModal(c)
}

// ResolveConversation handles the "Close" button logic.
func (h Handler) ResolveConversation(c echo.Context) error {
	ctx := grpcContext(c)
	convID := c.Param("id")

	req := &conversationpb.CloseConversationRequest{
		ConversationId: convID,
	}

	_, err := h.conversationAd.CloseConversation(ctx, req)
	if err != nil {
		return h.renderGRPCError(c, "ResolveConversation", err)
	}

	// Re-render the modal to show it's closed.
	return h.GetConversationModal(c)
}
