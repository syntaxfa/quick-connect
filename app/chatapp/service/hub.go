package service

import (
	"log/slog"
	"sync"
)

type Participant interface {
	GetID() string
	Send(message Message)
	Close()
}

type Hub struct {
	clients            map[string]Participant
	supports           map[string]Participant
	registerClient     chan Participant
	registerSupport    chan Participant
	unregisterClient   chan Participant
	unregisterSupport  chan Participant
	broadcastToSupport chan Message
	broadcastToClient  chan Message
	mu                 sync.RWMutex
	logger             *slog.Logger
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients:            make(map[string]Participant),
		supports:           make(map[string]Participant),
		registerClient:     make(chan Participant),
		registerSupport:    make(chan Participant),
		unregisterClient:   make(chan Participant),
		unregisterSupport:  make(chan Participant),
		broadcastToSupport: make(chan Message),
		broadcastToClient:  make(chan Message),
		mu:                 sync.RWMutex{},
		logger:             logger,
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.registerClient:
			h.mu.Lock()
			h.clients[client.GetID()] = client
			h.mu.Unlock()
			h.logger.Info("client registered", slog.String("id", client.GetID()))

		case support := <-h.registerSupport:
			h.mu.Lock()
			h.supports[support.GetID()] = support
			h.mu.Unlock()
			h.logger.Info("support registered", slog.String("id", support.GetID()))

		case client := <-h.unregisterClient:
			h.mu.Lock()
			if conn, ok := h.clients[client.GetID()]; ok {
				delete(h.clients, client.GetID())
				conn.Close()
				h.logger.Info("client unregistered", slog.String("id", client.GetID()))
			}
			h.mu.Unlock()

		case support := <-h.unregisterSupport:
			h.mu.Lock()
			if conn, ok := h.supports[support.GetID()]; ok {
				delete(h.supports, support.GetID())
				conn.Close()
				h.logger.Info("support unregistered", slog.String("id", support.GetID()))
			}
			h.mu.Unlock()

		case message := <-h.broadcastToSupport:
			h.SendMessageToSupport(message)

			//case message := <-h.broadcastToClient:
			//	if message.Recipient == "" {
			//		h.BroadcastToAllClient(message)
			//	} else {
			//		h.SendPrivateMessageToClient(message)
			//	}
		}
	}
}

func (h *Hub) RegisterClient(client Participant) {
	h.registerClient <- client
}

func (h *Hub) RegisterSupport(support Participant) {
	h.registerSupport <- support
}

func (h *Hub) UnregisterClient(client Participant) {
	h.unregisterClient <- client
}

func (h *Hub) UnregisterSupport(support Participant) {
	h.unregisterSupport <- support
}

func (h *Hub) BroadcastMessageToSupport(message Message) {
	h.broadcastToSupport <- message
}

func (h *Hub) BroadcastMessageToClient(message Message) {
	h.broadcastToClient <- message
}

func (h *Hub) SendMessageToSupport(message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, support := range h.supports {
		support.Send(message)
	}

	//if sender, ok := h.clients[message.Sender]; ok {
	//	message.Type = MessageTypeEcho
	//	sender.Send(message)
	//}
}

func (h *Hub) BroadcastToAllClient(message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		client.Send(message)
	}

	//if sender, ok := h.supports[message.Sender]; ok {
	//	message.Type = MessageTypeEcho
	//	sender.Send(message)
	//}
}

func (h *Hub) SendPrivateMessageToClient(message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	//if client, clientOk := h.clients[message.Recipient]; clientOk {
	//	client.Send(message)
	//
	//	if sender, supportOk := h.supports[message.Sender]; supportOk {
	//		message.Type = MessageTypeEcho
	//		sender.Send(message)
	//	}
	//}
}
