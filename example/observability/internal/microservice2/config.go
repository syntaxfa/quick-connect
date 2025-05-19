package microservice2

import (
	"github.com/syntaxfa/quick-connect/adapter/observability/otelcore"
	"github.com/syntaxfa/quick-connect/adapter/observability/traceotela"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/pkg/grpcserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

type Observability struct {
	Core  otelcore.Config   `koanf:"core"`
	Trace traceotela.Config `koanf:"trace"`
}

type Config struct {
	GRPCServer    grpcserver.Config `koanf:"grpc_server"`
	Logger        logger.Config     `koanf:"logger"`
	Postgres      postgres.Config   `koanf:"postgres"`
	Observability Observability     `koanf:"observability"`
}
