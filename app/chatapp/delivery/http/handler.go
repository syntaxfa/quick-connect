package http

import (
	"log/slog"

	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/pkg/websocket"
)

type Handler struct {
	upgrader *websocket.GorillaUpgrader
	logger   *slog.Logger
	svc      *service.Service
	t        *translation.Translate
}

func NewHandler(upgrader *websocket.GorillaUpgrader, logger *slog.Logger, svc *service.Service, t *translation.Translate) Handler {
	return Handler{
		upgrader: upgrader,
		logger:   logger,
		svc:      svc,
		t:        t,
	}
}
