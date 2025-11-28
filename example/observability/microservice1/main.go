package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"

	"github.com/syntaxfa/quick-connect/adapter/observability/metricotela"
	"github.com/syntaxfa/quick-connect/adapter/observability/otelcore"
	"github.com/syntaxfa/quick-connect/adapter/observability/traceotela"
	"github.com/syntaxfa/quick-connect/config"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice1"
	"github.com/syntaxfa/quick-connect/pkg/logger"
	"go.opentelemetry.io/otel/metric"
)

func getMemoryUsage() uint64 {
	var memStats runtime.MemStats

	runtime.ReadMemStats(&memStats)

	currentMemoryUsage := memStats.HeapAlloc

	return currentMemoryUsage
}

// main
//
//	@schemes					http https
//	@securityDefinitions.apiKey	JWT
//	@in							header
//	@name						Authorization
//	@description				JWT security accessToken. Please add it in the format "Bearer {AccessToken}" to authorize your requests.
func main() {
	var cfg microservice1.Config

	workingDir, gErr := os.Getwd()
	if gErr != nil {
		panic(gErr)
	}

	options := config.Option{
		Prefix:       "MICROSERVICE1_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "example", "observability", "microservice1", "config.yml"),
		CallBackEnv:  nil,
	}
	config.Load(options, &cfg, nil)

	ctx, cancelFunc := context.WithCancel(context.Background())

	resource, sErr := otelcore.NewResource(ctx, cfg.Observability.Core)
	if sErr != nil {
		panic(sErr)
	}

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	log := logger.New(cfg.Logger, nil, "microservice1")

	if mErr := metricotela.InitMetric(ctx, cfg.Observability.Metric, resource, log); mErr != nil {
		log.Error(mErr.Error())
	}

	metricotela.SetMeter(cfg.Observability.Core.ServiceName)

	_, err := metricotela.Meter().Int64ObservableGauge(
		"system.memory.heap",
		metric.WithDescription("Memory usage of the allocated heap objects khekhe"),
		metric.WithUnit("By"),
		metric.WithInt64Callback(func(_ context.Context, o metric.Int64Observer) error {
			memoryUsage := getMemoryUsage()
			o.Observe(int64(memoryUsage))

			return nil
		}),
	)
	if err != nil {
		panic(err)
	}

	traceCtx, tErr := traceotela.InitTracer(ctx, cfg.Observability.Trace, resource)
	if tErr != nil {
		panic(tErr)
	}

	traceotela.SetTracer(cfg.Observability.Core.ServiceName)

	app := microservice1.Setup(cfg, log, trap)

	app.Start()

	<-trap

	if tErr := traceCtx(ctx); tErr != nil {
		panic(tErr)
	}

	defer cancelFunc()

	fmt.Println("stopped")
}
