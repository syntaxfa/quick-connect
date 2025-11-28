package main

import (
	"context"
	"time"

	"github.com/syntaxfa/quick-connect/pkg/httpserver"
	"github.com/syntaxfa/quick-connect/pkg/logger"
)

func main() {
	log := logger.New(logger.Config{
		FilePath:         "logs.json",
		UseLocalTime:     false,
		FileMaxSizeInMB:  1,
		FileMaxAgeInDays: 10,
		MaxBackup:        0,
		Compress:         false,
	}, nil, "example")

	server := httpserver.New(httpserver.Config{
		Port: 5050,
		Cors: httpserver.Cors{AllowOrigins: []string{"localhost"}},
	}, log)

	go server.Start()

	time.Sleep(time.Second * 5)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()

	server.Stop(ctx)
}
