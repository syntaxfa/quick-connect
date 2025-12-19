package grpc

import (
	"log/slog"

	"github.com/syntaxfa/quick-connect/app/storageapp/service"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/storage/golang/storagepb"
)

type InternalHandler struct {
	storagepb.UnimplementedStorageInternalServiceServer

	svc    service.Service
	t      *translation.Translate
	logger *slog.Logger
}

func NewInternalHandler(svc service.Service, t *translation.Translate, logger *slog.Logger) InternalHandler {
	return InternalHandler{
		svc:    svc,
		t:      t,
		logger: logger,
	}
}
