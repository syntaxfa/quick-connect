package main

import "github.com/syntaxfa/quick-connect/pkg/logger"

func main() {
	log := logger.New(logger.Config{
		FilePath:         "/logs/example/logs.json",
		UseLocalTime:     false,
		FileMaxSizeInMB:  10,
		FileMaxAgeInDays: 30,
		MaxBackup:        0,
		Compress:         false,
	}, nil, true, "example")

	log.Info("Hello world!")
}
