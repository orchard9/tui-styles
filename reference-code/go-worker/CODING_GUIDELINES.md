# /var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker Coding Guidelines

**Version**: 1.0.0  
**Last Updated**: 2024-07-19  
**Toolkit Version**: go-core-http-toolkit v0.6.3

## Overview

This document defines the mandatory coding standards for the /var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker project. All code MUST strictly adhere to these guidelines to ensure consistency, maintainability, and compliance with the go-core-http-toolkit framework.

## Core Philosophy

The project follows the **opinionated architecture** of go-core-http-toolkit v0.6.3:
- **One way to do things** - No abstractions for multiple implementations
- **Simplicity over flexibility** - Use only the prescribed tools
- **Convention over configuration** - Follow toolkit patterns exactly

## 1. Logging Standards

### MANDATORY: Use slog Only

**❌ FORBIDDEN:**
```go
import "log"

log.Printf("Starting server on %s", addr)
log.Println("Server started")
log.Fatal("Critical error")
```

**✅ REQUIRED:**
```go
import "log/slog"

logger := slog.Default()
logger.Info("Starting server", slog.String("address", addr))
logger.Info("Server started")
logger.Error("Critical error", slog.String("error", err.Error()))
```

### Structured Logging Requirements

1. **Always use structured fields**:
```go
logger.Info("User authenticated",
    slog.String("user_id", userID),
    slog.String("method", "JWT"),
    slog.Duration("latency", latency))
```

2. **Log levels**:
   - `Debug`: Development/debugging information
   - `Info`: Normal operational messages
   - `Warn`: Warning conditions
   - `Error`: Error conditions (but service continues)

3. **Context propagation**:
```go
logger := slog.Default().With(
    slog.String("request_id", requestID),
    slog.String("user_id", userID))
```

## 2. Configuration Management

### MANDATORY: Use Viper Only

**❌ FORBIDDEN:**
```go
// Direct environment variable access
dbHost := os.Getenv("DB_HOST")
port := getEnv("PORT", "8080")

// Custom config structs
type Config struct {
    Port string
}
```

**✅ REQUIRED:**
```go
import "github.com/spf13/viper"

// Initialize Viper
viper.SetConfigName("config")
viper.SetConfigType("yaml")
viper.AddConfigPath(".")
viper.AutomaticEnv()

// Read configuration
dbHost := viper.GetString("database.host")
port := viper.GetInt("server.port")
```

### Configuration Structure

```yaml
# config.yaml
server:
  port: 8080
  timeout: 30s

database:
  host: localhost
  port: 5432
  name: /var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker_dev

```

## 3. HTTP & Routing

### MANDATORY: Chi Router Only

**❌ FORBIDDEN:**
```go
// No Gin, Fiber, Echo, or Gorilla Mux
import "github.com/gin-gonic/gin"
import "github.com/gofiber/fiber/v2"
import "github.com/labstack/echo/v4"
import "github.com/gorilla/mux"
```

**✅ REQUIRED:**
```go
import "github.com/go-chi/chi/v5"
import "github.com/go-chi/chi/v5/middleware"

r := chi.NewRouter()

// Mandatory middleware order
r.Use(middleware.RequestID)
r.Use(middleware.RealIP)
r.Use(middleware.Logger)
r.Use(middleware.Recoverer)
r.Use(middleware.Timeout(30 * time.Second))

// Routes with /api/v1 prefix
r.Route("/api/v1", func(r chi.Router) {
    r.Route("/users", func(r chi.Router) {
        r.Get("/", h.ListUsers)
        r.Post("/", h.CreateUser)
        r.Get("/{id}", h.GetUser)
    })
})
```

## 4. Authentication & Authorization

### MANDATORY: JWT with HMAC-SHA256 Only

**❌ FORBIDDEN:**
```go
// Multiple algorithm support
switch algorithm {
case "HS256":
    signingMethod = jwt.SigningMethodHS256
case "RS256":
    signingMethod = jwt.SigningMethodRS256
}
```

**✅ REQUIRED:**
```go
import "github.com/golang-jwt/jwt/v5"

// ONLY HMAC-SHA256
signingMethod := jwt.SigningMethodHS256

// Fixed implementation
func GenerateToken(claims jwt.MapClaims, secret string) (string, error) {
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(secret))
}
```

## 5. Database Standards

### MANDATORY: PostgreSQL with sqlx Only

**❌ FORBIDDEN:**
```go
// No ORMs
import "gorm.io/gorm"

// No other SQL libraries
import "database/sql"
```

**✅ REQUIRED:**
```go
import (
    "github.com/jmoiron/sqlx"
    _ "github.com/lib/pq"
)

// Connection
db, err := sqlx.Connect("postgres", dsn)

// Queries with sqlx
var users []User
err := db.Select(&users, "SELECT * FROM users WHERE active = $1", true)

// Transactions
tx, err := db.Beginx()
defer tx.Rollback()
// ... operations
tx.Commit()
```

### SQL Guidelines

1. **Always use parameterized queries**:
```go
// ✅ CORRECT
db.Get(&user, "SELECT * FROM users WHERE id = $1", id)

// ❌ WRONG - SQL injection risk
db.Get(&user, fmt.Sprintf("SELECT * FROM users WHERE id = %s", id))
```

2. **Use sqlx tags for struct mapping**:
```go
type User struct {
    ID        string    `db:"id"`
    Name      string    `db:"name"`
    CreatedAt time.Time `db:"created_at"`
}
```

## 6. Caching Standards

### MANDATORY: Redis Only

**✅ REQUIRED:**
```go
import "github.com/redis/go-redis/v9"

// Client initialization
rdb := redis.NewClient(&redis.Options{
    Addr:     "localhost:6379",
    Password: "",
    DB:       0,
})

// Usage pattern with TTL
ctx := context.Background()
err := rdb.Set(ctx, "key", "value", 5*time.Minute).Err()
val, err := rdb.Get(ctx, "key").Result()
```

### Caching Patterns

1. **Cache keys must be namespaced**:
```go
key := fmt.Sprintf("/var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker:users:%s", userID)
```

2. **Always set TTL**:
```go
rdb.Set(ctx, key, value, 15*time.Minute)
```

## 7. Error Handling

### Standard Error Pattern

Use the generated errors package:

```go
import "github.com/orchard9/peach/apps/email-worker/internal/errors"

// Return domain errors
if user == nil {
    return errors.NotFound("user")
}

// Wrap errors with context
if err != nil {
    return errors.WrapError(err, "failed to create user")
}

// Check error types
if errors.Is(err, errors.ErrNotFound) {
    // Handle specific error
}
```

### HTTP Error Responses

```go
// Use base handler helpers
func (h *BaseHandler) respondError(w http.ResponseWriter, err error) {
    switch {
    case errors.Is(err, errors.ErrNotFound):
        h.RespondJSON(w, http.StatusNotFound, map[string]string{"error": "Resource not found"})
    case errors.Is(err, errors.ErrInvalidInput):
        h.RespondJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
    default:
        h.RespondJSON(w, http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
    }
}
```

## 8. Testing Standards

### MANDATORY: Testify Framework

**✅ REQUIRED:**
```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
    "github.com/stretchr/testify/suite"
)

func TestUserCreation(t *testing.T) {
    // Use assert for non-critical checks
    assert.Equal(t, expected, actual)
    
    // Use require for critical checks
    require.NoError(t, err)
    require.NotNil(t, user)
}
```

### Test Organization

1. **Table-driven tests**:
```go
func TestValidateEmail(t *testing.T) {
    tests := []struct {
        name    string
        email   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"invalid email", "not-an-email", true},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := ValidateEmail(tt.email)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

## 9. Code Organization

### Project Structure

```
/var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker/
├── cmd/                    # Application entrypoints
│   ├── server/            # Main HTTP server
├── internal/              # Private application code
│   ├── config/           # Configuration management
│   ├── errors/           # Error definitions
│   ├── handlers/         # HTTP handlers
│   ├── repository/       # Data access
│   └── service/          # Business logic
├── migrations/          # Database migrations
├── test/               # Test utilities
└── scripts/            # Development scripts
```

### Import Ordering

```go
import (
    // Standard library
    "context"
    "fmt"
    "time"
    
    // Third-party
    "github.com/go-chi/chi/v5"
    "github.com/jmoiron/sqlx"
    
    // Internal
    "github.com/orchard9/peach/apps/email-worker/internal/config"
    "github.com/orchard9/peach/apps/email-worker/internal/errors"
    "github.com/orchard9/peach/apps/email-worker/internal/service"
)
```

## 10. Security Guidelines

### Required Security Practices

1. **Input Validation**:
```go
import "github.com/go-playground/validator/v10"

validate := validator.New()
err := validate.Struct(input)
```

2. **SQL Injection Prevention**:
   - Always use parameterized queries
   - Never concatenate user input into SQL

3. **Security Headers**:
```go
// Add security middleware in server setup
r.Use(middleware.SetHeader("X-Content-Type-Options", "nosniff"))
r.Use(middleware.SetHeader("X-Frame-Options", "DENY"))
r.Use(middleware.SetHeader("X-XSS-Protection", "1; mode=block"))
```

## Enforcement

### Pre-commit Checks

All code must pass:
```bash
# Format
go fmt ./...

# Vet
go vet ./...

# Lint
golangci-lint run

# Tests
go test ./... -v

# CI check
make ci
```

### Code Review Checklist

- [ ] Uses slog for all logging
- [ ] Uses Viper for configuration
- [ ] Uses Chi router exclusively
- [ ] JWT uses only HMAC-SHA256
- [ ] Uses sqlx for database access
- [ ] Uses generated errors package
- [ ] Uses testify for tests
- [ ] Follows security guidelines
- [ ] Passes all linting checks

## References

- [go-core-http-toolkit Documentation](https://github.com/orchard9/go-core-http-toolkit)
- [Chi Router Documentation](https://github.com/go-chi/chi)
- [sqlx Documentation](https://github.com/jmoiron/sqlx)
- [Viper Documentation](https://github.com/spf13/viper)
