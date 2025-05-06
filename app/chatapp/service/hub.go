package service

import (
	"log/slog"
	"sync"
)

type Client interface {
	GetID() string
	Send(message Message)
	Close()
}

type Support interface {
	GetID() string
	Send(message Message)
	Close()
}

type Hub struct {
	clients            map[string]Client
	supports           map[string]Support
	registerClient     chan Client
	registerSupport    chan Support
	unregisterClient   chan Client
	unregisterSupport  chan Support
	broadcastToSupport chan Message
	broadcastToClient  chan Message
	mu                 sync.RWMutex
	logger             *slog.Logger
}

func NewHub(logger *slog.Logger) *Hub {
	return &Hub{
		clients:            make(map[string]Client),
		supports:           make(map[string]Support),
		registerClient:     make(chan Client),
		registerSupport:    make(chan Support),
		unregisterClient:   make(chan Client),
		unregisterSupport:  make(chan Support),
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
			if client, ok := h.clients[client.GetID()]; ok {
				delete(h.clients, client.GetID())
				client.Close()
				h.logger.Info("client unregistered", slog.String("id", client.GetID()))
			}
			h.mu.Unlock()

		case support := <-h.unregisterSupport:
			h.mu.Lock()
			if support, ok := h.supports[support.GetID()]; ok {
				delete(h.supports, support.GetID())
				support.Close()
				h.logger.Info("support unregistered", slog.String("id", support.GetID()))
			}
			h.mu.Unlock()

		case message := <-h.broadcastToSupport:
			h.SendMessageToSupport(message)

		case message := <-h.broadcastToClient:
			if message.Recipient == "" {
				h.BroadcastToAllClient(message)
			} else {
				h.SendPrivateMessageToClient(message)
			}
		}
	}
}

func (h *Hub) RegisterClient(client Client) {
	h.registerClient <- client
}

func (h *Hub) RegisterSupport(support Support) {
	h.registerSupport <- support
}

func (h *Hub) UnregisterClient(client Client) {
	h.unregisterClient <- client
}

func (h *Hub) UnregisterSupport(support Support) {
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

	if sender, ok := h.clients[message.Sender]; ok {
		sender.Send(message)
	}
}

func (h *Hub) BroadcastToAllClient(message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		client.Send(message)
	}

	if sender, ok := h.supports[message.Sender]; ok {
		sender.Send(message)
	}
}

func (h *Hub) SendPrivateMessageToClient(message Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	if client, ok := h.clients[message.Recipient]; ok {
		client.Send(message)
	}

	if sender, ok := h.supports[message.Sender]; ok {
		sender.Send(message)
	}
}
