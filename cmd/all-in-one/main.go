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
	chatCommand "github.com/syntaxfa/quick-connect/cmd/chat/command"
	managerCommand "github.com/syntaxfa/quick-connect/cmd/manager/command"
	notificationCommand "github.com/syntaxfa/quick-connect/cmd/notification/command"
	"github.com/syntaxfa/quick-connect/config"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

func main() {
	serviceCfg := setServiceConfigs()

	// Loggers.
	serviceLog := setLoggers(serviceCfg)

	// commands.
	managerRoot := &cobra.Command{
		Use:   "manager",
		Short: "manager service commands",
	}
	managerRoot.AddCommand(managerCommand.Migrate{}.Command(serviceCfg.ManagerCfg.Postgres, serviceLog.ManagerLog))
	managerRoot.AddCommand(managerCommand.CreateUser{}.Command(serviceCfg.ManagerCfg, serviceLog.ManagerLog, nil))

	chatRoot := &cobra.Command{
		Use:   "chat",
		Short: "chat service commands",
	}
	chatRoot.AddCommand(chatCommand.Migrate{}.Command(serviceCfg.ChatCfg.Postgres, serviceLog.ChatLog))

	notificationRoot := &cobra.Command{
		Use:   "notification",
		Short: "notification service commands",
	}
	notificationRoot.AddCommand(notificationCommand.Migrate{}.Command(serviceCfg.NotificationCfg.Postgres, serviceLog.NotificationLog))

	root := &cobra.Command{
		Use:     "quick-connect",
		Short:   "quick connect all in one",
		Version: "1.0.0",
	}

	trap := make(chan os.Signal, 1)
	signal.Notify(trap, syscall.SIGINT, syscall.SIGTERM)

	root.AddCommand(
		managerRoot,
		chatRoot,
		notificationRoot,
		command.Server{}.Command(serviceCfg, serviceLog, trap),
	)

	if exErr := root.Execute(); exErr != nil {
		panic(fmt.Sprintf("failed to execute root command, error: %s", exErr.Error()))
	}
}

func setServiceConfigs() command.ServiceConfig {
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

	return command.ServiceConfig{
		ManagerCfg:      managerCfg,
		ChatCfg:         chatCfg,
		NotificationCfg: notificationCfg,
		AdminCfg:        adminCfg,
	}
}

func setLoggers(serviceCfg command.ServiceConfig) command.Logger {
	return command.Logger{
		ManagerLog:      logger.New(serviceCfg.ManagerCfg.Logger, nil, "admin"),
		ChatLog:         logger.New(serviceCfg.ChatCfg.Logger, nil, "chat"),
		NotificationLog: logger.New(serviceCfg.NotificationCfg.Logger, nil, "notification"),
		AdminLog:        logger.New(serviceCfg.AdminCfg.Logger, nil, "admin"),
	}
}
