package http

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	echoSwagger "github.com/swaggo/echo-swagger"
	"github.com/syntaxfa/quick-connect/app/notificationapp/docs"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
)

type Server struct {
	httpServer        httpserver.Server
	handler           Handler
	getExternalUserID string
	logger            *slog.Logger
}

func New(httpServer httpserver.Server, handler Handler, getExternalUserID string, logger *slog.Logger) Server {
	return Server{
		httpServer:        httpServer,
		handler:           handler,
		getExternalUserID: getExternalUserID,
		logger:            logger,
	}
}

func (s Server) Start() error {
	s.registerRoutes()

	return s.httpServer.Start()
}

func (s Server) Stop(ctx context.Context) error {
	return s.httpServer.Stop(ctx)
}

func (s Server) registerRoutes() {
	s.registerSwagger()

	s.httpServer.Router.GET("/health-check", s.handler.healthCheck)

	v1 := s.httpServer.Router.Group("/v1")

	httpClient := &http.Client{Timeout: time.Second * 10}

	notifications := v1.Group("/notifications")
	notifications.POST("", s.handler.sendNotification)
	notifications.POST("/list", s.handler.findNotifications)
	notifications.GET("/:notificationID/mark-as-read", s.handler.markNotificationAsRead)
	notifications.GET("/:externalUserID/mark-all-as-read", s.handler.markAllNotificationAsRead)

	notifications.GET("/ws", s.handler.wsNotification, validateExternalToken(s.getExternalUserID, s.logger, httpClient))
}

func (s Server) registerSwagger() {
	docs.SwaggerInfo.Title = "Notification API"
	docs.SwaggerInfo.Description = "Notification restfull API documentation"
	docs.SwaggerInfo.Version = "1.0.0"

	s.httpServer.Router.GET("/swagger/*any", echoSwagger.WrapHandler)
}
