package http

import (
	"github.com/gorilla/websocket"
	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"log/slog"
)

type Handler struct {
	upgrader websocket.Upgrader
	logger   *slog.Logger
	svc      service.Service
}

func NewHandler(upgrader websocket.Upgrader, logger *slog.Logger, svc service.Service) Handler {
	return Handler{
		upgrader: upgrader,
		logger:   logger,
		svc:      svc,
	}
}
