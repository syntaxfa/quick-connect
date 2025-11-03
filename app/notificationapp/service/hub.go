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
				h.mu.Lock()
				if _, ok := h.clients[client.userID]; ok {
					if len(h.clients[client.userID]) > h.cfg.UserConnectionLimit {
						for _, connection := range h.clients[client.userID] {
							close(connection.send)
							if cErr := connection.conn.Close(); cErr != nil {
								errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), h.logger)
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
								errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), h.logger)
							}

							h.clients[client.userID] = append(h.clients[client.userID][:i], h.clients[client.userID][i+1:]...)
							h.logger.Info("client unregistered", slog.String("user_id", string(client.userID)),
								slog.String("addr", client.conn.RemoteAddr()))
						}
					}
				}
				h.mu.Unlock()
			}
		}
	}()

	for {
		message, rErr := subscribe.ReceiveMessage(ctx)
		if rErr != nil {
			errlog.WithoutErrContext(ctx, richerror.New(op).WithMessage("can't receive message").WithWrapError(rErr).WithKind(richerror.KindUnexpected), h.logger)
		}

		var notification NotificationMessage
		if uErr := json.Unmarshal(message, &notification); uErr != nil {
			errlog.WithoutErrContext(ctx, richerror.New(op).WithMessage("can't unmarshalling message").WithWrapError(uErr).WithKind(richerror.KindUnexpected), h.logger)

			continue
		}

		h.mu.RLock()
		if connections, ok := h.clients[notification.UserID]; ok {
			for _, client := range connections {
				select {
				case client.send <- &notification:
					h.logger.DebugContext(ctx, "notification sent to client", slog.String("notification_id",
						string(notification.ID)))
				default:
					h.logger.WarnContext(ctx, "failed to send notification to client, client send buffer full",
						slog.String("user_id", string(notification.UserID)))
				}
			}
		}
		h.mu.RUnlock()
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
