package service

import (
	"context"
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/paginate/cursorbased"
	paginate "github.com/syntaxfa/quick-connect/pkg/paginate/limitoffset"
	"github.com/syntaxfa/quick-connect/pkg/pubsub"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

// Repository defines the persistence interface required by the chat service.
// It includes methods needed by both conversation.go and service.go.
type Repository interface {
	IsUserHaveActiveConversation(ctx context.Context, userID types.ID) (bool, error)
	CreateActiveConversation(ctx context.Context, id, userID types.ID, conversationStatus ConversationStatus) (Conversation, error)
	GetUserActiveConversation(ctx context.Context, userID types.ID) (Conversation, error)
	GetConversationList(ctx context.Context, paginated paginate.RequestBase, assignedSupportID types.ID,
		statuses []ConversationStatus) ([]Conversation, paginate.ResponseBase, error)

	GetConversationByID(ctx context.Context, conversationID types.ID) (Conversation, error)
	GetConversationParticipants(ctx context.Context, conversationID types.ID) ([]types.ID, error)
	CheckUserInConversation(ctx context.Context, userID, conversationID types.ID) (bool, error)
	SaveMessage(ctx context.Context, message Message) (Message, error)
	UpdateConversationSnippet(ctx context.Context, conversationID, lastMessageSenderID types.ID, snippet string) error
	IsConversationExistByID(ctx context.Context, conversationID types.ID) (bool, error)
	AssignConversation(ctx context.Context, conversationID, supportID types.ID) error

	GetChatHistory(ctx context.Context, conversationID types.ID, req cursorbased.Request) (ChatHistoryResponse, error)
}

// Service handles the business logic for the real-time chat and conversations.
// Methods for this struct are defined in service.go, conversation.go, etc.
type Service struct {
	cfg       Config
	hub       *Hub
	repo      Repository
	logger    *slog.Logger
	vld       Validate
	publisher pubsub.Publisher
}

// New creates a new chat Service.
func New(cfg Config, repo Repository, hub *Hub, publisher pubsub.Publisher, logger *slog.Logger, vld Validate) *Service {
	return &Service{
		cfg:       cfg,
		hub:       hub,
		repo:      repo,
		publisher: publisher,
		logger:    logger,
		vld:       vld,
	}
}

// HandleNewConnection is the entry point for a new websocket connection.
func (s *Service) HandleNewConnection(ctx context.Context, conn Connection, userID types.ID) {
	participant := NewParticipant(s.cfg, s.hub, conn, userID, s.logger)
	s.hub.Register(participant)

	go participant.WritePump(ctx)
	go participant.ReadPump(ctx, s.HandleIncomingMessage)

	s.logger.InfoContext(ctx, "new participant connection handled", slog.String("user_id", string(userID)))
}

// HandleIncomingMessage is the callback for participant.readPump.
func (s *Service) HandleIncomingMessage(ctx context.Context, senderID types.ID, messageType int, payload []byte) {
	const op = "service.Service.HandleIncomingMessage"

	if messageType != websocket.TextMessage {
		s.logger.WarnContext(ctx, "received non-text message", slog.String("op", op), slog.String("user_id", string(senderID)))

		return
	}

	var clientMsg ClientMessage
	if err := json.Unmarshal(payload, &clientMsg); err != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindInvalid), s.logger)

		return
	}

	clientMsg.SenderID = senderID

	if err := s.vld.ValidateClientMessage(clientMsg); err != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindInvalid), s.logger)

		return
	}

	allowed, err := s.repo.CheckUserInConversation(ctx, senderID, clientMsg.ConversationID)
	if err != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindForbidden), s.logger)

		return
	}
	if !allowed {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithMessage("user not allowed in conversation").
			WithKind(richerror.KindForbidden), s.logger)

		return
	}

	conv, err := s.repo.GetConversationByID(ctx, clientMsg.ConversationID)
	if err != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindNotFound), s.logger)

		return
	}
	if conv.Status == ConversationStatusClosed {
		s.logger.WarnContext(ctx, "message sent to closed conversation", slog.String("op", op),
			slog.String("senderID", string(senderID)), slog.String("conversationID", string(clientMsg.ConversationID)))

		// We stop processing here. No message is sent or saved.
		return
	}

	switch clientMsg.Type {
	case MessageTypeText:
		s.handleTextMessage(ctx, clientMsg)
	case MessageTypeSystem:
		s.handleSystemMessage(ctx, clientMsg)
	case MessageTypeMedia:
		// TODO: complete this!!!
		return
	default:
		s.logger.WarnContext(ctx, "unknown client message type", slog.String("op", op), slog.String("type", string(clientMsg.Type)))
	}
}

// handleTextMessage processes a text message: save to DB, update snippet, and publish.
func (s *Service) handleTextMessage(ctx context.Context, msg ClientMessage) {
	const op = "service.Service.handleTextMessage"

	dbMsg := Message{
		ID:             types.ID(ulid.Make().String()),
		ConversationID: msg.ConversationID,
		SenderID:       msg.SenderID,
		MessageType:    MessageTypeText,
		Content:        msg.Content,
		CreatedAt:      time.Now(),
	}

	savedMsg, sErr := s.repo.SaveMessage(ctx, dbMsg)
	if sErr != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(sErr).WithKind(richerror.KindUnexpected), s.logger)

		return
	}

	snippet := msg.Content
	runes := []rune(snippet)
	if len(runes) > s.cfg.MessageSnippetCharNumber {
		snippet = string(runes[:50]) + "..."
	}

	if uErr := s.repo.UpdateConversationSnippet(ctx, msg.ConversationID, msg.SenderID, snippet); uErr != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(uErr).WithKind(richerror.KindUnexpected), s.logger)
	}

	wsMsg := NewTextMessage(savedMsg, msg.ClientMessageID)

	// Publish to all conversation participants (including the sender)
	s.publishToConversation(ctx, msg.ConversationID, wsMsg)
}

// handleSystemMessage processes a system message: publish to other participants.
func (s *Service) handleSystemMessage(ctx context.Context, msg ClientMessage) {
	const op = "service.Service.handleSystemMessage"

	var wsMsg ServerMessage
	switch msg.SubType {
	case "typing_started":
		wsMsg = NewSystemMessage(msg.ConversationID, msg.SenderID, "typing_started")
	case "typing_stopped":
		wsMsg = NewSystemMessage(msg.ConversationID, msg.SenderID, "typing_stopped")
	default:
		s.logger.WarnContext(ctx, "unknown system message sub_type", slog.String("op", op), slog.String("sub_type", msg.SubType))
		return
	}

	// Publish system messages to all participants *except* the original sender
	s.publishToConversation(ctx, msg.ConversationID, wsMsg, msg.SenderID)
}

// publishToConversation finds all participants and publishes one message to the global chat topic.
func (s *Service) publishToConversation(ctx context.Context, conversationID types.ID, message ServerMessage, excludeIDs ...types.ID) {
	const op = "service.Service.publishToConversation"

	participants, gErr := s.repo.GetConversationParticipants(ctx, conversationID)
	if gErr != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(gErr).WithKind(richerror.KindNotFound), s.logger)
		return
	}

	// Wrap the message in a PubSubMessage
	pubsubPayload := PubSubMessage{
		RecipientIDs: participants,
		ServerMsg:    message,
		ExcludeIDs:   excludeIDs,
	}

	payload, mErr := json.Marshal(pubsubPayload)
	if mErr != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(mErr).WithKind(richerror.KindUnexpected), s.logger)
		return
	}

	// Publish to the single, global chat channel
	if pErr := s.publisher.Publish(ctx, s.cfg.ChatChannelName, payload); pErr != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(pErr).
			WithMeta(map[string]interface{}{"channel": s.cfg.ChatChannelName}),
			s.logger)
	}
}
