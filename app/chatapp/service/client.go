package service

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"time"
)

func (s Service) ClientStartChat(conn *websocket.Conn) {
	id := uuid.New().String()
	s.clients[id] = conn

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

		delete(s.clients, id)

		return
	}

	if wErr := conn.WriteMessage(websocket.TextMessage, data); wErr != nil {
		s.logger.Error(mErr.Error())

		delete(s.clients, id)

		return
	}

	if bErr := s.BroadcastNewClientJoin(id); bErr != nil {
		s.logger.Error(bErr.Error())
	}

	for {
		messageType, data, rErr := conn.ReadMessage()
		if rErr != nil {
			s.logger.Error(rErr.Error())
		}

		if wErr := conn.WriteMessage(messageType, data); wErr != nil {
			s.logger.Error(wErr.Error())
		}
	}
}
