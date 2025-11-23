package grpc

import (
	"log/slog"

	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/translation"
	"github.com/syntaxfa/quick-connect/protobuf/manager/golang/userinternalpb"
)

type HandlerInternal struct {
	userinternalpb.UnimplementedUserInternalServiceServer

	logger  *slog.Logger
	userSvc userservice.Service
	t       *translation.Translate
}

func NewHandlerInternal(logger *slog.Logger, userSvc userservice.Service, t *translation.Translate) HandlerInternal {
	return HandlerInternal{
		logger:  logger,
		userSvc: userSvc,
		t:       t,
	}
}
