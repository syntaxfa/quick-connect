package main

import (
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/app/notificationapp"
	"github.com/syntaxfa/quick-connect/cmd/notification/command"
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
	var cfg notificationapp.Config

	workingDir, gErr := os.Getwd()
	if gErr != nil {
		panic(gErr)
	}

	options := config.Option{
		Prefix:       "NOTIFICATION_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "notification", "config.yml"),
		CallBackEnv:  nil,
	}
	config.Load(options, &cfg, nil)

	log := logger.New(cfg.Logger, nil, true, "notification")

	root := &cobra.Command{
		Use:   "notification",
		Short: "notification application",
		Long: "The application notification is used for managing notification, and it also include a UI interface" +
			"for monitoring notifications.",
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
