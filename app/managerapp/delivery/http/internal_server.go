package http

import (
	"context"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/servermsg"
	"github.com/syntaxfa/quick-connect/pkg/translation"
)

type InternalServer struct {
	httpServer httpserver.Server
	handler    InternalHandler
	logger     *slog.Logger
	apiKey     string
}

func NewInternalServer(httpServer httpserver.Server, handler InternalHandler, logger *slog.Logger, apiKey string) InternalServer {
	return InternalServer{
		httpServer: httpServer,
		handler:    handler,
		logger:     logger,
		apiKey:     apiKey,
	}
}

type InternalHandler struct {
	userSvc userservice.Service
	t       *translation.Translate
}

func NewInternalHandler(userSvc userservice.Service, t *translation.Translate) InternalHandler {
	return InternalHandler{
		userSvc: userSvc,
		t:       t,
	}
}

func (s InternalServer) Start() error {
	s.registerRoutes()

	return s.httpServer.Start()
}

func (s InternalServer) Stop(ctx context.Context) error {
	return s.httpServer.Stop(ctx)
}

func (s InternalServer) registerRoutes() {
	rootGr := s.httpServer.Router.Group("", AccessToPrivateEndpoint(s.apiKey))

	rootGr.POST("/auth/identify-client", s.handler.IdentifyClient)
}

// IdentifyClient docs
// @Summary IdentifyClient
// @Description register guest user and generate QCToken (long expire time)
// @Tags Internal
// @Accept json
// @Produce json
// @Param APIKey query string true "API Key for authorization"
// @Param Request body userservice.IdentifyClientRequest true "check token validation"
// @Success 201 {object} userservice.IdentifyClientResponse
// @Failure 400 {string} string Bad Request
// @Failure 422 {object} servermsg.ErrorResponse
// @Failure 500 {string} something went wrong
// @Router /auth/identify-client [POST].
func (h InternalHandler) IdentifyClient(c echo.Context) error {
	var req userservice.IdentifyClientRequest
	if bErr := c.Bind(&req); bErr != nil {
		return echo.NewHTTPError(http.StatusBadRequest)
	}

	resp, sErr := h.userSvc.IdentifyClient(c.Request().Context(), req)
	if sErr != nil {
		return servermsg.HTTPMsg(c, sErr, h.t)
	}

	return c.JSON(http.StatusOK, resp)
}
