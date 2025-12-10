package http

import (
	"log/slog"

	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type Handler struct {
	t      *translation.Translate
	logger *slog.Logger
}

func NewHandler(t *translation.Translate, logger *slog.Logger) Handler {
	return Handler{
		t:      t,
		logger: logger,
	}
}
