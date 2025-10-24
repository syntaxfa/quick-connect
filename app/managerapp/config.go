package managerapp

import (
	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"time"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

type Config struct {
	ShutdownTimeout time.Duration       `koanf:"shutdown_timeout"`
	HTTPServer      httpserver.Config   `koanf:"http_server"`
	Logger          logger.Config       `koanf:"logger"`
	Token           tokenservice.Config `koanf:"token"`
	Postgres        postgres.Config     `koanf:"postgres"`
	GRPCServer      grpcserver.Config   `koanf:"grpc_server"`
}
