package microservice1

import (
	"github.com/syntaxfa/quick-connect/adapter/observability/otelcore"
	"github.com/syntaxfa/quick-connect/adapter/observability/traceotela"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
	"time"
)

type Observability struct {
	Core  otelcore.Config   `koanf:"core"`
	Trace traceotela.Config `koanf:"trace"`
}

type Config struct {
	HTTPServer      httpserver.Config `koanf:"http_server"`
	ShutdownTimeout time.Duration     `koanf:"shutdown_timeout"`
	Logger          logger.Config     `koanf:"logger"`
	Observability   Observability     `koanf:"observability"`
}
