package redis

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/syntaxfa/quick-connect/pkg/cachemanager"
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

// New creates a new instance of Redis Adapter.
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

// Client get redis client.
func (a *Adapter) Client() *redis.Client {
	return a.client
}

// Close redis connections.
func (a *Adapter) Close() error {
	return a.client.Close()
}

// Set implement the Set method of the cachemanager.CacheClient interface for Redis.
func (a *Adapter) Set(ctx context.Context, key string, value []byte, expiration time.Duration) error {
	return a.client.Set(ctx, key, value, expiration).Err()
}

// Get implement the Get method of the cachemanager.CacheClient interface for Redis.
// It translates redis.Nil error into cachemanager.ErrKeyNotFound.
func (a *Adapter) Get(ctx context.Context, key string) ([]byte, error) {
	data, gErr := a.client.Get(ctx, key).Bytes()
	if errors.Is(gErr, redis.Nil) {
		return nil, cachemanager.ErrKeyNotFound
	}

	return data, gErr
}

// MGet implement the MGet method of the cachemanager.CacheClient interfcae for redis.
func (a *Adapter) MGet(ctx context.Context, keys ...string) ([]interface{}, error) {
	return a.client.MGet(ctx, keys...).Result()
}

// Delete implement the Delete method of the cachemanager.CacheClient interface for Redis.
func (a *Adapter) Delete(ctx context.Context, keys ...string) error {
	return a.client.Del(ctx, keys...).Err()
}

// GetTTL implement the GetTTL method of the cachemanager.CacheClient interface for Redis.
// It translates Redis's TTL return values (-1 for persistent, -2 for not found)
// into the generic time.Duration and cachemanager.ErrKeyNotFound.
func (a *Adapter) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := a.client.TTL(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	if ttl == -2*time.Nanosecond { // Redis returns -2 for non-existent key
		return 0, cachemanager.ErrKeyNotFound
	}

	if ttl == -1*time.Nanosecond { // Redis returns -1 for a key with no expirations (persistent)
		return 0, nil // cachemanager.CacheClient defines 0 duration for persistent keys
	}

	return ttl, nil
}
