package service

import (
	"github.com/syntaxfa/quick-connect/pkg/paginate/cursorbased"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	"github.com/syntaxfa/quick-connect/types"
)

// ListConversationsRequest defines parameters for listing conversations.
type ListConversationsRequest struct {
	AssignedSupportID types.ID             `json:"-"`
	Statuses          []ConversationStatus `json:"statuses"`
	Paginated         paginate.RequestBase `json:"paginated"`
}

// ListConversationsResponse holds the result of listing conversations.
type ListConversationsResponse struct {
	Results  []Conversation        `json:"results"`
	Paginate paginate.ResponseBase `json:"paginate"`
}

type ChatHistoryRequest struct {
	UserID         types.ID            `json:"-"`
	UserRoles      []types.Role        `json:"-"`
	ConversationID types.ID            `json:"conversation_id"`
	Pagination     cursorbased.Request `json:"pagination"`
}

type ChatHistoryResponse struct {
	Results  []Message            `json:"results"`
	Paginate cursorbased.Response `json:"paginate"`
}
