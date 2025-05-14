package service

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

type WSClient struct {
	id       string
	hub      *Hub
	conn     Connection
	send     chan Message
	username string
	logger   *slog.Logger
	cfg      Config
}

func NewClient(cfg Config, id string, hub *Hub, conn Connection, username string, logger *slog.Logger) *WSClient {
	return &WSClient{
		cfg:      cfg,
		id:       id,
		hub:      hub,
		conn:     conn,
		send:     make(chan Message, 256),
		username: username,
		logger:   logger,
	}
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
	const op = "service.client.WritePump"

	ticker := time.NewTicker(c.cfg.PingPeriod)

	defer func() {
		ticker.Stop()
		if cErr := c.conn.Close(); cErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), c.logger)
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				if wErr := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); wErr != nil {
					errlog.WithoutErr(richerror.New(op).WithMessage("error sending websocket close message").WithKind(richerror.KindUnexpected), c.logger)
				}

				return
			}

			msgBytes, mErr := json.Marshal(message)
			if mErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(mErr).WithMessage("error marshalling message"), c.logger)

				continue
			}

			if wErr := c.conn.WriteMessage(websocket.TextMessage, msgBytes); wErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(wErr).WithMessage("error writing message"), c.logger)
			}

		case <-ticker.C:
			if wErr := c.conn.WriteMessage(websocket.PingMessage, nil); wErr != nil {
				return
			}
		}
	}
}
