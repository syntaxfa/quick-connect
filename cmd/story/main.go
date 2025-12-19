package main

import (
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/app/storyapp"
	"github.com/syntaxfa/quick-connect/cmd/story/command"
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
	var cfg storyapp.Config

	workingDir, gErr := os.Getwd()
	if gErr != nil {
		panic(gErr)
	}

	options := config.Option{
		Prefix:       "STORY_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "story", "config.yml"),
		CallBackEnv:  nil,
	}
	config.Load(options, &cfg, nil)

	log := logger.New(cfg.Logger, nil, "story")

	root := &cobra.Command{
		Use:     "story",
		Short:   "story application",
		Version: "1.x.x",
	}

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	root.AddCommand(
		command.Server{}.Command(cfg, log, trap),
		command.Migrate{}.Command(cfg.Postgres, log),
	)

	if exErr := root.Execute(); exErr != nil {
		log.Error("failed to execute root command", slog.String("error", exErr.Error()))

		os.Exit(1)
	}
}
