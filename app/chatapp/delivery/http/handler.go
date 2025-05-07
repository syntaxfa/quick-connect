package http

import (
	"log/slog"

	"github.com/syntaxfa/quick-connect/adapter/websocket"
	"github.com/syntaxfa/quick-connect/app/chatapp/service"
)

type Handler struct {
	upgrader *websocket.GorillaUpgrader
	logger   *slog.Logger
	svc      *service.Service
}

func NewHandler(upgrader *websocket.GorillaUpgrader, logger *slog.Logger, svc *service.Service) Handler {
	return Handler{
		upgrader: upgrader,
		logger:   logger,
		svc:      svc,
	}
}
