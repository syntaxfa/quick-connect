package http

import (
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type Handler struct {
	svc service.Service
	t   *translation.Translate
}

func NewHandler(svc service.Service, t *translation.Translate) Handler {
	return Handler{
		svc: svc,
		t:   t,
	}
}
