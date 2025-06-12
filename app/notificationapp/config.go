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
	ShutdownTimeout time.Duration     `koanf:"shutdown_timeout"`
	HTTPServer      httpserver.Config `koanf:"http_server"`
	Logger          logger.Config     `koanf:"logger"`
	Postgres        postgres.Config   `koanf:"postgres"`
	Notification    service.Config    `koanf:"notification"`
	Redis           redis.Config      `koanf:"redis"`
	Hub             service.HubConfig `koanf:"hub"`
	Websocket       websocket.Config  `koanf:"websocket"`
}
