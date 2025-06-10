package redis

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Host     string `koanf:"host"`
	Port     int    `koanf:"port"`
	Password string `koanf:"password"`
	DB       int    `koanf:"db"`
}

type Adapter struct {
	client *redis.Client
}

func New(cfg Config, logger *slog.Logger) *Adapter {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		logger.Error("redis connection failed", slog.String("error", err.Error()))

		panic("redis connection error, check logs")
	}

	return &Adapter{client: rdb}
}

func (a *Adapter) Client() *redis.Client {
	return a.client
}

func (a *Adapter) Close() error {
	return a.client.Close()
}
