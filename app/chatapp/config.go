package chatapp

import (
	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"time"

	"github.com/syntaxfa/quick-connect/adapter/websocket"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

type Config struct {
	ShutdownTimeout time.Duration     `koanf:"shutdown_timeout"`
	HTTPServer      httpserver.Config `koanf:"http_server"`
	Logger          logger.Config     `koanf:"logger"`
	Websocket       websocket.Config  `koanf:"websocket"`
	ChatService     service.Config    `koanf:"chat"`
}
