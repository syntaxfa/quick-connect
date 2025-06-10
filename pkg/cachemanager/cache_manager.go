package cachemanager

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

	// Delete removes one or more keys from the cache.
	Delete(ctx context.Context, keys ...string) error

	// GetTTL retrieves the remaining time-to-live (TTL) for a key.
	// If the key is not found, it returns ErrKeyNotFound.
	// If the key exists nut has no TTL (it's persistent), it returns 0.
	GetTTL(ctx context.Context, key string) (time.Duration, error)
}

// CacheManager creates a new instance of CacheManager.
type CacheManager struct {
	client CacheClient
}

// New creates a new instance of CacheManager.
// It accepts an implementation of CacheClient interface.
func New(client CacheClient) *CacheManager {
	return &CacheManager{
		client: client,
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
