# Caching Strategy Guide

This guide outlines best practices for implementing and using caching in <no value>, following hexagonal architecture principles and production-ready patterns.

## Table of Contents

- [Architecture Overview](#architecture-overview)
- [Cache Key Conventions](#cache-key-conventions)
- [TTL Configuration](#ttl-configuration)
- [Invalidation Strategies](#invalidation-strategies)
- [Error Handling](#error-handling)
- [Cache Stampede Prevention](#cache-stampede-prevention)
- [Testing Patterns](#testing-patterns)
- [Best Practices](#best-practices)
- [Anti-Patterns](#anti-patterns)

## Architecture Overview

### Hexagonal Architecture Alignment

Caching in <no value> follows strict hexagonal architecture boundaries:

```
Domain (Core)
    ↑
    | Port Interface (ports/cache.go)
    ↓
Adapter (Infrastructure)
    ├── Redis Implementation
    ├── In-Memory Implementation
    └── Composite/Multi-Layer
```

**Key Principles:**
- Domain logic NEVER imports cache implementations directly
- Use cases depend on cache ports (interfaces), not concrete implementations
- Cache failures should not break domain operations (fail open)
- Cache is an optimization, not a requirement for correctness

### When to Use Caching

**Good Use Cases:**
- Expensive database queries that change infrequently
- External API responses with predictable staleness tolerance
- Computed aggregations or reports
- Session data or user preferences
- Rate limiting counters

**Poor Use Cases:**
- Data that changes frequently (>1 update/second)
- Data requiring strong consistency guarantees
- Large objects that exceed network/memory efficiency thresholds
- Security-sensitive data without encryption

## Cache Key Conventions

### Naming Pattern

Follow this hierarchical naming convention:

```
{service}:{domain}:{entity}:{identifier}[:version]
```

**Examples:**
```go
// Single entity by ID
"<no value>:<no value>:<no value>:550e8400-e29b-41d4-a716-446655440000"

// Entity with version
"<no value>:<no value>:<no value>:550e8400-e29b-41d4-a716-446655440000:v2"

// List/collection cache
"<no value>:<no value>:<no value>s:list:active:page1"

// Aggregation
"<no value>:<no value>:stats:daily:2025-10-26"

// User-scoped data
"<no value>:<no value>:user:123:preferences"
```

### Key Generation Helper

Create a key builder in your adapter for consistency:

```go
// internal/adapter/repository/cache_keys.go
package repository

import (
    "fmt"
    "strings"
)

const (
    servicePrefix = "<no value>"
    domainPrefix  = "<no value>"
)

type CacheKeyBuilder struct {
    parts []string
}

func NewCacheKey() *CacheKeyBuilder {
    return &CacheKeyBuilder{
        parts: []string{servicePrefix, domainPrefix},
    }
}

func (b *CacheKeyBuilder) Entity(name string) *CacheKeyBuilder {
    b.parts = append(b.parts, name)
    return b
}

func (b *CacheKeyBuilder) ID(id string) *CacheKeyBuilder {
    b.parts = append(b.parts, id)
    return b
}

func (b *CacheKeyBuilder) Scope(scopes ...string) *CacheKeyBuilder {
    b.parts = append(b.parts, scopes...)
    return b
}

func (b *CacheKeyBuilder) Build() string {
    return strings.Join(b.parts, ":")
}

// Usage examples:
// key := NewCacheKey().Entity("<no value>").ID("123").Build()
// key := NewCacheKey().Entity("<no value>s").Scope("list", "active").Build()
```

### Key Pattern Prefixes

Use consistent patterns for different access patterns:

```go
const (
    // Single entity retrieval
    entityByIDPattern = "{service}:{domain}:{entity}:{id}"

    // List/query results
    entityListPattern = "{service}:{domain}:{entity}s:list:{filter}:{page}"

    // Aggregations
    aggregationPattern = "{service}:{domain}:{entity}s:{aggregation}:{period}"

    // User-scoped data
    userScopedPattern = "{service}:{domain}:user:{userID}:{resource}"
)
```

## TTL Configuration

### TTL Selection Guidelines

Choose TTL based on data characteristics and business requirements:

| Data Type | Recommended TTL | Rationale |
|-----------|----------------|-----------|
| Immutable entities | 24 hours - 7 days | Data never changes once created |
| Frequently read, rare writes | 5-60 minutes | Balance staleness vs database load |
| User sessions | 15-30 minutes | Security vs UX trade-off |
| External API responses | Per API SLA | Respect upstream rate limits |
| Computed aggregations | 1-6 hours | Expensive to compute, acceptable staleness |
| Rate limit counters | 1 minute - 1 hour | Match rate limit window |
| Real-time data | 10-60 seconds | Near real-time requirements |

### Configuration Pattern

Define TTLs as constants with clear documentation:

```go
// internal/adapter/repository/cache_ttls.go
package repository

import "time"

const (
    // <no value>TTL is the cache duration for individual <no value> entities.
    // Rationale: <no value>s change infrequently and are read-heavy.
    <no value>TTL = 30 * time.Minute

    // <no value>ListTTL is the cache duration for <no value> list queries.
    // Rationale: Lists are more volatile than individual entities due to creation/deletion.
    <no value>ListTTL = 5 * time.Minute

    // <no value>StatsTTL is the cache duration for <no value> aggregations.
    // Rationale: Stats are expensive to compute and tolerate staleness.
    <no value>StatsTTL = 1 * time.Hour

    // SessionTTL is the cache duration for user sessions.
    // Rationale: Balance between security (shorter is better) and UX (longer is better).
    SessionTTL = 15 * time.Minute
)
```

### Dynamic TTL Calculation

For time-sensitive data, calculate TTL based on content:

```go
func calculateTTL(createdAt time.Time) time.Duration {
    age := time.Since(createdAt)

    switch {
    case age < 1*time.Hour:
        return 5 * time.Minute   // Fresh data, cache briefly
    case age < 24*time.Hour:
        return 30 * time.Minute  // Recent data, medium cache
    default:
        return 2 * time.Hour     // Old data, longer cache
    }
}
```

## Invalidation Strategies

### Strategy Selection Matrix

| Strategy | Use Case | Pros | Cons |
|----------|----------|------|------|
| **TTL-only** | Acceptable staleness window | Simple, no coordination | Stale reads until expiry |
| **Explicit invalidation** | Immediate consistency needed | Precise control | Coupling between write and cache |
| **Pattern-based** | Related entities (cascade) | Handles relationships | Potentially broad invalidation |
| **Write-through** | Strong consistency | Always up-to-date | Higher write latency |
| **Event-driven** | Distributed systems | Decoupled, scalable | Added complexity |

### Explicit Invalidation

Invalidate specific keys on write operations:

```go
// internal/adapter/repository/postgres_<no value>.go
func (r *<no value>Repository) Update(
    ctx context.Context,
    <no value> *domain.<no value>,
) error {
    // Perform database update
    if err := r.db.Update<no value>(ctx, <no value>); err != nil {
        return fmt.Errorf("update <no value>: %w", err)
    }

    // Invalidate specific cache entry
    cacheKey := NewCacheKey().
        Entity("<no value>").
        ID(<no value>.ID.String()).
        Build()

    if err := r.cache.Delete(ctx, cacheKey); err != nil {
        // Log but don't fail - cache invalidation is best-effort
        r.logger.Warn("failed to invalidate cache",
            "key", cacheKey,
            "error", err,
        )
    }

    return nil
}
```

### Pattern-Based Invalidation

Invalidate multiple related keys using patterns:

```go
func (r *<no value>Repository) Delete(
    ctx context.Context,
    id uuid.UUID,
) error {
    if err := r.db.Delete<no value>(ctx, id); err != nil {
        return fmt.Errorf("delete <no value>: %w", err)
    }

    // Invalidate specific entity
    entityKey := NewCacheKey().Entity("<no value>").ID(id.String()).Build()
    _ = r.cache.Delete(ctx, entityKey)

    // Invalidate all list caches (pattern-based)
    listPattern := NewCacheKey().Entity("<no value>s").Scope("list").Build() + ":*"
    if err := r.cache.DeletePattern(ctx, listPattern); err != nil {
        r.logger.Warn("failed to invalidate list caches",
            "pattern", listPattern,
            "error", err,
        )
    }

    return nil
}
```

### Cache-Aside Pattern (Lazy Loading)

Most common pattern for reads:

```go
func (r *<no value>Repository) GetByID(
    ctx context.Context,
    id uuid.UUID,
) (*domain.<no value>, error) {
    cacheKey := NewCacheKey().Entity("<no value>").ID(id.String()).Build()

    // 1. Try cache first
    var <no value> domain.<no value>
    err := r.cache.Get(ctx, cacheKey, &<no value>)
    if err == nil {
        // Cache hit
        return &<no value>, nil
    }

    // 2. Cache miss - fetch from database
    <no value>Ptr, err := r.db.Find<no value>ByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("find <no value>: %w", err)
    }

    // 3. Populate cache (fire-and-forget to avoid blocking)
    go func() {
        setCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
        defer cancel()

        if err := r.cache.Set(setCtx, cacheKey, <no value>Ptr, <no value>TTL); err != nil {
            r.logger.Warn("failed to cache <no value>",
                "id", id,
                "error", err,
            )
        }
    }()

    return <no value>Ptr, nil
}
```

## Error Handling

### Fail Open Philosophy

**Golden Rule**: Cache failures should NEVER break application functionality.

```go
// Good: Fail open on cache errors
func (r *<no value>Repository) GetByID(
    ctx context.Context,
    id uuid.UUID,
) (*domain.<no value>, error) {
    cacheKey := NewCacheKey().Entity("<no value>").ID(id.String()).Build()

    var <no value> domain.<no value>
    err := r.cache.Get(ctx, cacheKey, &<no value>)
    if err == nil {
        return &<no value>, nil
    }

    // Log cache miss/error but continue to database
    if err != ports.ErrCacheMiss {
        r.logger.Warn("cache get failed, falling back to database",
            "key", cacheKey,
            "error", err,
        )
    }

    // Proceed to authoritative data source
    return r.db.Find<no value>ByID(ctx, id)
}
```

### Circuit Breaker Pattern

Protect against cascading cache failures:

```go
// internal/adapter/cache/circuit_breaker.go
type CircuitBreakerCache struct {
    cache   ports.CacheService
    breaker *CircuitBreaker
    logger  *slog.Logger
}

func (c *CircuitBreakerCache) Get(
    ctx context.Context,
    key string,
    dest interface{},
) error {
    if c.breaker.IsOpen() {
        return ports.ErrCacheMiss // Treat as cache miss when circuit is open
    }

    err := c.cache.Get(ctx, key, dest)
    if err != nil {
        c.breaker.RecordFailure()
        return err
    }

    c.breaker.RecordSuccess()
    return nil
}
```

### Timeout Configuration

Always set aggressive timeouts for cache operations:

```go
const (
    // CacheGetTimeout is the maximum time to wait for a cache read.
    // Rationale: Cache should be fast; slow cache defeats the purpose.
    CacheGetTimeout = 50 * time.Millisecond

    // CacheSetTimeout is the maximum time to wait for a cache write.
    // Rationale: Sets are often async; don't block on them.
    CacheSetTimeout = 100 * time.Millisecond

    // CacheDeleteTimeout is the maximum time to wait for cache invalidation.
    // Rationale: Invalidation is best-effort; don't block requests.
    CacheDeleteTimeout = 100 * time.Millisecond
)

func (r *<no value>Repository) getWithTimeout(
    ctx context.Context,
    key string,
    dest interface{},
) error {
    ctx, cancel := context.WithTimeout(ctx, CacheGetTimeout)
    defer cancel()

    return r.cache.Get(ctx, key, dest)
}
```

## Cache Stampede Prevention

### Problem Description

Cache stampede occurs when:
1. Popular cache entry expires
2. Multiple requests simultaneously detect cache miss
3. All requests query database concurrently
4. Database experiences spike in load

### Single-Flight Pattern

Use single-flight to coalesce concurrent requests:

```go
// internal/adapter/repository/singleflight_<no value>.go
package repository

import (
    "context"
    "fmt"
    "golang.org/x/sync/singleflight"
    "github.com/google/uuid"
)

type <no value>Repository struct {
    db         *PostgresDB
    cache      ports.CacheService
    logger     *slog.Logger
    sfGroup    singleflight.Group  // Single-flight group
}

func (r *<no value>Repository) GetByID(
    ctx context.Context,
    id uuid.UUID,
) (*domain.<no value>, error) {
    cacheKey := NewCacheKey().Entity("<no value>").ID(id.String()).Build()

    // Try cache first
    var <no value> domain.<no value>
    err := r.cache.Get(ctx, cacheKey, &<no value>)
    if err == nil {
        return &<no value>, nil
    }

    // Use single-flight to prevent stampede
    result, err, shared := r.sfGroup.Do(cacheKey, func() (interface{}, error) {
        // Only one goroutine executes this for the same key
        <no value>Ptr, err := r.db.Find<no value>ByID(ctx, id)
        if err != nil {
            return nil, err
        }

        // Populate cache asynchronously
        go r.setCacheAsync(cacheKey, <no value>Ptr)

        return <no value>Ptr, nil
    })

    if err != nil {
        return nil, fmt.Errorf("get <no value>: %w", err)
    }

    if shared {
        r.logger.Debug("request coalesced via single-flight",
            "key", cacheKey,
        )
    }

    return result.(*domain.<no value>), nil
}

func (r *<no value>Repository) setCacheAsync(key string, value interface{}) {
    ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
    defer cancel()

    if err := r.cache.Set(ctx, key, value, <no value>TTL); err != nil {
        r.logger.Warn("async cache set failed",
            "key", key,
            "error", err,
        )
    }
}
```

### Probabilistic Early Expiration

Prevent thundering herd by randomly refreshing before expiration:

```go
func (r *<no value>Repository) shouldRefreshEarly(ttl time.Duration) bool {
    // Refresh early with probability that increases as TTL decreases
    // Example: 10% chance when 90% of TTL remains, 90% chance when 10% remains
    remainingRatio := rand.Float64()
    threshold := 1.0 - (ttl.Seconds() / <no value>TTL.Seconds())

    return remainingRatio < threshold
}
```

## Testing Patterns

### Mocking Cache Port

Create test doubles for cache port:

```go
// internal/ports/mocks/cache_mock.go
package mocks

import (
    "context"
    "sync"
)

type CacheMock struct {
    mu      sync.RWMutex
    storage map[string]interface{}

    GetFunc    func(ctx context.Context, key string, dest interface{}) error
    SetFunc    func(ctx context.Context, key string, value interface{}, ttl time.Duration) error
    DeleteFunc func(ctx context.Context, key string) error
}

func NewCacheMock() *CacheMock {
    return &CacheMock{
        storage: make(map[string]interface{}),
    }
}

func (m *CacheMock) Get(ctx context.Context, key string, dest interface{}) error {
    if m.GetFunc != nil {
        return m.GetFunc(ctx, key, dest)
    }

    m.mu.RLock()
    defer m.mu.RUnlock()

    val, exists := m.storage[key]
    if !exists {
        return ports.ErrCacheMiss
    }

    // Simple reflection-based copy for testing
    reflect.ValueOf(dest).Elem().Set(reflect.ValueOf(val).Elem())
    return nil
}
```

### Testing Cache Hit/Miss Scenarios

```go
func TestGetByID_CacheHit(t *testing.T) {
    // Arrange
    cacheMock := mocks.NewCacheMock()
    dbMock := mocks.NewDBMock()
    repo := NewRepository(dbMock, cacheMock)

    expected<no value> := &domain.<no value>{
        ID: uuid.New(),
        // ... other fields
    }

    cacheKey := NewCacheKey().Entity("<no value>").ID(expected<no value>.ID.String()).Build()
    cacheMock.storage[cacheKey] = expected<no value>

    // Act
    result, err := repo.GetByID(context.Background(), expected<no value>.ID)

    // Assert
    require.NoError(t, err)
    assert.Equal(t, expected<no value>, result)
    assert.Zero(t, dbMock.FindByIDCallCount) // Database should not be called
}

func TestGetByID_CacheMiss(t *testing.T) {
    // Arrange
    cacheMock := mocks.NewCacheMock()
    dbMock := mocks.NewDBMock()
    repo := NewRepository(dbMock, cacheMock)

    expected<no value> := &domain.<no value>{
        ID: uuid.New(),
    }

    dbMock.Find<no value>ByIDFunc = func(ctx context.Context, id uuid.UUID) (*domain.<no value>, error) {
        return expected<no value>, nil
    }

    // Act
    result, err := repo.GetByID(context.Background(), expected<no value>.ID)

    // Assert
    require.NoError(t, err)
    assert.Equal(t, expected<no value>, result)
    assert.Equal(t, 1, dbMock.FindByIDCallCount) // Database should be called once
}
```

### Testing Invalidation

```go
func TestUpdate_InvalidatesCache(t *testing.T) {
    // Arrange
    cacheMock := mocks.NewCacheMock()
    dbMock := mocks.NewDBMock()
    repo := NewRepository(dbMock, cacheMock)

    <no value> := &domain.<no value>{
        ID: uuid.New(),
    }

    cacheKey := NewCacheKey().Entity("<no value>").ID(<no value>.ID.String()).Build()
    cacheMock.storage[cacheKey] = <no value>

    var deletedKey string
    cacheMock.DeleteFunc = func(ctx context.Context, key string) error {
        deletedKey = key
        delete(cacheMock.storage, key)
        return nil
    }

    // Act
    err := repo.Update(context.Background(), <no value>)

    // Assert
    require.NoError(t, err)
    assert.Equal(t, cacheKey, deletedKey)
    assert.NotContains(t, cacheMock.storage, cacheKey)
}
```

## Best Practices

### 1. Cache Only Serializable Data

```go
// Good: Simple, serializable domain model
type <no value> struct {
    ID        uuid.UUID
    Name      string
    CreatedAt time.Time
}

// Bad: Contains channels, mutexes, or other non-serializable types
type <no value>Bad struct {
    ID        uuid.UUID
    Name      string
    stateChan chan State  // Cannot be serialized
    mu        sync.Mutex  // Cannot be serialized
}
```

### 2. Version Cache Keys for Schema Changes

```go
const (
    <no value>CacheVersion = "v2" // Increment when <no value> structure changes
)

func (b *CacheKeyBuilder) Version(v string) *CacheKeyBuilder {
    b.parts = append(b.parts, v)
    return b
}

// Usage:
key := NewCacheKey().
    Entity("<no value>").
    ID(id.String()).
    Version(<no value>CacheVersion).
    Build()
```

### 3. Monitor Cache Performance

```go
func (r *<no value>Repository) GetByID(
    ctx context.Context,
    id uuid.UUID,
) (*domain.<no value>, error) {
    start := time.Now()
    defer func() {
        duration := time.Since(start)
        r.metrics.RecordCacheOperation("get", duration)
    }()

    // ... cache logic
}
```

### 4. Use Compression for Large Objects

```go
func (r *<no value>Repository) cacheWithCompression(
    ctx context.Context,
    key string,
    value interface{},
    ttl time.Duration,
) error {
    data, err := json.Marshal(value)
    if err != nil {
        return err
    }

    // Compress if larger than threshold
    if len(data) > 1024 { // 1KB threshold
        compressed := compress(data)
        return r.cache.Set(ctx, key+":gz", compressed, ttl)
    }

    return r.cache.Set(ctx, key, value, ttl)
}
```

### 5. Document Cache Behavior

```go
// GetByID retrieves a <no value> by ID.
//
// Caching behavior:
//   - Cache key: <no value>:<no value>:<no value>:{id}
//   - TTL: 30 minutes
//   - Invalidation: On Update() and Delete()
//   - Failure mode: Falls back to database on cache errors
func (r *<no value>Repository) GetByID(
    ctx context.Context,
    id uuid.UUID,
) (*domain.<no value>, error) {
    // ...
}
```

## Anti-Patterns

### 1. DON'T Cache Without TTL

```go
// Bad: No expiration
r.cache.Set(ctx, key, value, 0)

// Good: Always set appropriate TTL
r.cache.Set(ctx, key, value, <no value>TTL)
```

### 2. DON'T Fail on Cache Errors

```go
// Bad: Returns error when cache fails
if err := r.cache.Get(ctx, key, &result); err != nil {
    return nil, err // Breaks on cache failure
}

// Good: Fail open
if err := r.cache.Get(ctx, key, &result); err == nil {
    return &result, nil
}
// Continue to database
```

### 3. DON'T Cache Domain Logic

```go
// Bad: Caching computed business logic
func (s *<no value>Service) CalculatePrice(ctx context.Context, id uuid.UUID) (float64, error) {
    cacheKey := fmt.Sprintf("price:%s", id)
    var price float64
    if err := s.cache.Get(ctx, cacheKey, &price); err == nil {
        return price, nil
    }

    // Business logic should not be in adapter layer
    price = s.calculateComplexPrice(id)
    s.cache.Set(ctx, cacheKey, price, time.Hour)
    return price, nil
}

// Good: Cache repository data, compute in domain
func (s *<no value>Service) CalculatePrice(ctx context.Context, id uuid.UUID) (float64, error) {
    // Repository handles caching of entity
    <no value>, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return 0, err
    }

    // Domain logic in service
    return <no value>.CalculatePrice(), nil
}
```

### 4. DON'T Over-Invalidate

```go
// Bad: Nuclear invalidation
func (r *<no value>Repository) Update(ctx context.Context, <no value> *domain.<no value>) error {
    r.cache.DeletePattern(ctx, "<no value>:*") // Clears EVERYTHING
    return r.db.Update(ctx, <no value>)
}

// Good: Targeted invalidation
func (r *<no value>Repository) Update(ctx context.Context, <no value> *domain.<no value>) error {
    entityKey := NewCacheKey().Entity("<no value>").ID(<no value>.ID.String()).Build()
    r.cache.Delete(ctx, entityKey)

    listPattern := NewCacheKey().Entity("<no value>s").Scope("list").Build() + ":*"
    r.cache.DeletePattern(ctx, listPattern)

    return r.db.Update(ctx, <no value>)
}
```

### 5. DON'T Store Pointers in Cache

```go
// Bad: Storing pointer can cause race conditions
var <no value> *domain.<no value>
r.cache.Set(ctx, key, <no value>, ttl) // Stores reference

// Good: Store value or make defensive copy
<no value>Copy := *<no value>
r.cache.Set(ctx, key, <no value>Copy, ttl)
```

---

## Summary

Effective caching in <no value> requires:

1. **Architectural discipline**: Keep cache adapters separate from domain logic
2. **Defensive coding**: Always fail open on cache errors
3. **Clear conventions**: Use consistent key naming and TTL strategies
4. **Proper invalidation**: Balance consistency with cache efficiency
5. **Comprehensive testing**: Test hit/miss scenarios and error conditions
6. **Monitoring**: Track cache hit rates and performance metrics

Remember: **Cache is an optimization, not a crutch**. Your application should work correctly without cache; cache makes it faster.

For specific implementation examples, see:
- `/internal/ports/cache.go` - Cache port interface
- `/internal/adapter/cache/redis.go` - Redis cache adapter
- `/internal/adapter/repository/postgres_<no value>.go` - Repository with caching
