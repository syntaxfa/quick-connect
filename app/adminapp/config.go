package adminapp

import (
	"time"

	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

type Config struct {
	ShutdownTimeout time.Duration     `koanf:"shutdown_timeout"`
	HTTPServer      httpserver.Config `koanf:"http_server"`
	Logger          logger.Config     `koanf:"logger"`
	TemplatePath    string            `koanf:"template_path"`
	ManagerAppGRPC  grpcclient.Config `koanf:"manager_app_grpc"`
	ChatAppGRPC     grpcclient.Config `koanf:"chat_app_grpc"`
}
