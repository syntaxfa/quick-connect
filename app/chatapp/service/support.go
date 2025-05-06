package service

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"log/slog"
	"time"
)

type WSSupport struct {
	id       string
	hub      *Hub
	conn     Connection
	send     chan Message
	username string
	logger   *slog.Logger
	cfg      Config
}

func NewSupport(cfg Config, id string, hub *Hub, conn Connection, username string, logger *slog.Logger) *WSSupport {
	return &WSSupport{
		cfg:      cfg,
		id:       id,
		hub:      hub,
		conn:     conn,
		send:     make(chan Message, 256),
		username: username,
		logger:   logger,
	}
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
			errlog.ErrLog(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), c.logger)
		}
	}()

	for {
		_, message, rErr := c.conn.ReadMessage()
		if rErr != nil {
			errlog.ErrLog(richerror.New(op).WithWrapError(rErr).WithKind(richerror.KindUnexpected), c.logger)
			break
		}

		var msg Message
		if uErr := json.Unmarshal(message, &msg); uErr != nil {
			errlog.ErrLog(richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected), c.logger)
			continue
		}

		msg.Sender = c.id

		c.hub.broadcastToClient <- msg
	}
}

func (c *WSSupport) WritePump() {
	const op = "service.support.WritePump"

	ticker := time.NewTicker(c.cfg.PingPeriod)

	defer func() {
		ticker.Stop()
		if cErr := c.conn.Close(); cErr != nil {
			errlog.ErrLog(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), c.logger)
		}
	}()

	for {
		select {
		case message, ok := <-c.send:
			if !ok {
				if wErr := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); wErr != nil {
					errlog.ErrLog(richerror.New(op).WithMessage("error sending websocket close message").WithKind(richerror.KindUnexpected), c.logger)
				}

				return
			}

			msgBytes, mErr := json.Marshal(message)
			if mErr != nil {
				errlog.ErrLog(richerror.New(op).WithWrapError(mErr).WithMessage("error marshalling message"), c.logger)

				continue
			}

			if wErr := c.conn.WriteMessage(websocket.TextMessage, msgBytes); wErr != nil {
				errlog.ErrLog(richerror.New(op).WithWrapError(wErr).WithMessage("error writing message"), c.logger)
			}

		case <-ticker.C:
			if wErr := c.conn.WriteMessage(websocket.PingMessage, nil); wErr != nil {
				return
			}
		}
	}
}
