package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/oklog/ulid/v2"
	"github.com/syntaxfa/quick-connect/adapter/pubsub/redispubsub"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	"github.com/syntaxfa/quick-connect/pkg/pubsub"
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

	redisPubSub := redispubsub.New(re)

	svc := NewNotificationService(redisPubSub)

	go svc.publish()

	svc.Receive()
}

type NotificationService struct {
	pubsub pubsub.PubSub
}

func NewNotificationService(pubSub pubsub.PubSub) NotificationService {
	return NotificationService{
		pubsub: pubSub,
	}
}

func (n NotificationService) publish() {
	ctx := context.Background()
	channel := "notification"

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

		if pErr := n.pubsub.Publish(ctx, channel, jsonData); pErr != nil {
			panic(pErr)
		}

		time.Sleep(time.Millisecond * 500)
	}
}

func (n NotificationService) Receive() {
	ctx := context.Background()
	channels := []string{"notification"}

	pubSub := n.pubsub.Subscribe(ctx, channels...)

	for {
		message, rErr := pubSub.ReceiveMessage(ctx)
		if rErr != nil {
			panic(rErr)
		}

		var notification NotificationMessage
		uErr := json.Unmarshal(message, &notification)
		if uErr != nil {
			panic(uErr)
		}

		fmt.Printf("%+v\n", notification)
	}
}
