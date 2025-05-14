package http

import (
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type Handler struct {
	t        *translation.Translate
	tokenSvc tokenservice.Service
	userSvc  userservice.Service
}

func NewHandler(t *translation.Translate, tokenSvc tokenservice.Service, userSvc userservice.Service) Handler {
	return Handler{
		t:        t,
		tokenSvc: tokenSvc,
		userSvc:  userSvc,
	}
}
