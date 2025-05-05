package main

import (
	"log/slog"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/app/chatapp"
	"github.com/syntaxfa/quick-connect/cmd/chat/command"
	"github.com/syntaxfa/quick-connect/config"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

//	@schemes					http https
//	@securityDefinitions.apiKey	JWT
//	@in							header
//	@name						Authorization
//	@description				JWT security accessToken. Please add it in the format "Bearer {AccessToken}" to authorize your requests.
func main() {
	var cfg chatapp.Config

	workingDir, gErr := os.Getwd()
	if gErr != nil {
		panic(gErr)
	}

	options := config.Option{
		Prefix:       "CHAT_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "chat", "development", "config.yml"),
		CallBackEnv:  nil,
	}
	config.Load(options, &cfg, nil)

	log := logger.New(cfg.Logger, nil, true, "chat")

	root := &cobra.Command{
		Use:     "chat",
		Short:   "chat application",
		Version: "1.0.0",
	}

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	root.AddCommand(
		command.Server{}.Command(cfg, log, trap),
	)

	if eErr := root.Execute(); eErr != nil {
		logger.L().Error("failed to execute root command", slog.String("error", eErr.Error()))

		os.Exit(1)
	}
}
