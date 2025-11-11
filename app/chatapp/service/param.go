package service

import (
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
