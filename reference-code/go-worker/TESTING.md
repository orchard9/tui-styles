# Testing Guide

This guide covers testing strategies for hexagonal architecture, including unit tests, integration tests, and end-to-end tests.

## Table of Contents

- [Testing Philosophy](#testing-philosophy)
- [Testing Pyramid](#testing-pyramid)
- [Domain Layer Tests](#domain-layer-tests)
- [Service Layer Tests](#service-layer-tests)
- [Repository Tests](#repository-tests)
- [Handler Tests](#handler-tests)
- [Test Organization](#test-organization)
- [Mocking Strategy](#mocking-strategy)
- [Coverage Requirements](#coverage-requirements)
- [Running Tests](#running-tests)

## Testing Philosophy

### Test at the Right Level

Each layer requires a different testing approach:

| Layer | Test Type | Dependencies | Speed | Coverage Goal |
|-------|-----------|--------------|-------|---------------|
| Domain | Unit | None | Very Fast (<1ms) | 90%+ |
| Service | Unit | Mocked Ports | Fast (<10ms) | 80%+ |
| Repository | Integration | Real DB | Slow (~100ms) | 70%+ |
| Handler | E2E | Real Server | Slowest (~200ms) | Key paths only |

### Key Principles

1. **Test behavior, not implementation**
2. **Write tests before fixing bugs**
3. **Keep tests simple and readable**
4. **Mock at boundaries (ports)**
5. **Integration tests for critical paths**

## Testing Pyramid

```
        ┌─────────┐
        │   E2E   │  Few tests, test critical user journeys
        │ Handler │  Full HTTP stack, real dependencies
        └─────────┘
       ┌───────────┐
       │Integration│  Some tests, verify adapters work
       │Repository │  Real database, SQL correctness
       └───────────┘
      ┌─────────────┐
      │     Unit    │  Many tests, fast feedback
      │Service+Domain│  Mocked dependencies
      └─────────────┘
```

## Domain Layer Tests

Domain tests are **pure unit tests** with **zero dependencies**.

### Entity Tests

**File**: `internal/domain/entities/user_test.go`

```go
package entities_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
)

func TestNewUser_Success(t *testing.T) {
    user, err := entities.NewUser("test@example.com", "John Doe")

    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email().String())
    assert.Equal(t, "John Doe", user.Name())
    assert.NotEqual(t, uuid.Nil, user.ID())
}

func TestNewUser_InvalidEmail(t *testing.T) {
    user, err := entities.NewUser("invalid-email", "John Doe")

    assert.Error(t, err)
    assert.Nil(t, user)
    assert.Contains(t, err.Error(), "invalid email")
}

func TestNewUser_EmptyName(t *testing.T) {
    user, err := entities.NewUser("test@example.com", "")

    assert.Error(t, err)
    assert.Nil(t, user)
}

func TestUser_ChangeEmail(t *testing.T) {
    user, _ := entities.NewUser("old@example.com", "John")

    err := user.ChangeEmail("new@example.com")

    assert.NoError(t, err)
    assert.Equal(t, "new@example.com", user.Email().String())
}

func TestUser_ChangeEmail_Invalid(t *testing.T) {
    user, _ := entities.NewUser("old@example.com", "John")

    err := user.ChangeEmail("invalid")

    assert.Error(t, err)
    assert.Equal(t, "old@example.com", user.Email().String()) // Unchanged
}
```

### Value Object Tests

**File**: `internal/domain/entities/email_test.go`

```go
package entities_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
)

func TestNewEmail_Valid(t *testing.T) {
    tests := []string{
        "test@example.com",
        "user+tag@domain.co.uk",
        "first.last@company.io",
    }

    for _, email := range tests {
        t.Run(email, func(t *testing.T) {
            e, err := entities.NewEmail(email)
            assert.NoError(t, err)
            assert.Equal(t, email, e.String())
        })
    }
}

func TestNewEmail_Invalid(t *testing.T) {
    tests := []string{
        "",
        "notanemail",
        "@example.com",
        "user@",
        "user @example.com",
    }

    for _, email := range tests {
        t.Run(email, func(t *testing.T) {
            _, err := entities.NewEmail(email)
            assert.Error(t, err)
        })
    }
}
```

### Domain Service Tests

**File**: `internal/domain/services/user_service_test.go`

```go
package services_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
    "github.com/orchard9/peach/apps/email-worker/internal/domain/services"
)

func TestUserService_ValidateEmailUniqueness(t *testing.T) {
    service := &services.UserService{}

    user1, _ := entities.NewUser("test@example.com", "User 1")
    user2, _ := entities.NewUser("other@example.com", "User 2")
    existingUsers := []*entities.User{user1, user2}

    // Test duplicate
    err := service.ValidateEmailUniqueness("test@example.com", existingUsers)
    assert.Error(t, err)

    // Test unique
    err = service.ValidateEmailUniqueness("new@example.com", existingUsers)
    assert.NoError(t, err)
}
```

## Service Layer Tests

Service tests use **mocked ports** to test orchestration logic.

### Setting Up Mocks

First, generate mocks:

```bash
go generate ./internal/ports/...
```

This creates `internal/ports/mocks/` with gomock implementations.

### Service Test Example

**File**: `internal/service/user_service_test.go`

```go
package service_test

import (
    "context"
    "errors"
    "testing"

    "github.com/golang/mock/gomock"
    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"

    "github.com/orchard9/peach/apps/email-worker/internal/domain"
    "github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
    "github.com/orchard9/peach/apps/email-worker/internal/ports/mocks"
    "github.com/orchard9/peach/apps/email-worker/internal/service"
)

func TestUserService_CreateUser_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    // Create mocks
    mockRepo := mocks.NewMockUserRepository(ctrl)
    mockCache := mocks.NewMockCacheService(ctrl)

    // Set expectations
    mockRepo.EXPECT().
        FindByEmail(gomock.Any(), "test@example.com").
        Return(nil, domain.ErrNotFound) // Email not taken

    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(nil)

    mockCache.EXPECT().
        DeletePattern(gomock.Any(), "users:*").
        Return(nil)

    // Create service
    svc := service.NewUserService(mockRepo, mockCache, nil)

    // Execute
    user, err := svc.CreateUser(context.Background(), "test@example.com", "Test User")

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, "test@example.com", user.Email().String())
    assert.Equal(t, "Test User", user.Name())
}

func TestUserService_CreateUser_DuplicateEmail(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)
    existingUser, _ := entities.NewUser("test@example.com", "Existing")

    mockRepo.EXPECT().
        FindByEmail(gomock.Any(), "test@example.com").
        Return(existingUser, nil) // Email already exists

    svc := service.NewUserService(mockRepo, nil, nil)

    user, err := svc.CreateUser(context.Background(), "test@example.com", "New User")

    assert.Error(t, err)
    assert.Nil(t, user)
    assert.True(t, errors.Is(err, domain.ErrConflict))
}

func TestUserService_CreateUser_InvalidEmail(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)
    svc := service.NewUserService(mockRepo, nil, nil)

    user, err := svc.CreateUser(context.Background(), "invalid-email", "Test")

    assert.Error(t, err)
    assert.Nil(t, user)
    // No repository calls should be made for invalid input
}

func TestUserService_GetUser_CacheHit(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)
    mockCache := mocks.NewMockCacheService(ctrl)

    userID := uuid.New()
    cachedData := `{"id":"` + userID.String() + `","email":"test@example.com"}`

    mockCache.EXPECT().
        Get(gomock.Any(), "user:"+userID.String()).
        Return(cachedData, nil)

    // Repository should NOT be called on cache hit
    mockRepo.EXPECT().FindByID(gomock.Any(), gomock.Any()).Times(0)

    svc := service.NewUserService(mockRepo, mockCache, nil)
    user, err := svc.GetUser(context.Background(), userID)

    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### Testing with Transactions

**File**: `internal/service/transfer_service_test.go`

```go
func TestTransferService_TransferMoney_Success(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockTxMgr := mocks.NewMockTransactionManager(ctrl)
    mockAccountRepo := mocks.NewMockAccountRepository(ctrl)

    // Expect transaction to be started
    mockTxMgr.EXPECT().
        WithTransaction(gomock.Any(), gomock.Any()).
        DoAndReturn(func(ctx context.Context, fn func(context.Context) error) error {
            // Call the function with transaction context
            return fn(ctx)
        })

    // Expect operations within transaction
    mockAccountRepo.EXPECT().
        Withdraw(gomock.Any(), fromID, amount).
        Return(nil)

    mockAccountRepo.EXPECT().
        Deposit(gomock.Any(), toID, amount).
        Return(nil)

    svc := service.NewTransferService(mockTxMgr, mockAccountRepo)
    err := svc.TransferMoney(ctx, fromID, toID, amount)

    assert.NoError(t, err)
}
```

## Repository Tests

Repository tests are **integration tests** using a real database.

### Test Database Setup

Use `dockertest` or `testcontainers-go` for ephemeral databases:

**File**: `internal/repository/postgres/test_helpers.go`

```go
package postgres_test

import (
    "database/sql"
    "testing"

    "github.com/ory/dockertest/v3"
    "github.com/ory/dockertest/v3/docker"
    _ "github.com/lib/pq"
)

func setupTestDB(t *testing.T) (*sql.DB, func()) {
    pool, err := dockertest.NewPool("")
    if err != nil {
        t.Fatalf("Could not connect to docker: %s", err)
    }

    // Start PostgreSQL container
    resource, err := pool.RunWithOptions(&dockertest.RunOptions{
        Repository: "postgres",
        Tag:        "15-alpine",
        Env: []string{
            "POSTGRES_PASSWORD=test",
            "POSTGRES_DB=testdb",
        },
    }, func(config *docker.HostConfig) {
        config.AutoRemove = true
        config.RestartPolicy = docker.RestartPolicy{Name: "no"}
    })
    if err != nil {
        t.Fatalf("Could not start resource: %s", err)
    }

    var db *sql.DB
    if err := pool.Retry(func() error {
        var err error
        db, err = sql.Open("postgres",
            fmt.Sprintf("postgres://postgres:test@localhost:%s/testdb?sslmode=disable",
                resource.GetPort("5432/tcp")))
        if err != nil {
            return err
        }
        return db.Ping()
    }); err != nil {
        t.Fatalf("Could not connect to database: %s", err)
    }

    // Run migrations
    runMigrations(t, db)

    cleanup := func() {
        db.Close()
        pool.Purge(resource)
    }

    return db, cleanup
}
```

### Repository Test Example

**File**: `internal/repository/postgres/user_test.go`

```go
package postgres_test

import (
    "context"
    "testing"

    "github.com/google/uuid"
    "github.com/stretchr/testify/assert"

    "github.com/orchard9/peach/apps/email-worker/internal/domain"
    "github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
    "github.com/orchard9/peach/apps/email-worker/internal/repository/postgres"
)

func TestPostgresUserRepository_Create(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := postgres.NewUserRepository(db)
    user, _ := entities.NewUser("test@example.com", "Test User")

    err := repo.Create(context.Background(), user)

    assert.NoError(t, err)
}

func TestPostgresUserRepository_FindByID(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := postgres.NewUserRepository(db)
    user, _ := entities.NewUser("test@example.com", "Test User")
    repo.Create(context.Background(), user)

    found, err := repo.FindByID(context.Background(), user.ID())

    assert.NoError(t, err)
    assert.NotNil(t, found)
    assert.Equal(t, user.ID(), found.ID())
    assert.Equal(t, user.Email(), found.Email())
}

func TestPostgresUserRepository_FindByID_NotFound(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := postgres.NewUserRepository(db)
    randomID := uuid.New()

    found, err := repo.FindByID(context.Background(), randomID)

    assert.Error(t, err)
    assert.Nil(t, found)
    assert.True(t, errors.Is(err, domain.ErrNotFound))
}

func TestPostgresUserRepository_Create_DuplicateEmail(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := postgres.NewUserRepository(db)
    user1, _ := entities.NewUser("test@example.com", "User 1")
    user2, _ := entities.NewUser("test@example.com", "User 2")

    err := repo.Create(context.Background(), user1)
    assert.NoError(t, err)

    err = repo.Create(context.Background(), user2)
    assert.Error(t, err)
    assert.True(t, errors.Is(err, domain.ErrConflict))
}

func TestPostgresUserRepository_Update(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := postgres.NewUserRepository(db)
    user, _ := entities.NewUser("test@example.com", "Original Name")
    repo.Create(context.Background(), user)

    // Update user
    user.ChangeEmail("new@example.com")
    err := repo.Update(context.Background(), user)
    assert.NoError(t, err)

    // Verify update
    found, _ := repo.FindByID(context.Background(), user.ID())
    assert.Equal(t, "new@example.com", found.Email().String())
}
```

## Handler Tests

Handler tests are **end-to-end tests** using a test HTTP server.

**File**: `internal/handlers/user_handler_test.go`

```go
package handlers_test

import (
    "bytes"
    "encoding/json"
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi/v5"
    "github.com/stretchr/testify/assert"

    "github.com/orchard9/peach/apps/email-worker/internal/handlers"
    "github.com/orchard9/peach/apps/email-worker/internal/service"
    "github.com/orchard9/peach/apps/email-worker/internal/repository/memory"
)

func setupTestServer(t *testing.T) *chi.Mux {
    // Use in-memory implementations for testing
    userRepo := memory.NewUserRepository()
    cacheService := memory.NewCacheService()
    userService := service.NewUserService(userRepo, cacheService, nil)
    userHandler := handlers.NewUserHandler(userService)

    r := chi.NewRouter()
    r.Post("/users", userHandler.CreateUser)
    r.Get("/users/{id}", userHandler.GetUser)

    return r
}

func TestUserHandler_CreateUser_Success(t *testing.T) {
    router := setupTestServer(t)

    reqBody := map[string]string{
        "email": "test@example.com",
        "name":  "Test User",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusCreated, rr.Code)

    var resp map[string]interface{}
    json.NewDecoder(rr.Body).Decode(&resp)
    assert.Equal(t, "test@example.com", resp["email"])
    assert.Equal(t, "Test User", resp["name"])
    assert.NotEmpty(t, resp["id"])
}

func TestUserHandler_CreateUser_InvalidEmail(t *testing.T) {
    router := setupTestServer(t)

    reqBody := map[string]string{
        "email": "invalid-email",
        "name":  "Test User",
    }
    body, _ := json.Marshal(reqBody)

    req := httptest.NewRequest("POST", "/users", bytes.NewReader(body))
    req.Header.Set("Content-Type", "application/json")
    rr := httptest.NewRecorder()

    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHandler_GetUser_NotFound(t *testing.T) {
    router := setupTestServer(t)

    req := httptest.NewRequest("GET", "/users/"+uuid.New().String(), nil)
    rr := httptest.NewRecorder()

    router.ServeHTTP(rr, req)

    assert.Equal(t, http.StatusNotFound, rr.Code)
}
```

## Test Organization

### File Naming

- Domain tests: `*_test.go` in same package
- Service tests: `*_test.go` in `_test` package
- Integration tests: `test/integration/*_test.go`

### Package Naming

```go
// Domain tests - same package (white box)
package entities

func TestNewUser(t *testing.T) { ... }

// Service tests - separate package (black box)
package service_test

func TestUserService_CreateUser(t *testing.T) { ... }
```

## Mocking Strategy

### When to Mock

✅ **Mock**:
- External dependencies (database, cache, APIs)
- At port boundaries
- In service layer tests

❌ **Don't Mock**:
- Domain entities
- Value objects
- Pure functions

### Mock Alternatives

For testing, you can also create in-memory implementations:

**File**: `internal/repository/memory/user.go`

```go
package memory

type InMemoryUserRepository struct {
    users map[uuid.UUID]*entities.User
    mu    sync.RWMutex
}

func (r *InMemoryUserRepository) Create(ctx context.Context, user *entities.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.users[user.ID()]; exists {
        return domain.ErrConflict
    }
    r.users[user.ID()] = user
    return nil
}
```

Use in-memory for handler tests and local development.

## Coverage Requirements

### Minimum Coverage

- **Domain layer**: 90%+
- **Service layer**: 80%+
- **Repository layer**: 70%+
- **Handlers**: Critical paths only

### Checking Coverage

```bash
# Run tests with coverage
make test-coverage

# View HTML report
open coverage.html

# Check coverage percentage
go tool cover -func=coverage.out | grep total
```

### Coverage Best Practices

- Don't chase 100% coverage
- Focus on business logic
- Cover error paths
- Test critical user journeys

## Running Tests

### All Tests

```bash
make test
```

### Unit Tests Only

```bash
go test -short ./...
```

### Integration Tests

```bash
go test -run Integration ./test/integration/...
```

### With Coverage

```bash
make test-coverage
```

### Watch Mode

```bash
make test-watch
```

### Specific Package

```bash
go test ./internal/domain/entities/...
```

### Verbose Output

```bash
go test -v ./...
```

## Common Testing Patterns

### Table-Driven Tests

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid email", "test@example.com", false},
        {"invalid email", "invalid", true},
        {"empty email", "", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            _, err := NewEmail(tt.input)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }
        })
    }
}
```

### Test Fixtures

```go
func createTestUser(t *testing.T) *entities.User {
    user, err := entities.NewUser("test@example.com", "Test User")
    assert.NoError(t, err)
    return user
}
```

### Cleanup with defer

```go
func TestWithCleanup(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    // Test logic...
}
```

## Best Practices

1. **Test one thing per test**
2. **Use descriptive test names**
3. **Follow AAA pattern** (Arrange, Act, Assert)
4. **Don't test implementation details**
5. **Make tests readable**
6. **Keep tests fast**
7. **Avoid test interdependencies**
8. **Use helper functions**
9. **Clean up resources**
10. **Test error cases**

## Further Reading

- [ARCHITECTURE.md](./ARCHITECTURE.md) - Understand the layers
- [DEVELOPMENT.md](./DEVELOPMENT.md) - Development workflows
- [testify documentation](https://github.com/stretchr/testify)
- [gomock documentation](https://github.com/golang/mock)
