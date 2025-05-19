package main

import (
	"context"
	"fmt"
	"github.com/syntaxfa/quick-connect/adapter/observability/otelcore"
	"github.com/syntaxfa/quick-connect/adapter/observability/traceotela"
	"github.com/syntaxfa/quick-connect/config"
	"github.com/syntaxfa/quick-connect/example/observability/internal/microservice1"
	"github.com/syntaxfa/quick-connect/pkg/logger"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
)

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

	traceCtx, tErr := traceotela.InitTracer(ctx, cfg.Observability.Trace, resource)
	if tErr != nil {
		panic(tErr)
	}

	traceotela.SetTracer(cfg.Observability.Core.ServiceName)

	log := logger.New(cfg.Logger, nil, false, "microservice1")

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	app := microservice1.Setup(cfg, log, trap)

	app.Start()

	<-trap

	if tErr := traceCtx(ctx); tErr != nil {
		panic(tErr)
	}

	cancelFunc()

	fmt.Println("stopped")
}
