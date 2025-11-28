package http

import (
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
)

type ChatHistoryJSONResponse struct {
	Messages   []MessageJSON `json:"messages"`
	NextCursor string        `json:"next_cursor"`
	HasMore    bool          `json:"has_more"`
}

type MessageJSON struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	SenderID  string `json:"sender_id"`
	CreatedAt string `json:"created_at"`
}

const messageLimitNumber = 30

// GetConversationModal renders the chat modal with details and history.
func (h Handler) GetConversationModal(c echo.Context) error {
	ctx := grpcContext(c)
	user, _ := getUserFromContext(c)
	convID := c.Param("id")

	// 1. Get Conversation Details
	detailReq := &conversationpb.ConversationDetailRequest{
		ConversationId: convID,
	}
	detailResp, err := h.conversationSvc.ConversationDetail(ctx, detailReq)
	if err != nil {
		return h.renderGRPCError(c, "GetConversationDetail", err)
	}

	// 2. Get Chat History
	historyReq := &conversationpb.ChatHistoryRequest{
		ConversationId: convID,
		Limit:          messageLimitNumber,
	}
	historyResp, err := h.conversationSvc.ChatHistory(ctx, historyReq)
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
		"NextCursor":   historyResp.GetNextCursor(),
		"HasMore":      historyResp.GetHasMore(),
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

	_, err := h.conversationSvc.OpenConversation(ctx, req)
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

	_, err := h.conversationSvc.CloseConversation(ctx, req)
	if err != nil {
		return h.renderGRPCError(c, "ResolveConversation", err)
	}

	// Re-render the modal to show it's closed.
	return h.GetConversationModal(c)
}

// GetChatHistory returns older messages as JSON for pagination.
func (h Handler) GetChatHistory(c echo.Context) error {
	ctx := grpcContext(c)
	convID := c.Param("id")
	cursor := c.QueryParam("cursor")

	req := &conversationpb.ChatHistoryRequest{
		ConversationId: convID,
		Cursor:         cursor,
		Limit:          int32(messageLimitNumber),
	}

	resp, err := h.conversationSvc.ChatHistory(ctx, req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	// Convert protobuf messages to JSON struct
	messages := make([]MessageJSON, 0, len(resp.GetResults()))
	for _, msg := range resp.GetResults() {
		createdAt := ""
		if msg.GetCreatedAt() != nil {
			createdAt = msg.GetCreatedAt().AsTime().Format(time.RFC3339)
		}
		messages = append(messages, MessageJSON{
			ID:        msg.GetId(),
			Content:   msg.GetContent(),
			SenderID:  msg.GetSenderId(),
			CreatedAt: createdAt,
		})
	}

	return c.JSON(http.StatusOK, ChatHistoryJSONResponse{
		Messages:   messages,
		NextCursor: resp.GetNextCursor(),
		HasMore:    resp.GetHasMore(),
	})
}
