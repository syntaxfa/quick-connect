package service

import "time"

type Message struct {
	Type      MessageType       `json:"type"`
	Body      string            `json:"body"`
	Meta      map[string]string `json:"meta"`
	CreatedAt time.Time         `json:"created_at"`
	UpdatedAt time.Time         `json:"updated_at"`
}

type MessageType string

const (
	MessageTypeText          MessageType = "text"
	MessageTypeNewID         MessageType = "new_id"
	MessageTypeNewClientJoin MessageType = "new_client_join"
	MessageTypeSupportOnline MessageType = "support_online"
	MessageTypeClientOnline  MessageType = "client_online"
)
