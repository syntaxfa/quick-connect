package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"

	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/pubsub"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

// Hub manages active websocket connections on this specific server instance.
type Hub struct {
	cfg          Config
	participants map[types.ID][]*Participant // Map of UserID to their active connections
	register     chan *Participant
	unregister   chan *Participant
	logger       *slog.Logger
	mu           sync.RWMutex
	subscriber   pubsub.Subscriber
}

func NewHub(cfg Config, logger *slog.Logger, subscriber pubsub.Subscriber) *Hub {
	return &Hub{
		cfg:          cfg,
		participants: make(map[types.ID][]*Participant),
		register:     make(chan *Participant),
		unregister:   make(chan *Participant),
		logger:       logger,
		subscriber:   subscriber,
	}
}

// Register registers a new participant connection with the hub.
func (h *Hub) Register(p *Participant) {
	h.register <- p
}

// Unregister removes a participant connection from the hub.
func (h *Hub) Unregister(p *Participant) {
	h.unregister <- p
}

// Run starts the hub's connection management and Pub/Sub listener.
func (h *Hub) Run(ctx context.Context) {
	const op = "service.hub.Run"

	if h.cfg.ChatChannelName == "" {
		h.logger.ErrorContext(ctx, "chat channel name is not configured", slog.String("op", op))
		return
	}

	receiver := h.subscriber.Subscribe(ctx, h.cfg.ChatChannelName)
	h.logger.InfoContext(ctx, "hub started and subscribed to pubsub", slog.String("channel", h.cfg.ChatChannelName))

	// Run the registration manager in a separate goroutine
	go h.manageRegistrations(ctx)

	// Main goroutine: Listen for messages from Pub/Sub
	for {
		select {
		case <-ctx.Done():
			h.logger.InfoContext(ctx, "context done, stopping hub run loop")
			return
		default:
			message, rErr := receiver.ReceiveMessage(ctx)
			if rErr != nil {
				if ctx.Err() == nil { // Don't log errors if context was just cancelled
					errlog.WithoutErrContext(ctx, richerror.New(op).WithMessage("can't receive message").WithWrapError(rErr).
						WithKind(richerror.KindUnexpected), h.logger)
				}
				continue
			}

			h.processPubSubMessage(ctx, message)
		}
	}
}

// processPubSubMessage unmarshals the pubsub message and fans it out to local clients.
func (h *Hub) processPubSubMessage(ctx context.Context, payload []byte) {
	const op = "service.hub.processPubSubMessage"

	var pubsubPayload PubSubMessage
	if err := json.Unmarshal(payload, &pubsubPayload); err != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindInvalid), h.logger)
		return
	}

	// Re-marshal the inner ServerMessage to send to clients
	serverMsgPayload, err := json.Marshal(pubsubPayload.ServerMsg)
	if err != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected), h.logger)
		return
	}

	// Create a map for efficient exclusion
	excludeMap := make(map[types.ID]struct{}, len(pubsubPayload.ExcludeIDs))
	for _, id := range pubsubPayload.ExcludeIDs {
		excludeMap[id] = struct{}{}
	}

	for _, userID := range pubsubPayload.RecipientIDs {
		if _, excluded := excludeMap[userID]; excluded {
			continue // Skip excluded users
		}

		h.sendMessageToUser(ctx, userID, serverMsgPayload)
	}
}

// manageRegistrations handles register and unregister events.
func (h *Hub) manageRegistrations(ctx context.Context) {
	for {
		select {
		case client := <-h.register:
			h.handleRegister(ctx, client)
		case client := <-h.unregister:
			h.handleUnregister(ctx, client)
		case <-ctx.Done():
			h.logger.InfoContext(ctx, "context done, stopping registration manager")
			return
		}
	}
}

func (h *Hub) handleRegister(ctx context.Context, client *Participant) {
	h.mu.Lock()
	defer h.mu.Unlock()

	connections, exists := h.participants[client.userID]
	if !exists {
		h.participants[client.userID] = []*Participant{client}
	} else {
		// TODO: Add connection limiting logic (e.g., cfg.UserConnectionLimit)
		h.participants[client.userID] = append(connections, client)
	}
	h.logger.InfoContext(ctx, "participant registered", slog.String("user_id", string(client.userID)))
}

func (h *Hub) handleUnregister(ctx context.Context, client *Participant) {
	h.mu.Lock()
	defer h.mu.Unlock()

	connections, ok := h.participants[client.userID]
	if !ok {
		return
	}

	for i, connection := range connections {
		if connection == client {
			h.participants[client.userID] = append(connections[:i], connections[i+1:]...)
			if len(h.participants[client.userID]) == 0 {
				delete(h.participants, client.userID)
			}
			close(client.send)
			h.logger.InfoContext(ctx, "participant unregistered", slog.String("user_id", string(client.userID)))
			break
		}
	}
}

// sendMessageToUser sends a payload to all active connections for a specific user *on this server*.
func (h *Hub) sendMessageToUser(ctx context.Context, userID types.ID, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	connections, ok := h.participants[userID]
	if !ok {
		h.logger.DebugContext(ctx, "no active connections found for user to send message", slog.String("user_id", string(userID)))
		return
	}

	for _, client := range connections {
		select {
		case client.send <- payload:
		default:
			h.logger.WarnContext(ctx, "failed to send message to participant, send buffer full", slog.String("user_id", string(userID)))
		}
	}
}
