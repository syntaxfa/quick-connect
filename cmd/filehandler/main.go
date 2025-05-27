package main

import (
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/app/filehandlerapp"
	"github.com/syntaxfa/quick-connect/cmd/filehandler/command"
	"github.com/syntaxfa/quick-connect/config"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

func main() {
	var cfg filehandlerapp.Config

	workingDir, gErr := os.Getwd()
	if gErr != nil {
		panic(gErr)
	}

	options := config.Option{
		Prefix:       "FILE_HANDLER_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "filehandler", "config.yml"),
		CallBackEnv:  nil,
	}
	config.Load(options, &cfg, nil)

	log := logger.New(cfg.Logger, nil, true, "filehandler")

	root := &cobra.Command{
		Use:     "filehandler",
		Short:   "filehandler application",
		Long:    "The application is used for managing Quick-Connect files, and controls the client access for files",
		Version: "1.0.0",
	}

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	root.AddCommand(
		command.Migrate{}.Command(cfg.Postgres, log),
	)

	if eErr := root.Execute(); eErr != nil {
		log.Error("failed to execute root command", slog.String("error", eErr.Error()))

		os.Exit(1)
	}
}
