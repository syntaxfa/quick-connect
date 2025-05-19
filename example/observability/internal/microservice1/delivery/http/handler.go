package http

import "github.com/syntaxfa/quick-connect/example/observability/internal/microservice1/service"

type Handler struct {
	svc service.Service
}

func NewHandler(svc service.Service) Handler {
	return Handler{
		svc: svc,
	}
}
