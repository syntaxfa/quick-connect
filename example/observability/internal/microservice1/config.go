package microservice1

import (
	"github.com/syntaxfa/quick-connect/adapter/observability/metricotela"
	"github.com/syntaxfa/quick-connect/adapter/observability/otelcore"
	"github.com/syntaxfa/quick-connect/adapter/observability/traceotela"
	"github.com/syntaxfa/quick-connect/adapter/postgres"
	"github.com/syntaxfa/quick-connect/pkg/grpcclient"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
	"time"
)

type Observability struct {
	Core   otelcore.Config    `koanf:"core"`
	Trace  traceotela.Config  `koanf:"trace"`
	Metric metricotela.Config `koanf:"metric"`
}

type Config struct {
	HTTPServer      httpserver.Config `koanf:"http_server"`
	ShutdownTimeout time.Duration     `koanf:"shutdown_timeout"`
	Logger          logger.Config     `koanf:"logger"`
	Observability   Observability     `koanf:"observability"`
	Postgres        postgres.Config   `koanf:"postgres"`
	GRPCClient      grpcclient.Config `koanf:"grpc_client"`
}
