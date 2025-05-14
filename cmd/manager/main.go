package main

import (
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/app/managerapp"
	_ "github.com/syntaxfa/quick-connect/app/managerapp/docs"
	"github.com/syntaxfa/quick-connect/cmd/manager/command"
	"github.com/syntaxfa/quick-connect/config"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

// main
//
//	@schemes					http https
//	@securityDefinitions.apiKey	JWT
//	@in							header
//	@name						Authorization
//	@description				JWT security accessToken. Please add it in the format "Bearer {AccessToken}" to authorize your requests.
func main() {
	var cfg managerapp.Config

	workingDir, gErr := os.Getwd()
	if gErr != nil {
		panic(gErr)
	}

	options := config.Option{
		Prefix:       "MANAGER_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "manager", "config.yml"),
		CallBackEnv:  nil,
	}
	config.Load(options, &cfg, nil)

	log := logger.New(cfg.Logger, nil, true, "manager")

	root := &cobra.Command{
		Use:   "manager",
		Short: "manager application",
		Long: "The application manager is used for managing Quick-Connect, and it also includes a back-office" +
			" that controls the access management of administrators.",
		Version: "1.0.0",
	}

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	root.AddCommand(
		command.Server{}.Command(cfg, log, trap),
		command.Migrate{}.Command(cfg.Postgres, log),
	)

	if eErr := root.Execute(); eErr != nil {
		log.Error("failed to execute root command", slog.String("error", eErr.Error()))

		os.Exit(1)
	}
}
