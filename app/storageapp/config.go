package storageapp

import (
	"time"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/storage/aws"
	"github.com/syntaxfa/quick-connect/adapter/storage/local"
	"github.com/syntaxfa/quick-connect/app/storageapp/service"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

type Storage struct {
	Driver service.Driver `koanf:"driver"`
	AWS    aws.Config     `koanf:"aws"`
	Local  local.Config   `koanf:"local"`
}

type Config struct {
	ShutdownTimeout time.Duration     `koanf:"shutdown_timeout"`
	HTTPServer      httpserver.Config `koanf:"http_server"`
	Logger          logger.Config     `koanf:"logger"`
	Storage         Storage           `koanf:"storage"`
	Postgres        postgres.Config   `koanf:"postgres"`
	Service         service.Config    `koanf:"service"`
	ManagerAppGRPC  grpcclient.Config `koanf:"manager_app_grpc"`
}
