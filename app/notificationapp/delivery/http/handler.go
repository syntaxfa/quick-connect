package http

import (
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/pkg/websocket"
)

type Handler struct {
	svc      service.Service
	t        *translation.Translate
	upgrader *websocket.GorillaUpgrader
}

func NewHandler(svc service.Service, t *translation.Translate, upgrader *websocket.GorillaUpgrader) Handler {
	return Handler{
		svc:      svc,
		t:        t,
		upgrader: upgrader,
	}
}
