package service

import (
	"encoding/json"
	"log/slog"
	"time"

	"github.com/gorilla/websocket"
	"github.com/syntaxfa/quick-connect/pkg/errlog"
	"github.com/syntaxfa/quick-connect/pkg/richerror"
)

type wsConnection struct {
	conn   Connection
	send   chan Message
	cfg    Config
	logger *slog.Logger
}

func newWsConnection(cfg Config, conn Connection, send chan Message, logger *slog.Logger) *wsConnection {
	return &wsConnection{
		conn:   conn,
		send:   send,
		cfg:    cfg,
		logger: logger,
	}
}

func (w *wsConnection) writePump(op string) {
	ticker := time.NewTicker(w.cfg.PingPeriod)
	defer func() {
		ticker.Stop()
		closeConnection(w.conn, w.logger, op)
	}()

	for {
		select {
		case message, ok := <-w.send:
			if !ok {
				sendCloseMessage(w.conn, w.logger, op)
				return
			}
			sendMessage(w.conn, message, w.logger, op)

		case <-ticker.C:
			if !sendPing(w.conn) {
				return
			}
		}
	}
}

func closeConnection(conn Connection, logger *slog.Logger, op string) {
	if cErr := conn.Close(); cErr != nil {
		errlog.WithoutErr(richerror.New(op).WithWrapError(cErr).WithKind(richerror.KindUnexpected), logger)
	}
}

func sendCloseMessage(conn Connection, logger *slog.Logger, op string) {
	if wErr := conn.WriteMessage(websocket.CloseMessage, []byte{}); wErr != nil {
		errlog.WithoutErr(
			richerror.New(op).
				WithMessage("error sending websocket close message").
				WithKind(richerror.KindUnexpected),
			logger,
		)
	}
}

func sendMessage(conn Connection, message Message, logger *slog.Logger, op string) bool {
	msgBytes, mErr := json.Marshal(message)
	if mErr != nil {
		errlog.WithoutErr(
			richerror.New(op).WithWrapError(mErr).WithMessage("error marshalling message"),
			logger,
		)
		return false
	}

	if wErr := conn.WriteMessage(websocket.TextMessage, msgBytes); wErr != nil {
		errlog.WithoutErr(
			richerror.New(op).WithWrapError(wErr).WithMessage("error writing message"),
			logger,
		)
		return false
	}

	return true
}

func sendPing(conn Connection) bool {
	return conn.WriteMessage(websocket.PingMessage, nil) == nil
}
