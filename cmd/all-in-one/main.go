package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/spf13/cobra"
	"github.com/syntaxfa/quick-connect/app/adminapp"
	"github.com/syntaxfa/quick-connect/app/chatapp"
	"github.com/syntaxfa/quick-connect/app/managerapp"
	"github.com/syntaxfa/quick-connect/app/notificationapp"
	"github.com/syntaxfa/quick-connect/cmd/all-in-one/command"
	"github.com/syntaxfa/quick-connect/config"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

func main() {
	var managerCfg managerapp.Config
	var chatCfg chatapp.Config
	var notificationCfg notificationapp.Config
	var adminCfg adminapp.Config

	// Configs.
	workingDir, gErr := os.Getwd()
	if gErr != nil {
		panic(gErr)
	}

	managerOptions := config.Option{
		Prefix:       "MANAGER_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "manager", "config.yml"),
		CallBackEnv:  nil,
	}
	config.Load(managerOptions, &managerCfg, nil)

	chatOption := config.Option{
		Prefix:       "CHAT_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "chat", "config.yml"),
		CallBackEnv:  nil,
	}
	config.Load(chatOption, &chatCfg, nil)

	notificationOptions := config.Option{
		Prefix:       "NOTIFICATION_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "notification", "config.yml"),
		CallBackEnv:  nil,
	}
	config.Load(notificationOptions, &notificationCfg, nil)

	adminOptions := config.Option{
		Prefix:       "ADMIN_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: filepath.Join(workingDir, "deploy", "admin", "config.yml"),
		CallBackEnv:  nil,
	}
	config.Load(adminOptions, &adminCfg, nil)

	cfg := command.Config{
		ManagerCfg:      managerCfg,
		ChatCfg:         chatCfg,
		NotificationCfg: notificationCfg,
		AdminCfg:        adminCfg,
	}

	// Loggers.
	managerLog := logger.New(managerCfg.Logger, nil, "admin")
	chatLog := logger.New(chatCfg.Logger, nil, "chat")
	notificationLog := logger.New(notificationCfg.Logger, nil, "notification")
	adminLog := logger.New(adminCfg.Logger, nil, "admin")

	log := command.Logger{
		ManagerLog:      managerLog,
		ChatLog:         chatLog,
		NotificationLog: notificationLog,
		AdminLog:        adminLog,
	}

	root := &cobra.Command{
		Use:     "quick-connect",
		Short:   "quick connect all in one",
		Version: "1.0.0",
	}

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	root.AddCommand(
		command.Server{}.Command(cfg, log, trap),
	)

	if exErr := root.Execute(); exErr != nil {
		panic(fmt.Sprintf("failed to execute root command, error: %s", exErr.Error()))
	}
}
