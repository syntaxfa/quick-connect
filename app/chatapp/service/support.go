package service

import (
	"log/slog"
)

const supportChannelSize = 256

type WSSupport struct {
	wsUser // Embed the common websocket user logic
}

func NewSupport(cfg Config, id string, hub *Hub, conn Connection, username string, logger *slog.Logger) *WSSupport {
	wsSupport := &WSSupport{
		wsUser: newWsUser(cfg, id, hub, conn, username, logger, supportChannelSize),
	}
	return wsSupport
}

// ReadPump defines the specific read behavior for a WSSupport.
func (c *WSSupport) ReadPump() {
	c.readPump("service.support.ReadPump", c.hub.broadcastToClient, func() {
		c.hub.UnregisterClient(c)
	})
}

// WritePump defines the specific write behavior for a WSSupport.
func (c *WSSupport) WritePump() {
	c.writePump("service.support.WritePump")
}
