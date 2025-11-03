package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/pubsub"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

type Connection interface {
	WriteMessage(messageType int, data []byte) error
	Close() error
	RemoteAddr() string
}

type Hub struct {
	cfg          Config
	clients      map[types.ID][]*Client
	register     chan *Client
	unregistered chan *Client
	notification chan *NotificationMessage
	logger       *slog.Logger
	mu           sync.RWMutex
	subscriber   pubsub.Subscriber
}

type Client struct {
	hub    *Hub
	conn   Connection
	send   chan *NotificationMessage
	userID types.ID
}

func NewHub(cfg Config, logger *slog.Logger, subscriber pubsub.Subscriber) *Hub {
	return &Hub{
		cfg:          cfg,
		clients:      make(map[types.ID][]*Client),
		register:     make(chan *Client),
		unregistered: make(chan *Client),
		notification: make(chan *NotificationMessage),
		logger:       logger,
		subscriber:   subscriber,
	}
}

func (h *Hub) Run(ctx context.Context) {
	const op = "service.hub.Run"

	subscribe := h.subscriber.Subscribe(ctx, h.cfg.ChannelName)

	go func() {
		for {
			select {
			case client := <-h.register:
				h.handleClientRegister(ctx, client, op)
			case client := <-h.unregistered:
				h.handleClientUnregister(ctx, client, op)
			}
		}
	}()

	for {
		message, rErr := subscribe.ReceiveMessage(ctx)
		if rErr != nil {
			errlog.WithoutErrContext(ctx, richerror.New(op).WithMessage("can't receive message").WithWrapError(rErr).
				WithKind(richerror.KindUnexpected), h.logger)

			continue
		}

		var notification NotificationMessage
		if uErr := json.Unmarshal(message, &notification); uErr != nil {
			errlog.WithoutErrContext(ctx, richerror.New(op).WithMessage("can't unmarshalling message").
				WithWrapError(uErr).WithKind(richerror.KindUnexpected), h.logger)

			continue
		}

		h.sendNotificationToClients(ctx, &notification)
	}
}

func (h *Hub) handleClientRegister(ctx context.Context, client *Client, op string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	connections, exists := h.clients[client.userID]

	if !exists {
		h.clients[client.userID] = []*Client{client}
		h.logger.InfoContext(ctx, "client registered", slog.String("user_id", string(client.userID)),
			slog.String("addr", client.conn.RemoteAddr()))

		return
	}

	if len(connections) > h.cfg.UserConnectionLimit {
		h.closeAllUserConnections(ctx, connections, op)
		h.clients[client.userID] = []*Client{client}
	} else {
		h.clients[client.userID] = append(connections, client)
	}

	h.logger.InfoContext(ctx, "client registered", slog.String("user_id", string(client.userID)),
		slog.String("addr", client.conn.RemoteAddr()))
}

func (h *Hub) closeAllUserConnections(ctx context.Context, connections []*Client, op string) {
	for _, connection := range connections {
		close(connection.send)
		if cErr := connection.conn.Close(); cErr != nil {
			errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), h.logger)
		}
	}
}

func (h *Hub) handleClientUnregister(ctx context.Context, client *Client, op string) {
	h.mu.Lock()
	defer h.mu.Unlock()

	connections, ok := h.clients[client.userID]
	if !ok {
		return
	}

	for i, connection := range connections {
		if connection == client {
			close(client.send)
			if cErr := client.conn.Close(); cErr != nil {
				errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), h.logger)
			}

			h.clients[client.userID] = append(connections[:i], connections[i+1:]...)
			h.logger.InfoContext(ctx, "client unregistered", slog.String("user_id", string(client.userID)),
				slog.String("addr", client.conn.RemoteAddr()))

			break
		}
	}
}

func (h *Hub) sendNotificationToClients(ctx context.Context, notification *NotificationMessage) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	connections, ok := h.clients[notification.UserID]
	if !ok {
		return
	}

	for _, client := range connections {
		select {
		case client.send <- notification:
			h.logger.DebugContext(ctx, "notification sent to client",
				slog.String("notification_id", string(notification.ID)))
		default:
			h.logger.WarnContext(ctx, "failed to send notification to client, client send buffer full",
				slog.String("user_id", string(notification.UserID)))
		}
	}
}

func (c *Client) WritePump() {
	const op = "service.hub.WritePump"

	ticker := time.NewTicker(c.hub.cfg.PingPeriod)
	defer c.cleanupConnection(ticker, op)

	for {
		select {
		case notification, ok := <-c.send:
			if !ok {
				c.sendCloseMessage(op)
				return
			}
			c.sendNotification(notification, op)

		case <-ticker.C:
			if !c.sendPing(op) {
				return
			}
		}
	}
}

func (c *Client) cleanupConnection(ticker *time.Ticker, op string) {
	ticker.Stop()
	if cErr := c.conn.Close(); cErr != nil {
		errlog.WithoutErr(
			richerror.New(op).
				WithWrapError(cErr).
				WithKind(richerror.KindUnexpected).
				WithMessage("failed to stop client websocket connection"),
			c.hub.logger,
		)
	}
	c.hub.unregistered <- c
}

func (c *Client) sendCloseMessage(op string) {
	if wErr := c.conn.WriteMessage(websocket.CloseMessage, []byte{}); wErr != nil {
		errlog.WithoutErr(
			richerror.New(op).
				WithMessage("error sending websocket close message").
				WithKind(richerror.KindUnexpected),
			c.hub.logger,
		)
	}
}

func (c *Client) sendNotification(notification *NotificationMessage, op string) {
	msgBytes, mErr := json.Marshal(notification)
	if mErr != nil {
		errlog.WithoutErr(
			richerror.New(op).
				WithWrapError(mErr).
				WithMessage("error marshalling notification"),
			c.hub.logger,
		)
		return
	}

	if wErr := c.conn.WriteMessage(websocket.TextMessage, msgBytes); wErr != nil {
		errlog.WithoutErr(
			richerror.New(op).
				WithWrapError(wErr).
				WithMessage("error writing message"),
			c.hub.logger,
		)
	}
}

func (c *Client) sendPing(op string) bool {
	if wErr := c.conn.WriteMessage(websocket.PingMessage, nil); wErr != nil {
		errlog.WithoutErr(
			richerror.New(op).
				WithWrapError(wErr).
				WithMessage("error writing ping message"),
			c.hub.logger,
		)
		return false
	}
	return true
}

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
