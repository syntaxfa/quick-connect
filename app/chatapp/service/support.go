package service

import (
	"encoding/json"
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

const supportChannelSize = 256

type WSSupport struct {
	id           string
	hub          *Hub
	conn         Connection
	send         chan Message
	username     string
	logger       *slog.Logger
	cfg          Config
	wsConnection *wsConnection
}

func NewSupport(cfg Config, id string, hub *Hub, conn Connection, username string, logger *slog.Logger) *WSSupport {
	wsSupport := &WSSupport{
		cfg:      cfg,
		id:       id,
		hub:      hub,
		conn:     conn,
		send:     make(chan Message, supportChannelSize),
		username: username,
		logger:   logger,
	}

	wsSupport.wsConnection = newWsConnection(cfg, conn, wsSupport.send, logger)

	return wsSupport
}

func (c *WSSupport) GetID() string {
	return c.id
}

func (c *WSSupport) Send(message Message) {
	c.send <- message
}

func (c *WSSupport) Close() {
	close(c.send)
}

func (c *WSSupport) ReadPump() {
	const op = "service.support.ReadPump"

	defer func() {
		c.hub.UnregisterClient(c)
		if cErr := c.conn.Close(); cErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), c.logger)
		}
	}()

	for {
		_, message, rErr := c.conn.ReadMessage()
		if rErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected), c.logger)

			break
		}

		var msg Message
		if uErr := json.Unmarshal(message, &msg); uErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected), c.logger)

			continue
		}

		msg.Sender = c.id

		c.hub.broadcastToClient <- msg
	}
}

func (c *WSSupport) WritePump() {
	c.wsConnection.writePump("service.support.WritePump")
}
