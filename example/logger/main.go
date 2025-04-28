package main

import "github.com/syntaxfa/quick-connect/pkg/logger"

func main() {
	logger.SetDefault(logger.Config{
		FilePath:         "/logs/example/logs.json",
		UseLocalTime:     false,
		FileMaxSizeInMB:  10,
		FileMaxAgeInDays: 30,
		MaxBackup:        0,
		Compress:         false,
	}, nil, true)

	logger.L().Info("Hello world!")
}
