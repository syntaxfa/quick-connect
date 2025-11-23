package service

import (
	"time"

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

type ClientInfo struct {
	ID           types.ID  `json:"id"`
	Fullname     string    `json:"fullname,omitempty"`
	PhoneNumber  string    `json:"phone_number,omitempty"`
	Email        string    `json:"email,omitempty"`
	Avatar       string    `json:"avatar"`
	LastOnlineAt time.Time `json:"last_online_at"`
}

type SupportInfo struct {
	ID           types.ID  `json:"id"`
	Fullname     string    `json:"fullname,omitempty"`
	Avatar       string    `json:"avatar"`
	LastOnlineAt time.Time `json:"last_online_at"`
}

type ConversationDetailResponse struct {
	Conversation

	ClientInfo  ClientInfo  `json:"client_info"`
	SupportInfo SupportInfo `json:"support_info"`
}
