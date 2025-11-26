package http

import (
	"context"
	"log/slog"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syntaxfa/quick-connect/app/notificationapp/docs"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
)

type AdminServer struct {
	httpServer httpserver.Server
	handler    Handler
	logger     *slog.Logger
}

func NewAdminServer(httpServer httpserver.Server, handler Handler, logger *slog.Logger) AdminServer {
	return AdminServer{
		httpServer: httpServer,
		handler:    handler,
		logger:     logger,
	}
}

func (s AdminServer) Start() error {
	s.registerRoutes()

	return s.httpServer.Start()
}

func (s AdminServer) Stop(ctx context.Context) error {
	return s.httpServer.Stop(ctx)
}

func (s AdminServer) registerRoutes() {
	s.registerSwagger()

	s.httpServer.Router.GET("/health-check", s.handler.healthCheck)

	v1 := s.httpServer.Router.Group("/v1")

	notifications := v1.Group("/notifications")
	notifications.POST("", s.handler.sendNotification)

	templates := v1.Group("/templates")
	templates.POST("", s.handler.createTemplate)
	templates.POST("/list", s.handler.ListTemplate)
	templates.PUT("/:templateID", s.handler.updateTemplate)
	templates.GET("/:templateID", s.handler.getDetailTemplate)

	settings := v1.Group("/settings")
	settings.POST("/:externalUserID", s.handler.updateUserSettingAdmin)
	settings.GET("/:externalUserID", s.handler.getUserSettingAdmin)
}

func (s AdminServer) registerSwagger() {
	docs.SwaggerInfonotification.Title = "Notification Admin API"
	docs.SwaggerInfonotification.Description = "Notification admin restfull API documentation"
	docs.SwaggerInfonotification.Version = "1.0.0"

	s.httpServer.Router.GET("/swagger/*any", echoSwagger.EchoWrapHandler(echoSwagger.InstanceName("notification")))
}
