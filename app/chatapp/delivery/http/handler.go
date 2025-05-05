package http

import (
	"github.com/gorilla/websocket"
	"log/slog"
)

type Handler struct {
	upgrader websocket.Upgrader
	logger   *slog.Logger
}

func NewHandler(upgrader websocket.Upgrader, logger *slog.Logger) Handler {
	return Handler{
		upgrader: upgrader,
		logger:   logger,
	}
}
