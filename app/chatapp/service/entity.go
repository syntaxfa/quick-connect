package service

type Message struct {
	Type      string `json:"type"`
	Content   string `json:"content"`
	Sender    string `json:"sender"`
	Recipient string `json:"recipient"`
}
