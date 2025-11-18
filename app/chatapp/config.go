package chatapp

import (
	"time"

	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	"github.com/syntaxfa/quick-connect/app/chatapp/service"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
	"github.com/syntaxfa/quick-connect/pkg/websocket"
)

type Config struct {
	ShutdownTimeout time.Duration     `koanf:"shutdown_timeout"`
	HTTPServer      httpserver.Config `koanf:"http_server"`
	Logger          logger.Config     `koanf:"logger"`
	Websocket       websocket.Config  `koanf:"websocket"`
	ChatService     service.Config    `koanf:"chat"`
	Postgres        postgres.Config   `koanf:"postgres"`
	ManagerAppGRPC  grpcclient.Config `koanf:"manager_app_grpc"`
	GRPCServer      grpcserver.Config `koanf:"grpc_server"`
	Redis           redis.Config      `koanf:"redis"`
}
