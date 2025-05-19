package main

import (
	"context"
	"fmt"
	"github.com/syntaxfa/quick-connect/adapter/observability/otelcore"
	"github.com/syntaxfa/quick-connect/adapter/observability/traceotela"
	"github.com/syntaxfa/quick-connect/config"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice2"
	"github.com/syntaxfa/quick-connect/pkg/logger"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
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

	traceCtx, tErr := traceotela.InitTracer(ctx, cfg.Observability.Trace, resource)
	if tErr != nil {
		panic(tErr)
	}

	traceotela.SetTracer(cfg.Observability.Core.ServiceName)

	log := logger.New(cfg.Logger, nil, true, "microservice2")

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	app := microservice2.Setup(cfg, log, trap)

	app.Start()

	<-trap

	if tErr := traceCtx(ctx); tErr != nil {
		panic(tErr)
	}

	cancelFunc()

	fmt.Println("stopped")
}
