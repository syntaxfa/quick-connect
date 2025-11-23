package managerapp

import (
	"time"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	"github.com/syntaxfa/quick-connect/app/managerapp/delivery/http"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/tokenservice"
	"github.com/syntaxfa/quick-connect/app/managerapp/service/userservice"
	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

type Config struct {
	ShutdownTimeout    time.Duration       `koanf:"shutdown_timeout"`
	HTTPServer         httpserver.Config   `koanf:"http_server"`
	Logger             logger.Config       `koanf:"logger"`
	Token              tokenservice.Config `koanf:"token"`
	Postgres           postgres.Config     `koanf:"postgres"`
	GRPCServer         grpcserver.Config   `koanf:"grpc_server"`
	User               userservice.Config  `koanf:"user"`
	Redis              redis.Config        `koanf:"redis"`
	Delivery           http.Config         `koanf:"delivery"`
	InternalHTTPServer httpserver.Config   `koanf:"internal_http_server"`
	APIKey             string              `koanf:"api_key"`
	GRPCServerInternal grpcserver.Config   `koanf:"grpc_server_internal"`
}
