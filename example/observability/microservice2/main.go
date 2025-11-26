package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/syntaxfa/quick-connect/adapter/observability/metricotela"
	"github.com/syntaxfa/quick-connect/adapter/observability/otelcore"
	"github.com/syntaxfa/quick-connect/adapter/observability/traceotela"
	"github.com/syntaxfa/quick-connect/config"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice2"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

func main() {
	var cfg microservice2.Config

	workingDir, gErr := os.Getwd()
	if gErr != nil {
		panic(gErr)
	}

	options := config.Option{
		Prefix:       "MICROSERVICE1_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "example", "observability", "microservice2", "config.yml"),
		CallBackEnv:  nil,
	}
	config.Load(options, &cfg, nil)

	ctx, cancelFunc := context.WithCancel(context.Background())

	resource, sErr := otelcore.NewResource(ctx, cfg.Observability.Core)
	if sErr != nil {
		panic(sErr)
	}

	log := logger.New(cfg.Logger, nil, "microservice2")

	traceCtx, tErr := traceotela.InitTracer(ctx, cfg.Observability.Trace, resource)
	if tErr != nil {
		panic(tErr)
	}

	traceotela.SetTracer(cfg.Observability.Core.ServiceName)

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	if mErr := metricotela.InitMetric(ctx, cfg.Observability.Metric, resource, log); mErr != nil {
		log.Error(mErr.Error())
	}

	metricotela.SetMeter(cfg.Observability.Core.ServiceName)

	app := microservice2.Setup(cfg, log, trap)

	app.Start()

	<-trap

	if tErr := traceCtx(ctx); tErr != nil {
		panic(tErr)
	}

	cancelFunc()

	fmt.Println("stopped")
}
