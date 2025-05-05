package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"time"
)

func (s Service) BroadcastNewClientJoin(requestID string) error {
	message := Message{
		Type:      MessageTypeNewClientJoin,
		Body:      "new client join to chat",
		Meta:      map[string]string{"request_id": requestID},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	data, mErr := json.Marshal(message)
	if mErr != nil {
		s.logger.Error(mErr.Error())
	}

	for key, conn := range s.supports {
		if wErr := conn.WriteMessage(websocket.TextMessage, data); wErr != nil {
			s.logger.Error(wErr.Error())

			delete(s.supports, key)

			if cErr := conn.Close(); cErr != nil {
				s.logger.Error(cErr.Error())
			}
		}
	}

	return nil
}

func (s Service) SupportStartChat(conn *websocket.Conn) {
	id := uuid.New().String()
	s.supports[id] = conn

	message := Message{
		Type:      MessageTypeNewID,
		Body:      id,
		Meta:      nil,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	data, mErr := json.Marshal(message)
	if mErr != nil {
		s.logger.Error(mErr.Error())

		delete(s.supports, id)

		return
	}

	if wErr := conn.WriteMessage(websocket.TextMessage, data); wErr != nil {
		s.logger.Error(mErr.Error())

		delete(s.clients, id)

		return
	}

	for {
		_, data, rErr := conn.ReadMessage()
		if rErr != nil {
			s.logger.Error(rErr.Error())
		}

		var message Message
		if uErr := json.Unmarshal(data, &message); uErr != nil {
			s.logger.Error(uErr.Error())
		}

		clientID, _ := message.Meta["id"]

		s.SendMessageToClient(clientID, data)
	}
}

func (s Service) SendMessageToClient(id string, data []byte) {
	conn, ok := s.clients[id]
	if ok {
		if wErr := conn.WriteMessage(websocket.TextMessage, data); wErr != nil {
			s.logger.Error(wErr.Error())
		}
	}
}
