package service

import (
	"encoding/json"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"log/slog"
)

const clientChannelSize = 256

type WSClient struct {
	id           string
	hub          *Hub
	conn         Connection
	send         chan Message
	username     string
	logger       *slog.Logger
	cfg          Config
	wsConnection *wsConnection
}

func NewClient(cfg Config, id string, hub *Hub, conn Connection, username string, logger *slog.Logger) *WSClient {
	wsClient := &WSClient{
		cfg:      cfg,
		id:       id,
		hub:      hub,
		conn:     conn,
		send:     make(chan Message, clientChannelSize),
		username: username,
		logger:   logger,
	}

	wsClient.wsConnection = newWsConnection(cfg, conn, wsClient.send, logger)

	return wsClient
}

func (c *WSClient) GetID() string {
	return c.id
}

func (c *WSClient) Send(message Message) {
	c.send <- message
}

func (c *WSClient) Close() {
	close(c.send)
}

func (c *WSClient) ReadPump() {
	const op = "service.client.ReadPump"

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

		c.hub.broadcastToSupport <- msg
	}
}

func (c *WSClient) WritePump() {
	c.wsConnection.writePump("service.client.WritePump")
}
