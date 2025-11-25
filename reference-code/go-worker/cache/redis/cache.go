// Package redis provides Redis implementation of cache service.
//
// Purpose:
// - Production-ready caching with Redis
// - Supports TTL, pattern-based deletion, distributed locking
// - Thread-safe and connection-pooled
package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/orchard9/go-core-http-toolkit/cache"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
	"github.com/orchard9/peach/apps/email-worker/internal/ports"
)

// CacheService implements ports.CacheService using Redis
//
// Thread-Safety: Safe for concurrent use (Redis client is thread-safe)
//
// Features:
// - Automatic TTL management
// - Pattern-based key deletion
// - Distributed locking with SetNX
// - Connection pooling via go-core-http-toolkit
type CacheService struct {
	client *cache.RedisClient
}

// NewCacheService creates a new Redis cache service
//
// Parameters:
// - client: Redis client from go-core-http-toolkit/cache
//
// Returns configured cache service ready for use
func NewCacheService(client *cache.RedisClient) ports.CacheService {
	return &CacheService{
		client: client,
	}
}

// Get retrieves a value from cache
//
// Returns:
// - domain.ErrNotFound if key doesn't exist
// - domain.ErrInternal for Redis errors
func (c *CacheService) Get(ctx context.Context, key string) (string, error) {
	var value string
	found, err := c.client.Get(ctx, key, &value)
	if err != nil {
		return "", fmt.Errorf("%w: failed to get cache key %s: %v", domain.ErrInternal, key, err)
	}
	if !found {
		return "", fmt.Errorf("%w: cache key %s not found", domain.ErrNotFound, key)
	}
	return value, nil
}

// Set stores a value in cache with expiration
//
// Parameters:
// - key: Cache key
// - value: Value to store
// - ttl: Time-to-live (0 = no expiration)
//
// Returns:
// - domain.ErrInternal for Redis errors
func (c *CacheService) Set(ctx context.Context, key string, value string, ttl time.Duration) error {
	err := c.client.Set(ctx, key, value, ttl)
	if err != nil {
		return fmt.Errorf("%w: failed to set cache key %s: %v", domain.ErrInternal, key, err)
	}
	return nil
}

// Delete removes a value from cache
//
// Idempotent: No error if key doesn't exist
//
// Returns:
// - domain.ErrInternal for Redis errors
func (c *CacheService) Delete(ctx context.Context, key string) error {
	if err := c.client.Delete(ctx, key); err != nil {
		return fmt.Errorf("%w: failed to delete cache key %s: %v", domain.ErrInternal, key, err)
	}
	return nil
}

// DeletePattern removes all keys matching a pattern
//
// Pattern Syntax (Redis SCAN):
// - "*" matches any characters
// - "?" matches single character
// - Example: "users:*" matches "users:123", "users:456"
//
// Implementation:
// - Uses SCAN for cursor-based iteration (safe for large keyspaces)
// - Avoids blocking Redis with KEYS command
//
// Returns:
// - domain.ErrInternal for Redis errors
func (c *CacheService) DeletePattern(ctx context.Context, pattern string) error {
	// Get underlying Redis client for SCAN operation
	redisClient := c.client.GetClient()

	// Use SCAN to iterate through matching keys
	var cursor uint64
	var keys []string

	for {
		var scanKeys []string
		var err error

		scanKeys, cursor, err = redisClient.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return fmt.Errorf("%w: failed to scan cache pattern %s: %v", domain.ErrInternal, pattern, err)
		}

		keys = append(keys, scanKeys...)

		// Break when cursor returns to 0
		if cursor == 0 {
			break
		}
	}

	// Delete all matching keys
	if len(keys) > 0 {
		if err := c.client.Delete(ctx, keys...); err != nil {
			return fmt.Errorf("%w: failed to delete cache keys for pattern %s: %v", domain.ErrInternal, pattern, err)
		}
	}

	return nil
}

// Exists checks if a key exists in cache
//
// Returns:
// - true if key exists
// - false if key doesn't exist or expired
// - domain.ErrInternal for Redis errors
func (c *CacheService) Exists(ctx context.Context, key string) (bool, error) {
	exists, err := c.client.Exists(ctx, key)
	if err != nil {
		return false, fmt.Errorf("%w: failed to check cache key %s: %v", domain.ErrInternal, key, err)
	}
	return exists, nil
}

// SetNX sets a value only if the key doesn't exist
//
// Use Cases:
// - Distributed locking
// - Race condition prevention
// - Idempotency checks
//
// Returns:
// - true if value was set (key didn't exist)
// - false if key already existed (no-op)
// - domain.ErrInternal for Redis errors
//
// Example - Distributed Lock:
//
//	lockKey := "lock:user:" + userID
//	acquired, err := cache.SetNX(ctx, lockKey, "locked", 10*time.Second)
//	if acquired {
//	    defer cache.Delete(ctx, lockKey)
//	    // Perform locked operation
//	} else {
//	    // Lock already held by another process
//	}
func (c *CacheService) SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error) {
	success, err := c.client.SetNX(ctx, key, value, ttl)
	if err != nil {
		return false, fmt.Errorf("%w: failed to setnx cache key %s: %v", domain.ErrInternal, key, err)
	}
	return success, nil
}

// Compile-time verification that implementation satisfies interface
var _ ports.CacheService = (*CacheService)(nil)
