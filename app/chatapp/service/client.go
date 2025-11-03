package service

import (
	"log/slog"
)

const clientChannelSize = 256

type WSClient struct {
	wsUser // Embed the common websocket user logic
}

func NewClient(cfg Config, id string, hub *Hub, conn Connection, username string, logger *slog.Logger) *WSClient {
	wsClient := &WSClient{
		wsUser: newWsUser(cfg, id, hub, conn, username, logger, clientChannelSize),
	}
	return wsClient
}

// ReadPump defines the specific read behavior for a WSClient.
func (c *WSClient) ReadPump() {
	c.readPump("service.client.ReadPump", c.hub.broadcastToSupport, func() {
		c.hub.UnregisterClient(c)
	})
}

// WritePump defines the specific write behavior for a WSClient.
func (c *WSClient) WritePump() {
	c.writePump("service.client.WritePump")
}

// GetID, Send, and Close are inherited from wsUser.
