package service

import (
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

type ConversationStatus string

const (
	ConversationStatusNew         ConversationStatus = "new"
	ConversationStatusOpen        ConversationStatus = "open"
	ConversationStatusClosed      ConversationStatus = "closed"
	ConversationStatusBotHandling ConversationStatus = "bot_handling"
)

var AllConversationStatus = []ConversationStatus{
	ConversationStatusNew,
	ConversationStatusOpen,
	ConversationStatusClosed,
	ConversationStatusBotHandling,
}

func IsValidConversationStatus(conversationStatus ConversationStatus) bool {
	for _, con := range AllConversationStatus {
		if conversationStatus == con {
			return true
		}
	}

	return false
}

type MessageType string

const (
	MessageTypeText   MessageType = "text"
	MessageTypeMedia  MessageType = "media"
	MessageTypeSystem MessageType = "system"
)

var AllMessageType = []MessageType{
	MessageTypeText,
	MessageTypeMedia,
	MessageTypeSystem,
}

func IsValidMessageType(messageType MessageType) bool {
	for _, msgType := range AllMessageType {
		if messageType == msgType {
			return true
		}
	}

	return false
}

type Conversation struct {
	ID                  types.ID           `json:"id"`
	ClientUserID        types.ID           `json:"client_user_id"`
	AssignedSupportID   types.ID           `json:"assigned_support_id"`
	Status              ConversationStatus `json:"status"`
	LastMessageSnippet  string             `json:"last_message_snippet"`
	LastMessageSenderID types.ID           `json:"last_message_sender_id"`
	CreatedAt           time.Time          `json:"created_at"`
	UpdatedAt           time.Time          `json:"updated_at"`
	ClosedAt            *time.Time         `json:"closed_at"`
}

type Message struct {
	ID                 types.ID          `json:"id"`
	ConversationID     types.ID          `json:"conversation_id"`
	SenderID           types.ID          `json:"sender_id"`
	MessageType        MessageType       `json:"message_type"`
	Content            string            `json:"content"`
	Metadata           map[string]string `json:"metadata"`
	RepliedToMessageID types.ID          `json:"replied_to_message_id"`
	CreatedAt          time.Time         `json:"created_at"`
	ReadAt             *time.Time        `json:"read_at"`
}
