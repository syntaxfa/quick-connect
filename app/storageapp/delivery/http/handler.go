package http

import (
	"log/slog"

	"github.com/syntaxfa/quick-connect/app/storageapp/service"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type Handler struct {
	svc           service.Service
	t             *translation.Translate
	localRootPath string
	maxSize       int64
	logger        *slog.Logger
}

func NewHandler(svc service.Service, t *translation.Translate, localRootPath string, maxSize int64, logger *slog.Logger) Handler {
	return Handler{
		svc:           svc,
		t:             t,
		localRootPath: localRootPath,
		maxSize:       maxSize,
		logger:        logger,
	}
}
