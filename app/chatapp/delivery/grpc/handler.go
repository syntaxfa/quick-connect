package grpc

import (
	"log/slog"

	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/chat/golang/conversationpb"
)

type Handler struct {
	conversationpb.UnimplementedConversationServiceServer

	chatSvc *service.Service
	t       *translation.Translate
	logger  *slog.Logger
}

func NewHandler(chatSvc *service.Service, t *translation.Translate, logger *slog.Logger) Handler {
	return Handler{
		chatSvc: chatSvc,
		t:       t,
		logger:  logger,
	}
}
