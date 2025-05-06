package service

type Message struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}

type MessageType string

const (
	MessageTypeEcho = "echo"
	MessageTypeNew  = "new_message"
)
