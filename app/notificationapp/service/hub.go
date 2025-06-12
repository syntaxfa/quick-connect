package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

type Connection interface {
	WriteMessage(messageType int, data []byte) error
	Close() error
	RemoteAddr() string
}

type HubConfig struct {
	UserConnectionLimit int           `koanf:"user_connection_limit"`
	WriteWait           time.Duration `koanf:"write_wait"`
	PongWait            time.Duration `koanf:"pong_wait"`
	PingPeriod          time.Duration
	MaxMessageSize      int `koanf:"max_message_size"`
}

type Hub struct {
	cfg          HubConfig
	clients      map[types.ID][]*Client
	register     chan *Client
	unregistered chan *Client
	notification chan *NotificationMessage
	logger       *slog.Logger
	mu           sync.RWMutex
}

type Client struct {
	hub    *Hub
	conn   Connection
	send   chan *NotificationMessage
	userID types.ID
}

type NotificationMessage struct {
	NotificationID types.ID         `json:"notification_id"`
	UserID         types.ID         `json:"user_id"`
	Type           NotificationType `json:"type"`
	Title          string           `json:"title"`
	Body           string           `json:"body"`
	Data           json.RawMessage  `json:"data"`
	Timestamp      int64            `json:"timestamp"`
}

func NewHub(cfg HubConfig, logger *slog.Logger) *Hub {
	return &Hub{
		cfg:          cfg,
		clients:      make(map[types.ID][]*Client),
		register:     make(chan *Client),
		unregistered: make(chan *Client),
		notification: make(chan *NotificationMessage),
		logger:       logger,
	}
}

func (h *Hub) Run() {
	const op = "service.hub.Run"

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if _, ok := h.clients[client.userID]; ok {
				if len(h.clients[client.userID]) > h.cfg.UserConnectionLimit {
					for _, connection := range h.clients[client.userID] {
						close(connection.send)
						if cErr := connection.conn.Close(); cErr != nil {
							errlog.WithoutErr(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), h.logger)
						}
					}
					delete(h.clients, client.userID)
					h.clients[client.userID] = []*Client{client}
				} else {
					h.clients[client.userID] = append(h.clients[client.userID], client)
				}
			} else {
				h.clients[client.userID] = []*Client{client}
			}
			h.logger.Info("client registered", slog.String("user_id", string(client.userID)),
				slog.String("addr", client.conn.RemoteAddr()))
			h.mu.Unlock()
		case client := <-h.unregistered:
			h.mu.Lock()
			if connections, ok := h.clients[client.userID]; ok {
				for i, connection := range connections {
					if connection == client {
						close(client.send)
						if cErr := client.conn.Close(); cErr != nil {
							errlog.WithoutErr(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), h.logger)
						}

						h.clients[client.userID] = append(h.clients[client.userID][:i], h.clients[client.userID][i+1:]...)
						h.logger.Info("client unregistered", slog.String("user_id", string(client.userID)),
							slog.String("addr", client.conn.RemoteAddr()))
					}
				}
			}
			h.mu.Unlock()
		case notification := <-h.notification:
			h.mu.RLock()
			if connections, ok := h.clients[notification.UserID]; ok {
				for _, client := range connections {
					select {
					case client.send <- notification:
						h.logger.Debug("notification sent to client", slog.String("notification_id",
							string(notification.NotificationID)))
					default:
						h.logger.Warn("failed to send notification to client, client send buffer full",
							slog.String("user_id", string(notification.UserID)))
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// WritePump pumps message from the hub to the websocket connection.
func (c *Client) WritePump() {
	const op = "service.hub.WritePump"

	ticker := time.NewTicker(c.hub.cfg.PingPeriod)

	defer func() {
		ticker.Stop()
		if cErr := c.conn.Close(); cErr != nil {
			errlog.WithoutErr(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected).
				WithMessage("failed to stop client websocket connection"), c.hub.logger)
		}
		c.hub.unregistered <- c
	}()

	for {
		select {
		case notification, ok := <-c.send:
			if !ok {
				if !ok {
					if wErr := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); wErr != nil {
						errlog.WithoutErr(richerror.New(op).WithMessage("error sending websocket close message").WithKind(richerror.KindUnexpected), c.hub.logger)
					}

					return
				}
			}

			msgBytes, mErr := json.Marshal(notification)
			if mErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(mErr).WithMessage("error marshalling notification"), c.hub.logger)

				continue
			}

			if wErr := c.conn.WriteMessage(websocket.TextMessage, msgBytes); wErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(wErr).WithMessage("error writing message"), c.hub.logger)
			}
		case <-ticker.C:
			if wErr := c.conn.WriteMessage(websocket.PingMessage, nil); wErr != nil {
				errlog.WithoutErr(richerror.New(op).WithWrapError(wErr).WithMessage("error writing ping message"), c.hub.logger)

				return
			}
		}
	}
}

// NewClient created a new Client for websocket connection.
func (s Service) NewClient(ctx context.Context, conn Connection, externalUserID string) *Client {
	const op = "service.service.NewClient"

	userID, err := s.getUserIDFromExternalUserID(ctx, externalUserID)
	if err != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected), s.logger)
	}

	return &Client{
		hub:    s.hub,
		conn:   conn,
		send:   make(chan *NotificationMessage),
		userID: userID,
	}
}
