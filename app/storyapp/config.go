package storyapp

import (
	"time"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

type ServiceAuthInfo struct {
	Username string `koanf:"username"`
	Password string `koanf:"password"`
}

type Config struct {
	ShutdownTimeout time.Duration     `koanf:"shutdown_timeout"`
	HTTPServer      httpserver.Config `koanf:"http_server"`
	Logger          logger.Config     `koanf:"logger"`
	Postgres        postgres.Config   `koanf:"postgres"`
	StorageAppGRPC  grpcclient.Config `koanf:"storage_app_grpc"`
	ServiceAuthInfo ServiceAuthInfo   `koanf:"service_auth_info"`
	ManagerAppGRPC  grpcclient.Config `koanf:"manager_app_grpc"`
}
