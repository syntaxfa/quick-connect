package main

import (
	"context"
	"time"

	"github.com/syntaxfa/quick-connect/pkg/logger"
	"github.com/syntaxfa/quick-connect/pkg/websocket"
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

	server := websocket.New(websocket.Config{
		Host: "localhost",
		Port: "5000",
	}, logger.L(), echoHandler)

	go server.Start()

	time.Sleep(time.Second * 5)

	ctx, cancelFunc := context.WithTimeout(context.Background(), time.Second*2)
	defer cancelFunc()

	server.Stop(ctx)
}
