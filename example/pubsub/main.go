package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	"github.com/syntaxfa/quick-connect/types"
)

type NotificationMessage struct {
	ID        types.ID `json:"id"`
	Title     string   `json:"title"`
	Body      string   `json:"body"`
	Timestamp int64    `json:"timestamp"`
}

func main() {
	cfg := redis.Config{
		Host:     "localhost",
		Port:     12434,
		Password: "Z9265UQfrFiCYWMMJF4uvTEmJA7rEauJ",
		DB:       0,
	}

	logger := slog.Default()

	re := redis.New(cfg, logger)
	defer func() {
		if cErr := re.Close(); cErr != nil {
			panic(cErr)
		}
	}()

	fmt.Println("redis connection established")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	channel := "notification"

	go publisher(ctx, channel, re)

	subscriber(ctx, channel, re)
}

func publisher(ctx context.Context, channel string, re *redis.Adapter) {
	for i := 0; ; i++ {
		message := NotificationMessage{
			ID:        types.ID(ulid.Make().String()),
			Title:     fmt.Sprintf("title %d", i),
			Body:      fmt.Sprintf("body %d", i),
			Timestamp: time.Now().Unix(),
		}
		jsonData, mErr := json.Marshal(message)
		if mErr != nil {
			panic(mErr)
		}

		re.Client().Publish(ctx, channel, jsonData)

		time.Sleep(time.Millisecond * 500)
	}
}

func subscriber(ctx context.Context, channel string, re *redis.Adapter) {
	pubSub := re.Client().Subscribe(ctx, channel)

	for {
		message, rErr := pubSub.ReceiveMessage(ctx)
		if rErr != nil {
			panic(rErr)
		}

		var notification NotificationMessage
		uErr := json.Unmarshal([]byte(message.Payload), &notification)
		if uErr != nil {
			panic(uErr)
		}

		fmt.Printf("%+v\n", notification)
	}
}
