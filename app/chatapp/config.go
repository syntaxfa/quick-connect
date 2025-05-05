package chatapp

import (
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
	"time"
)

type Config struct {
	ShutdownTimeout time.Duration     `koanf:"shutdown_timeout"`
	HTTPServer      httpserver.Config `koanf:"http_server"`
	Logger          logger.Config     `koanf:"logger"`
}
