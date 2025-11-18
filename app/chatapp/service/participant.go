package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
	"github.com/syntaxfa/quick-connect/types"
)

const (
	SendChanSize = 256
)

// Connection defines the interface for a websocket connection.
type Connection interface {
	ReadMessage() (messageType int, p []byte, err error)
	WriteMessage(messageType int, data []byte) error
	Close() error
	SetReadLimit(limit int64)
	SetReadDeadline(t time.Time) error
	SetWriteDeadline(t time.Time) error
	SetPongHandler(h func(appData string) error)
	SetCloseHandler(h func(code int, text string) error)
}

// MessageHandlerCallback is a function callback executed by readPump for processing.
type MessageHandlerCallback func(ctx context.Context, senderID types.ID, messageType int, payload []byte)

// Participant represents a single active websocket connection.
type Participant struct {
	userID types.ID
	hub    *Hub
	conn   Connection
	send   chan []byte // Buffered channel for outbound messages (raw bytes)
	logger *slog.Logger
	cfg    Config
}

// NewParticipant creates a new participant manager for a websocket connection.
func NewParticipant(cfg Config, hub *Hub, conn Connection, userID types.ID, logger *slog.Logger) *Participant {
	return &Participant{
		userID: userID,
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, SendChanSize),
		logger: logger.With(slog.String("user_id", string(userID))),
		cfg:    cfg,
	}
}

// ReadPump pumps messages from the websocket connection to the Service handler.
func (p *Participant) ReadPump(ctx context.Context, handler MessageHandlerCallback) {
	const op = "service.participant.ReadPump"
	defer func() {
		p.hub.Unregister(p)
		if err := p.conn.Close(); err != nil {
			errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected), p.logger)
		}
		p.logger.DebugContext(ctx, "read pump stopped")
	}()

	p.conn.SetReadLimit(int64(p.cfg.MaxMessageSize))
	if err := p.conn.SetReadDeadline(time.Now().Add(p.cfg.PongWait)); err != nil {
		errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected), p.logger)
		return
	}
	p.conn.SetPongHandler(func(string) error {
		p.logger.DebugContext(ctx, "pong received")
		return p.conn.SetReadDeadline(time.Now().Add(p.cfg.PongWait))
	})
	p.conn.SetCloseHandler(func(code int, text string) error {
		p.logger.InfoContext(ctx, "websocket connection closed by client", slog.Int("code", code), slog.String("text", text))
		return nil
	})

	for {
		messageType, payload, err := p.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected), p.logger)
			}
			break
		}

		handler(ctx, p.userID, messageType, payload)
	}
}

// WritePump pumps messages from the hub to the websocket connection.
func (p *Participant) WritePump(ctx context.Context) {
	const op = "service.participant.WritePump"
	ticker := time.NewTicker(p.cfg.PingPeriod)
	defer func() {
		ticker.Stop()
		if err := p.conn.Close(); err != nil {
			errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected), p.logger)
		}
		p.logger.DebugContext(ctx, "write pump stopped")
	}()

	for {
		select {
		case payload, ok := <-p.send:
			if !ok {
				p.sendCloseMessage(op)
				return
			}

			if err := p.conn.SetWriteDeadline(time.Now().Add(p.cfg.WriteWait)); err != nil {
				errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected), p.logger)
				return
			}
			if err := p.conn.WriteMessage(websocket.TextMessage, payload); err != nil {
				errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected), p.logger)
				return
			}

		case <-ticker.C:
			if err := p.conn.SetWriteDeadline(time.Now().Add(p.cfg.WriteWait)); err != nil {
				errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected), p.logger)
				return
			}
			p.logger.DebugContext(ctx, "sending ping")
			if err := p.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				errlog.WithoutErrContext(ctx, richerror.New(op).WithWrapError(err).WithKind(richerror.KindUnexpected), p.logger)
				return
			}
		case <-ctx.Done():
			p.logger.InfoContext(ctx, "context done, stopping write pump")
			return
		}
	}
}

// sendCloseMessage writes a WebSocket close message.
func (p *Participant) sendCloseMessage(op string) {
	if err := p.conn.WriteMessage(websocket.CloseMessage, []byte{}); err != nil {
		errlog.WithoutErr(
			richerror.New(op).
				WithMessage("error sending websocket close message").
				WithKind(richerror.KindUnexpected),
			p.logger,
		)
	}
}
