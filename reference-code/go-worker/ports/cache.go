package ports

import (
	"context"
	"time"
)

//go:generate mockgen -source=cache.go -destination=mocks/cache.go -package=mocks

// CacheService defines the contract for caching operations.
//
// This is a PORT (interface) in hexagonal architecture.
// Implementations are ADAPTERS (e.g., Redis, in-memory).
//
// Design Principles:
// - Generic key-value operations
// - TTL support for expiration
// - Pattern-based deletion
// - No assumptions about cache backend
//
// Thread Safety:
// - Implementations MUST be safe for concurrent use
//
// Example Implementations:
// - adapters/cache/redis/cache.go
// - adapters/cache/memory/cache.go
type CacheService interface {
	// Get retrieves a value from cache.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - key: Cache key
	//
	// Returns:
	// - string: Cached value
	// - error: domain.ErrNotFound if key doesn't exist
	// - error: domain.ErrInternal for cache failures
	//
	// Example:
	//   value, err := cache.Get(ctx, "user:123")
	//   if errors.Is(err, domain.ErrNotFound) {
	//       // Cache miss - fetch from database
	//   }
	Get(ctx context.Context, key string) (string, error)

	// Set stores a value in cache with expiration.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - key: Cache key
	// - value: Value to store
	// - ttl: Time-to-live (0 = no expiration)
	//
	// Returns:
	// - error: domain.ErrInternal for cache failures
	//
	// Example:
	//   err := cache.Set(ctx, "user:123", userData, 5*time.Minute)
	Set(ctx context.Context, key string, value string, ttl time.Duration) error

	// Delete removes a value from cache.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - key: Cache key to delete
	//
	// Returns:
	// - error: domain.ErrInternal for cache failures
	//
	// Note: Deleting non-existent key is not an error
	//
	// Example:
	//   err := cache.Delete(ctx, "user:123")
	Delete(ctx context.Context, key string) error

	// DeletePattern removes all keys matching a pattern.
	//
	// Useful for cache invalidation when data changes.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - pattern: Pattern to match (e.g., "users:*")
	//
	// Returns:
	// - error: domain.ErrInternal for cache failures
	//
	// Pattern Syntax:
	// - "*" matches any characters
	// - "?" matches single character
	// - Example: "users:*" matches "users:123", "users:456"
	//
	// Example:
	//   // Invalidate all user caches
	//   err := cache.DeletePattern(ctx, "users:*")
	DeletePattern(ctx context.Context, pattern string) error

	// Exists checks if a key exists in cache.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - key: Cache key to check
	//
	// Returns:
	// - bool: true if key exists
	// - error: domain.ErrInternal for cache failures
	//
	// Example:
	//   exists, err := cache.Exists(ctx, "user:123")
	//   if !exists {
	//       // Populate cache
	//   }
	Exists(ctx context.Context, key string) (bool, error)

	// SetNX sets a value only if the key doesn't exist.
	//
	// Useful for distributed locking or preventing race conditions.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - key: Cache key
	// - value: Value to store
	// - ttl: Time-to-live
	//
	// Returns:
	// - bool: true if value was set, false if key already existed
	// - error: domain.ErrInternal for cache failures
	//
	// Example:
	//   // Distributed lock
	//   acquired, err := cache.SetNX(ctx, "lock:user:123", "locked", 10*time.Second)
	//   if acquired {
	//       defer cache.Delete(ctx, "lock:user:123")
	//       // Perform locked operation
	//   }
	SetNX(ctx context.Context, key string, value string, ttl time.Duration) (bool, error)
}
