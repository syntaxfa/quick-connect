package main

import (
	"context"
	"errors"
	"fmt"
	"github.com/syntaxfa/quick-connect/adapter/redis"
	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
	"log/slog"
	"time"
)

type Value struct {
	UserID string `json:"external_user_id"`
}

func main() {
	redisCfg := redis.Config{
		Host:     "localhost",
		Port:     12439,
		Password: "Z9265UQfrFiCYWM",
		DB:       0,
	}

	redisAdapter := redis.New(redisCfg, slog.Default())

	cache := cachemanager.New(redisAdapter)

	key := "notification:users:1"
	ctx := context.Background()
	value := Value{UserID: "user_id"}

	if err := cache.Set(ctx, key, &value, 0); err != nil {
		slog.Default().Error(err.Error())
	}

	var cacheValue Value
	if err := cache.Get(ctx, key, &cacheValue); err != nil {
		if errors.Is(err, cachemanager.ErrKeyNotFound) {
			fmt.Printf("key: %s not found\n", "sss")
		} else {
			slog.Default().Error(err.Error())
		}
	} else {
		fmt.Printf("%+v\n", cacheValue)
	}

	fmt.Println("check key ttl")

	ttl, err := cache.GetTTL(ctx, key)
	if err != nil {
		slog.Default().Error(err.Error())
	}

	fmt.Println(ttl)

	if err := cache.Delete(ctx, key); err != nil {
		slog.Default().Error(err.Error())
	}

	ttl, err = cache.GetTTL(ctx, key)
	if err != nil {
		slog.Default().Error(err.Error())
	}

	time.Sleep(time.Second * 5)
}
