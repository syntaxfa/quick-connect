package cachemanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"
)

// ErrKeyNotFound is a custom error returned when a key is not found in the cache.
var ErrKeyNotFound = errors.New("key not found in cache")

// CacheClient is a generic interface for basic cache operations.
// This interface is completely abstracted from implementation details (Redis, memcached ...).
type CacheClient interface {
	// Set stores a value associated with a given key and an expiration duration.
	// The 'value' here is expected to be a byte slice, allowing the CacheManager
	// to handle marshalling of various data types into a storable format.
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error

	// Get retrieves the value associated with a key.
	// It returns the value as a byte slice, If the key is not found, it returns ErrKeyNotFound.
	Get(ctx context.Context, key string) ([]byte, error)

	MGet(ctx context.Context, keys ...string) ([]interface{}, error)

	// Delete removes one or more keys from the cache.
	Delete(ctx context.Context, keys ...string) error

	// GetTTL retrieves the remaining time-to-live (TTL) for a key.
	// If the key is not found, it returns ErrKeyNotFound.
	// If the key exists nut has no TTL (it's persistent), it returns 0.
	GetTTL(ctx context.Context, key string) (time.Duration, error)

	// Incr atomically increments the integer value of a key by one.
	Incr(ctx context.Context, key string) (int64, error)

	// Decr atomically decrements the integer value of a key by one.
	Decr(ctx context.Context, key string) (int64, error)

	// Expire sets a timeout on a key.
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

// CacheManager creates a new instance of CacheManager.
type CacheManager struct {
	client CacheClient
	logger *slog.Logger
}

// New creates a new instance of CacheManager.
// It accepts an implementation of CacheClient interface.
func New(client CacheClient, logger *slog.Logger) *CacheManager {
	return &CacheManager{
		client: client,
		logger: logger,
	}
}

// Set marshals the provided value to JSON and then stores it in the cache.
// It handles the conversion of any Go type to a byte slice for storage.
func (c *CacheManager) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	marshalValue, jErr := json.Marshal(value)
	if jErr != nil {
		return fmt.Errorf("cachemanager: failed to marshal value for key '%s': '%s'", key, jErr.Error())
	}

	if sErr := c.client.Set(ctx, key, marshalValue, expiration); sErr != nil {
		return fmt.Errorf("cachemanager: failed to set key '%s' in cache: %s", key, sErr.Error())
	}

	return nil
}

// Get retrieves a value from the cache, unmarshal it from JSON into the 'dest' pointers.
// If the key is not found, it returns ErrKeyFound.
func (c *CacheManager) Get(ctx context.Context, key string, dest any) error {
	data, gErr := c.client.Get(ctx, key)
	if gErr != nil {
		if errors.Is(gErr, ErrKeyNotFound) {
			return ErrKeyNotFound
		}

		return fmt.Errorf("cachemanager: failed to get key '%s' from cache: %s", key, gErr.Error())
	}

	// Unmarshal the JSON bytes into the destination variables.
	if uErr := json.Unmarshal(data, dest); uErr != nil {
		return fmt.Errorf("cachemanager: failed to unmarshal value for key '%s': %s", key, uErr.Error())
	}

	return nil
}

// MGet retrieves multiple values from the cache for the given keys.
// It populates the `destMap` for keys are found(cache hint).
// The keys of destMap must be the cache keys, and values must be pointers
// to the destination structs
// It returns a slice of keys that were not found in the cache(cache misses),
// which can then be fetched from the primary data source.
func (c *CacheManager) MGet(ctx context.Context, destMap map[string]any, keys ...string) (missedKeys []string, err error) {
	if len(keys) == 0 {
		return nil, nil
	}

	results, gErr := c.client.MGet(ctx, keys...)
	if gErr != nil {
		return nil, fmt.Errorf("cachemanager: failed to mget keys from cache: %s", gErr.Error())
	}

	if len(results) != len(keys) {
		return nil, errors.New("mismatched number of keys and results")
	}

	for i, result := range results {
		key := keys[i]

		if result == nil {
			missedKeys = append(missedKeys, key)

			continue
		}

		dest, ok := destMap[key]
		if !ok {
			// the caller must provider a destination for each key.
			return nil, fmt.Errorf("destination for key `%s` not provided in destMap", key)
		}

		var data []byte
		switch v := result.(type) {
		case []byte:
			data = v
		case string:
			data = []byte(v)
		default:
			return nil, fmt.Errorf("unsupported cache value type for key `%s`: %T", key, v)
		}

		// Unmarshal the JSON bytes into the destination pointer.
		if uErr := json.Unmarshal(data, dest); uErr != nil {
			missedKeys = append(missedKeys, key)
			c.logger.WarnContext(ctx, fmt.Sprintf("failed to unmarshal cache for key `%s`, treating as miss", key),
				slog.String("error", uErr.Error()))

			continue
		}
	}

	return missedKeys, nil
}

// Delete removes one or more keys from cache.
func (c *CacheManager) Delete(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	err := c.client.Delete(ctx, keys...)
	if err != nil {
		return fmt.Errorf("cachemanager: failed to delete keys '%v' from cache: %s", keys, err.Error())
	}

	return nil
}

// GetTTL retrieves the remaining time-to-live for a key.
// It returns 0 if the exists but has no expirations (persistent).
func (c *CacheManager) GetTTL(ctx context.Context, key string) (time.Duration, error) {
	ttl, err := c.client.GetTTL(ctx, key)
	if err != nil {
		if errors.Is(err, ErrKeyNotFound) {
			return 0, ErrKeyNotFound
		}

		return 0, fmt.Errorf("cachemanager: failed to get TTL for key '%s':%s", key, err.Error())
	}

	return ttl, nil
}

// Incr atomically increments the integer value of a key by one.
// This is a pass-through method and does not involve JSON marshalling.
func (c *CacheManager) Incr(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Incr(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("cachemanager: failed to incr key '%s': %s", key, err.Error())
	}

	return val, nil
}

// Decr atomically decrements the integer value of a key by one.
// This is a pass-through method and does not involve JSON marshalling.
func (c *CacheManager) Decr(ctx context.Context, key string) (int64, error) {
	val, err := c.client.Decr(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("cachemanager: failed to decr key '%s': %s", key, err.Error())
	}

	return val, nil
}

// Expire sets a timeout on a key.
// This is a pass-through method.
func (c *CacheManager) Expire(ctx context.Context, key string, expiration time.Duration) error {
	err := c.client.Expire(ctx, key, expiration)
	if err != nil {
		return fmt.Errorf("cachemanager: failed to expire key '%s': %s", key, err.Error())
	}

	return nil
}
