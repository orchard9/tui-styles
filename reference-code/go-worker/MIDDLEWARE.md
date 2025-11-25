# Middleware Guide

This guide covers middleware usage in your /var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker service, including built-in toolkit middleware, custom implementations, configuration, and testing strategies.

## Table of Contents

- [Overview](#overview)
- [Built-in Middleware](#built-in-middleware)
- [Custom Middleware](#custom-middleware)
- [Configuration](#configuration)
- [Middleware Chaining](#middleware-chaining)
- [Testing Middleware](#testing-middleware)
- [Best Practices](#best-practices)
- [Common Patterns](#common-patterns)

## Overview

### What is Middleware?

Middleware is a function that wraps HTTP handlers to add cross-cutting concerns like:
- Request logging
- Authentication/authorization
- Request ID tracking
- Error recovery
- Rate limiting
- CORS handling
- Metrics collection

### Middleware Signature

In Go, HTTP middleware follows this pattern:

```go
func Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Before request processing

        next.ServeHTTP(w, r)  // Call next handler

        // After request processing
    })
}
```

### Execution Order

Middleware executes in the order it's registered:

```
Request
   ↓
Middleware 1 (before)
   ↓
Middleware 2 (before)
   ↓
Middleware 3 (before)
   ↓
Handler
   ↓
Middleware 3 (after)
   ↓
Middleware 2 (after)
   ↓
Middleware 1 (after)
   ↓
Response
```

## Built-in Middleware

The `go-core-http-toolkit` provides production-ready middleware for common concerns.

### Request ID Middleware

Adds a unique request ID to each request for tracing.

**Usage:**

```go
import (
    "github.com/orchard9/go-core-http-toolkit/pkg/middleware"
)

func setupRouter() *chi.Mux {
    r := chi.NewRouter()

    // Add request ID middleware
    r.Use(middleware.RequestID)

    return r
}
```

**Features:**
- Generates UUID v4 for each request
- Adds `X-Request-ID` header to request and response
- Uses existing `X-Request-ID` from request if present
- Stores request ID in context for handler access

**Accessing Request ID in Handlers:**

```go
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    requestID := middleware.GetRequestID(r.Context())

    log.Printf("[%s] Creating user", requestID)

    // Rest of handler logic...
}
```

### Logger Middleware

Structured request/response logging with duration, status, and request details.

**Usage:**

```go
import (
    "github.com/orchard9/go-core-http-toolkit/pkg/middleware"
    "go.uber.org/zap"
)

func setupRouter(logger *zap.Logger) *chi.Mux {
    r := chi.NewRouter()

    // Add structured logging
    r.Use(middleware.Logger(logger))

    return r
}
```

**Log Output:**

```json
{
  "level": "info",
  "ts": "2025-10-26T10:30:45.123Z",
  "msg": "HTTP request",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "method": "POST",
  "path": "/api/v1/users",
  "status": 201,
  "duration_ms": 45.2,
  "size_bytes": 256,
  "user_agent": "curl/7.64.1",
  "remote_addr": "192.168.1.100:52345"
}
```

**Features:**
- Structured JSON logging
- Request/response metrics
- Duration tracking
- Error logging
- Request ID correlation
- User agent and IP tracking

### Recovery Middleware

Recovers from panics and returns 500 errors gracefully.

**Usage:**

```go
import (
    "github.com/orchard9/go-core-http-toolkit/pkg/middleware"
    "go.uber.org/zap"
)

func setupRouter(logger *zap.Logger) *chi.Mux {
    r := chi.NewRouter()

    // Add panic recovery (should be first!)
    r.Use(middleware.Recoverer(logger))

    return r
}
```

**Features:**
- Catches panics in handlers
- Logs panic with stack trace
- Returns 500 Internal Server Error
- Prevents server crashes
- Includes request ID in error logs

**Example Panic Handling:**

```go
func (h *Handler) DangerousEndpoint(w http.ResponseWriter, r *http.Request) {
    // If this panics, Recoverer middleware catches it
    result := someDangerousOperation()

    response.Success(w, result)
}
```

### CORS Middleware

Handles Cross-Origin Resource Sharing for browser-based clients.

**Usage:**

```go
import (
    "github.com/orchard9/go-core-http-toolkit/pkg/middleware"
)

func setupRouter() *chi.Mux {
    r := chi.NewRouter()

    // Add CORS support
    r.Use(middleware.CORS(middleware.CORSConfig{
        AllowedOrigins:   []string{"https://app.example.com", "https://admin.example.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
        ExposedHeaders:   []string{"X-Request-ID"},
        AllowCredentials: true,
        MaxAge:           3600, // 1 hour
    }))

    return r
}
```

**Development Configuration:**

```go
// Allow all origins in development
r.Use(middleware.CORS(middleware.CORSConfig{
    AllowedOrigins: []string{"*"},
    AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
    AllowedHeaders: []string{"*"},
}))
```

**Production Configuration:**

```go
// Strict origin control in production
r.Use(middleware.CORS(middleware.CORSConfig{
    AllowedOrigins:   []string{os.Getenv("ALLOWED_ORIGIN")},
    AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
    AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
    AllowCredentials: true,
}))
```

### Rate Limiting Middleware

Protects endpoints from abuse with request rate limits.

**Usage:**

```go
import (
    "github.com/orchard9/go-core-http-toolkit/pkg/middleware"
    "time"
)

func setupRouter() *chi.Mux {
    r := chi.NewRouter()

    // Global rate limit: 100 requests per minute per IP
    r.Use(middleware.RateLimit(middleware.RateLimitConfig{
        RequestsPerWindow: 100,
        WindowDuration:    time.Minute,
        KeyFunc: middleware.IPBasedKey,
    }))

    return r
}
```

**Per-Route Rate Limiting:**

```go
r := chi.NewRouter()

// Strict rate limit for login endpoint
loginRateLimit := middleware.RateLimit(middleware.RateLimitConfig{
    RequestsPerWindow: 5,
    WindowDuration:    time.Minute,
    KeyFunc: middleware.IPBasedKey,
})

r.With(loginRateLimit).Post("/auth/login", authHandler.Login)

// Normal rate limit for other endpoints
r.Post("/users", userHandler.CreateUser)
```

**User-Based Rate Limiting:**

```go
userRateLimit := middleware.RateLimit(middleware.RateLimitConfig{
    RequestsPerWindow: 1000,
    WindowDuration:    time.Hour,
    KeyFunc: func(r *http.Request) string {
        // Rate limit by authenticated user ID
        claims := middleware.GetJWTClaims(r.Context())
        if userID, ok := claims["user_id"].(string); ok {
            return userID
        }
        // Fallback to IP for unauthenticated requests
        return middleware.IPBasedKey(r)
    },
})
```

### Timeout Middleware

Enforces request timeout limits to prevent long-running requests.

**Usage:**

```go
import (
    "github.com/orchard9/go-core-http-toolkit/pkg/middleware"
    "time"
)

func setupRouter() *chi.Mux {
    r := chi.NewRouter()

    // Global 30-second timeout
    r.Use(middleware.Timeout(30 * time.Second))

    return r
}
```

**Per-Route Timeouts:**

```go
// Long timeout for file uploads
uploadTimeout := middleware.Timeout(5 * time.Minute)
r.With(uploadTimeout).Post("/files/upload", fileHandler.Upload)

// Short timeout for health checks
healthTimeout := middleware.Timeout(5 * time.Second)
r.With(healthTimeout).Get("/health", healthHandler.Check)
```

## Custom Middleware

### Creating Custom Middleware

Custom middleware follows the standard Go HTTP middleware pattern.

**Basic Custom Middleware:**

```go
// internal/middleware/custom.go
package middleware

import (
    "context"
    "net/http"
)

// TenantID extracts tenant ID from subdomain or header
func TenantID(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Extract tenant ID from header
        tenantID := r.Header.Get("X-Tenant-ID")

        if tenantID == "" {
            http.Error(w, "missing tenant ID", http.StatusBadRequest)
            return
        }

        // Add to context
        ctx := context.WithValue(r.Context(), "tenant_id", tenantID)

        // Call next handler with updated context
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

**Using Custom Middleware:**

```go
import "github.com/orchard9/peach/apps/email-worker/internal/middleware"

func setupRouter() *chi.Mux {
    r := chi.NewRouter()

    r.Use(middleware.TenantID)

    return r
}
```

**Accessing Context Values:**

```go
func (h *UserHandler) ListUsers(w http.ResponseWriter, r *http.Request) {
    tenantID := r.Context().Value("tenant_id").(string)

    users, err := h.userService.ListUsersByTenant(r.Context(), tenantID)
    // ...
}
```

### Advanced Custom Middleware Examples

**Permission Checking Middleware:**

```go
// RequirePermission returns middleware that checks user permissions
func RequirePermission(permission string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims := GetJWTClaims(r.Context())

            // Extract permissions from JWT claims
            permissions, ok := claims["permissions"].([]interface{})
            if !ok {
                http.Error(w, "unauthorized", http.StatusForbidden)
                return
            }

            // Check if user has required permission
            hasPermission := false
            for _, p := range permissions {
                if p.(string) == permission {
                    hasPermission = true
                    break
                }
            }

            if !hasPermission {
                http.Error(w, "insufficient permissions", http.StatusForbidden)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

**Usage:**

```go
// Only admin users can delete
r.With(RequirePermission("users:delete")).Delete("/users/{id}", userHandler.Delete)

// Any authenticated user can read
r.With(RequirePermission("users:read")).Get("/users/{id}", userHandler.Get)
```

**Request Validation Middleware:**

```go
func ValidateContentType(contentType string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ct := r.Header.Get("Content-Type")

            if ct != contentType {
                http.Error(w,
                    fmt.Sprintf("invalid content type: expected %s, got %s", contentType, ct),
                    http.StatusUnsupportedMediaType)
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}
```

**Usage:**

```go
// Require JSON for POST/PUT endpoints
r.With(ValidateContentType("application/json")).Post("/users", userHandler.Create)
r.With(ValidateContentType("application/json")).Put("/users/{id}", userHandler.Update)
```

**Audit Logging Middleware:**

```go
func AuditLog(logger *zap.Logger) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Capture response status
            wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

            next.ServeHTTP(wrapped, r)

            // Log after request completes
            if r.Method != "GET" && r.Method != "HEAD" {
                claims := GetJWTClaims(r.Context())
                userID := claims["user_id"].(string)

                logger.Info("audit_log",
                    zap.String("user_id", userID),
                    zap.String("method", r.Method),
                    zap.String("path", r.URL.Path),
                    zap.Int("status", wrapped.statusCode),
                    zap.String("request_id", GetRequestID(r.Context())),
                )
            }
        })
    }
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}
```

## Configuration

### Environment-Based Configuration

Configure middleware based on environment variables.

**Example Configuration:**

```go
// internal/config/middleware.go
package config

import (
    "os"
    "strconv"
    "time"
)

type MiddlewareConfig struct {
    EnableCORS      bool
    AllowedOrigins  []string
    RateLimitRPS    int
    RequestTimeout  time.Duration
    JWTSecret       string
}

func LoadMiddlewareConfig() *MiddlewareConfig {
    return &MiddlewareConfig{
        EnableCORS:     os.Getenv("ENABLE_CORS") == "true",
        AllowedOrigins: strings.Split(os.Getenv("ALLOWED_ORIGINS"), ","),
        RateLimitRPS:   getEnvInt("RATE_LIMIT_RPS", 100),
        RequestTimeout: time.Duration(getEnvInt("REQUEST_TIMEOUT_SECONDS", 30)) * time.Second,
        JWTSecret:      os.Getenv("JWT_SECRET"),
    }
}

func getEnvInt(key string, defaultValue int) int {
    if val := os.Getenv(key); val != "" {
        if i, err := strconv.Atoi(val); err == nil {
            return i
        }
    }
    return defaultValue
}
```

**Using Configuration:**

```go
func setupRouter(cfg *config.MiddlewareConfig, logger *zap.Logger) *chi.Mux {
    r := chi.NewRouter()

    r.Use(middleware.Recoverer(logger))
    r.Use(middleware.RequestID)
    r.Use(middleware.Logger(logger))
    r.Use(middleware.Timeout(cfg.RequestTimeout))

    if cfg.EnableCORS {
        r.Use(middleware.CORS(middleware.CORSConfig{
            AllowedOrigins: cfg.AllowedOrigins,
        }))
    }

    r.Use(middleware.RateLimit(middleware.RateLimitConfig{
        RequestsPerWindow: cfg.RateLimitRPS,
        WindowDuration:    time.Minute,
    }))

    return r
}
```

## Middleware Chaining

### Global Middleware

Applied to all routes.

```go
r := chi.NewRouter()

// Global middleware (order matters!)
r.Use(middleware.Recoverer(logger))      // 1. Catch panics first
r.Use(middleware.RequestID)               // 2. Generate request IDs
r.Use(middleware.Logger(logger))          // 3. Log requests
r.Use(middleware.CORS(corsConfig))        // 4. Handle CORS
r.Use(middleware.RateLimit(rateLimitCfg)) // 5. Rate limiting
r.Use(middleware.Timeout(30*time.Second)) // 6. Request timeouts
```

### Route Groups

Apply middleware to specific route groups.

```go
r := chi.NewRouter()

// Public routes - no auth
r.Group(func(r chi.Router) {
    r.Post("/auth/login", authHandler.Login)
    r.Get("/health", healthHandler.Check)
})

// Protected routes - require JWT
r.Group(func(r chi.Router) {
    r.Use(jwtMiddleware)

    r.Get("/users/me", userHandler.GetMe)
    r.Put("/users/me", userHandler.UpdateMe)
})

// Admin routes - require JWT + admin permission
r.Group(func(r chi.Router) {
    r.Use(jwtMiddleware)
    r.Use(RequirePermission("admin"))

    r.Get("/admin/users", adminHandler.ListUsers)
    r.Delete("/admin/users/{id}", adminHandler.DeleteUser)
})
```

### Per-Route Middleware

Apply middleware to individual routes.

```go
r := chi.NewRouter()

// Most routes have normal rate limit
r.Post("/users", userHandler.Create)

// Login has strict rate limit
strictRateLimit := middleware.RateLimit(middleware.RateLimitConfig{
    RequestsPerWindow: 5,
    WindowDuration:    time.Minute,
})
r.With(strictRateLimit).Post("/auth/login", authHandler.Login)

// File upload has longer timeout
longTimeout := middleware.Timeout(5 * time.Minute)
r.With(longTimeout).Post("/files", fileHandler.Upload)
```

## Testing Middleware

### Testing Custom Middleware

**Example Test:**

```go
// internal/middleware/tenant_test.go
package middleware_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/stretchr/testify/assert"
    "github.com/orchard9/peach/apps/email-worker/internal/middleware"
)

func TestTenantID_Success(t *testing.T) {
    // Create test handler
    var capturedTenantID string
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        capturedTenantID = r.Context().Value("tenant_id").(string)
        w.WriteHeader(http.StatusOK)
    })

    // Wrap with middleware
    wrapped := middleware.TenantID(handler)

    // Create request
    req := httptest.NewRequest("GET", "/test", nil)
    req.Header.Set("X-Tenant-ID", "tenant-123")
    rr := httptest.NewRecorder()

    // Execute
    wrapped.ServeHTTP(rr, req)

    // Assert
    assert.Equal(t, http.StatusOK, rr.Code)
    assert.Equal(t, "tenant-123", capturedTenantID)
}

func TestTenantID_MissingHeader(t *testing.T) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        t.Error("handler should not be called")
    })

    wrapped := middleware.TenantID(handler)

    req := httptest.NewRequest("GET", "/test", nil)
    // No X-Tenant-ID header
    rr := httptest.NewRecorder()

    wrapped.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusBadRequest, rr.Code)
}
```

### Testing Handlers with Middleware

**Integration Test:**

```go
func TestUserHandler_WithAuth(t *testing.T) {
    // Setup
    router := setupTestRouter()

    // Create valid JWT token
    token := createTestJWT(t, map[string]interface{}{
        "user_id": "user-123",
        "role":    "user",
    })

    // Create request with auth header
    req := httptest.NewRequest("GET", "/users/me", nil)
    req.Header.Set("Authorization", "Bearer "+token)
    rr := httptest.NewRecorder()

    // Execute
    router.ServeHTTP(rr, req)

    // Assert
    assert.Equal(t, http.StatusOK, rr.Code)
}
```

## Best Practices

### 1. Order Middleware Carefully

```go
// ✅ Good: Recovery first, logging early, auth after basics
r.Use(middleware.Recoverer(logger))
r.Use(middleware.RequestID)
r.Use(middleware.Logger(logger))
r.Use(middleware.CORS(corsConfig))
r.Use(jwtMiddleware)

// ❌ Bad: Auth before recovery could miss panics
r.Use(jwtMiddleware)
r.Use(middleware.Recoverer(logger))
```

### 2. Use Context for Request-Scoped Data

```go
// ✅ Good: Use context for passing data
ctx := context.WithValue(r.Context(), "user_id", userID)
next.ServeHTTP(w, r.WithContext(ctx))

// ❌ Bad: Don't use global variables
globalUserID = userID  // Race condition!
```

### 3. Don't Write Response Before Calling Next

```go
// ❌ Bad: Writing headers breaks middleware chain
func BadMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Custom", "value")
        w.WriteHeader(http.StatusOK)  // Don't do this!
        next.ServeHTTP(w, r)
    })
}

// ✅ Good: Set headers, let handler write response
func GoodMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Custom", "value")
        next.ServeHTTP(w, r)  // Handler will write response
    })
}
```

### 4. Keep Middleware Focused

```go
// ✅ Good: Single responsibility
func RequestID(next http.Handler) http.Handler { /* ... */ }
func Logger(next http.Handler) http.Handler { /* ... */ }
func Auth(next http.Handler) http.Handler { /* ... */ }

// ❌ Bad: Too many responsibilities
func DoEverything(next http.Handler) http.Handler {
    // Logs, authenticates, rate limits, etc.
}
```

### 5. Make Middleware Configurable

```go
// ✅ Good: Configurable middleware
func RateLimit(config RateLimitConfig) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        // Use config...
    }
}

// ❌ Bad: Hard-coded values
func RateLimit(next http.Handler) http.Handler {
    limit := 100  // Hard-coded!
}
```

## Common Patterns

### Pattern 1: Response Writer Wrapper

Capture response status and size:

```go
type responseWriter struct {
    http.ResponseWriter
    status int
    size   int
}

func (rw *responseWriter) WriteHeader(status int) {
    rw.status = status
    rw.ResponseWriter.WriteHeader(status)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
    size, err := rw.ResponseWriter.Write(b)
    rw.size += size
    return size, err
}
```

### Pattern 2: Middleware with Dependencies

```go
type AuthMiddleware struct {
    userService ports.UserService
    logger      *zap.Logger
}

func NewAuthMiddleware(userService ports.UserService, logger *zap.Logger) *AuthMiddleware {
    return &AuthMiddleware{
        userService: userService,
        logger:      logger,
    }
}

func (m *AuthMiddleware) Middleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Can use m.userService and m.logger
        next.ServeHTTP(w, r)
    })
}
```

### Pattern 3: Conditional Middleware

```go
func ConditionalMiddleware(condition func(*http.Request) bool, middleware func(http.Handler) http.Handler) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            if condition(r) {
                middleware(next).ServeHTTP(w, r)
            } else {
                next.ServeHTTP(w, r)
            }
        })
    }
}
```

## Further Reading

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Project architecture overview
- [DEVELOPMENT.md](./DEVELOPMENT.md) - Development guide
- [go-core-http-toolkit Middleware Documentation](https://github.com/orchard9/go-core-http-toolkit)
- [Chi Router Middleware](https://go-chi.io/#/pages/middleware)
