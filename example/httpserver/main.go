package main

import (
	"context"
	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
	"time"
)

func main() {
	logger.SetDefault(logger.Config{
		FilePath:         "logs.json",
		UseLocalTime:     false,
		FileMaxSizeInMB:  1,
		FileMaxAgeInDays: 10,
		MaxBackup:        0,
		Compress:         false,
	}, nil, true, "example")

	server := httpserver.New(httpserver.Config{
		Port: 5050,
		Cors: httpserver.Cors{AllowOrigins: []string{"localhost"}},
	}, logger.L())

	go server.Start()

	time.Sleep(time.Second * 5)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()

	server.Stop(ctx)
}
