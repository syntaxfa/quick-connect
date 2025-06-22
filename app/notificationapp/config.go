package notificationapp

import (
	"time"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	"github.com/syntaxfa/quick-connect/app/notificationapp/service"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
	"github.com/syntaxfa/quick-connect/pkg/websocket"
)

type Config struct {
	ShutdownTimeout  time.Duration     `koanf:"shutdown_timeout"`
	ClientHTTPServer httpserver.Config `koanf:"client_http_server"`
	AdminHTTPServer  httpserver.Config `koanf:"admin_http_server"`
	Logger           logger.Config     `koanf:"logger"`
	Postgres         postgres.Config   `koanf:"postgres"`
	Notification     service.Config    `koanf:"notification"`
	Redis            redis.Config      `koanf:"redis"`
	Websocket        websocket.Config  `koanf:"websocket"`
	GetUserIDURL     string            `koanf:"get_user_id_url"`
}
