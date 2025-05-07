package service

import (
	"log/slog"

	"github.com/google/uuid"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

type Connection interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
}

type Repository interface {
	SaveMessage(message Message) error
}

type Service struct {
	cfg            Config
	hub            *Hub
	repo           Repository
	websocketConns map[string]Connection
	logger         *slog.Logger
}

func New(cfg Config, repo Repository, logger *slog.Logger) *Service {
	hub := NewHub(logger)
	go hub.Run()

	return &Service{
		cfg:            cfg,
		hub:            hub,
		repo:           repo,
		websocketConns: make(map[string]Connection),
		logger:         logger,
	}
}

func (s *Service) JoinClient(conn Connection, username string) {
	userID := uuid.New().String()

	client := NewClient(s.cfg, userID, s.hub, conn, username, s.logger)

	s.websocketConns[userID] = conn

	s.hub.RegisterClient(client)

	go client.ReadPump()
	go client.WritePump()
}

func (s *Service) JoinSupport(conn Connection, username string) {
	userID := uuid.New().String()

	support := NewSupport(s.cfg, userID, s.hub, conn, username, s.logger)

	s.websocketConns[userID] = conn

	s.hub.RegisterSupport(support)

	go support.ReadPump()
	go support.WritePump()
}

func (s *Service) ClientSendMessage(message Message) {
	// Save message to repository

	s.hub.BroadcastMessageToSupport(message)
}

func (s *Service) SupportSendMessage(message Message) {
	// Save message to repository

	s.hub.BroadcastMessageToClient(message)
}

func (s *Service) CloseConnection(userID string) {
	const op = "service.service.CloseConnection"

	if conn, exists := s.websocketConns[userID]; exists {
		if cErr := conn.Close(); cErr != nil {
			errlog.ErrLog(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), s.logger)
		}
		delete(s.websocketConns, userID)
	}
}
