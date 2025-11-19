package service

import (
	"time"

	"github.com/syntaxfa/quick-connect/types"
)

// ClientMessage defines the structure received from a websocket client.
type ClientMessage struct {
	Type            MessageType `json:"type"`
	SubType         string      `json:"sub_type"`
	ConversationID  types.ID    `json:"conversation_id"`
	Content         string      `json:"content"`
	ClientMessageID string      `json:"client_message_id"`
	SenderID        types.ID    `json:"-"` // Filled by the service
}

// ServerMessage defines the structure sent to a websocket client.
type ServerMessage struct {
	Type            MessageType `json:"type"`
	SubType         string      `json:"sub_type"`
	Timestamp       time.Time   `json:"timestamp"`
	Payload         interface{} `json:"payload"`
	ClientMessageID string      `json:"client_message_id"`
}

// SystemMessagePayload defines the payload for system-type ServerMessages.
type SystemMessagePayload struct {
	ConversationID types.ID `json:"conversation_id"`
	SenderID       types.ID `json:"sender_id"`
}

// PubSubMessage is the wrapper for messages sent over the Pub/Sub system.
type PubSubMessage struct {
	RecipientIDs []types.ID    `json:"recipient_ids"` // Target users for this message
	ServerMsg    ServerMessage `json:"server_msg"`    // The actual message to send
	ExcludeIDs   []types.ID    `json:"exclude_ids"`   // Users to explicitly exclude (for system messages)
}

// NewTextMessage creates a "text" type ServerMessage from a database Message object.
func NewTextMessage(dbMessage Message, clientMessageID string) ServerMessage {
	return ServerMessage{
		Type:            MessageTypeText,
		Timestamp:       dbMessage.CreatedAt,
		Payload:         dbMessage,
		ClientMessageID: clientMessageID,
	}
}

// NewSystemMessage creates a "system" type ServerMessage.
func NewSystemMessage(convoID, senderID types.ID, subType string) ServerMessage {
	return ServerMessage{
		Type:      MessageTypeSystem,
		SubType:   subType,
		Timestamp: time.Now(),
		Payload: SystemMessagePayload{
			ConversationID: convoID,
			SenderID:       senderID,
		},
	}
}
