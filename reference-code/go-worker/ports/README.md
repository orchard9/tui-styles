# Ports Package

## Overview

This package defines **ports** (interfaces) in the hexagonal architecture. Ports are contracts between layers that enable:
- **Dependency Inversion**: Application depends on abstractions, not concrete implementations
- **Testability**: Easy mocking for unit tests
- **Flexibility**: Swap implementations without changing business logic

## What are Ports?

Ports are **interfaces** that define how the application layer communicates with:
- **Inbound**: External requests (HTTP, gRPC) → Application
- **Outbound**: Application → External services (Database, Cache, APIs)

```
┌─────────────────────────────────────────┐
│         Adapters (Infrastructure)       │
│   HTTP, gRPC, PostgreSQL, Redis, SMTP   │
└──────────────┬──────────────────────────┘
               │ implements
┌──────────────▼──────────────────────────┐
│          Ports (Interfaces)             │ ← THIS PACKAGE
│  Repository, Cache, Email, Storage      │
└──────────────┬──────────────────────────┘
               │ used by
┌──────────────▼──────────────────────────┐
│      Application (Use Cases)            │
│   CreateUser, GetUser, etc.             │
└──────────────┬──────────────────────────┘
               │ calls
┌──────────────▼──────────────────────────┐
│       Domain (Core Business Logic)      │
│   Entities, Value Objects, Services     │
└─────────────────────────────────────────┘
```

## Port Files

### repositories.go
Defines repository interfaces for data persistence.

**Example**:
```go
type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
    // ...
}
```

**Implementations**:
- `adapters/repository/postgres/user.go` - PostgreSQL
- `adapters/repository/memory/user.go` - In-memory (testing)

### cache.go
Defines caching interface.

**Example**:
```go
type CacheService interface {
    Get(ctx context.Context, key string) (string, error)
    Set(ctx context.Context, key, value string, ttl time.Duration) error
    // ...
}
```

**Implementations**:
- `adapters/cache/redis/cache.go` - Redis
- `adapters/cache/memory/cache.go` - In-memory (testing)

### transaction.go
Defines transaction management interface.

**Example**:
```go
type TransactionManager interface {
    WithTransaction(ctx context.Context, fn func(txCtx context.Context) error) error
}
```

**Implementations**:
- `adapters/repository/postgres/transaction.go` - PostgreSQL transactions

### external.go
Defines external service interfaces.

**Example**:
```go
type EmailService interface {
    Send(ctx context.Context, to, subject, body string) error
}

type StorageService interface {
    Upload(ctx context.Context, key string, content []byte) error
    Download(ctx context.Context, key string) ([]byte, error)
}
```

**Implementations**:
- `adapters/external/email/smtp.go` - SMTP email
- `adapters/external/storage/s3.go` - AWS S3 storage

## Design Principles

### 1. Interface Segregation

Keep interfaces small and focused:

```go
// ✅ Good: Focused interface
type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
}

// ❌ Bad: Too many responsibilities
type Repository interface {
    CreateUser(...)
    CreatePost(...)
    CreateComment(...)
    // Too much!
}
```

### 2. Context First

Always pass `context.Context` as first parameter:

```go
// ✅ Good
FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)

// ❌ Bad
FindByID(id uuid.UUID, ctx context.Context) (*entities.User, error)
```

### 3. Domain Types Only

Ports reference domain types, never infrastructure types:

```go
// ✅ Good: Domain entity
Create(ctx context.Context, user *entities.User) error

// ❌ Bad: Database type
Create(ctx context.Context, user *sql.Row) error
```

### 4. Domain Errors

Return domain errors, not infrastructure errors:

```go
// ✅ Good: Domain error
return domain.ErrNotFound

// ❌ Bad: Infrastructure error
return sql.ErrNoRows
```

## Usage Examples

### In Use Cases

Inject ports as dependencies:

```go
type CreateUserUseCase struct {
    userRepo  ports.UserRepository  // Interface!
    cache     ports.CacheService    // Interface!
    emailSvc  ports.EmailService    // Interface!
}

func (uc *CreateUserUseCase) Execute(ctx context.Context, req dto.CreateUserRequest) error {
    // Create domain entity
    user, err := entities.NewUser(req.Email, req.Name)
    if err != nil {
        return err
    }

    // Persist via repository port
    if err := uc.userRepo.Create(ctx, user); err != nil {
        return err
    }

    // Invalidate cache via cache port
    if err := uc.cache.DeletePattern(ctx, "users:*"); err != nil {
        // Log but don't fail
    }

    // Send welcome email via email port
    if err := uc.emailSvc.Send(ctx, user.Email(), "Welcome!", "..."); err != nil {
        // Log but don't fail
    }

    return nil
}
```

### In Tests

Mock ports for unit testing:

```go
func TestCreateUserUseCase_Success(t *testing.T) {
    // Setup mocks
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)
    mockCache := mocks.NewMockCacheService(ctrl)
    mockEmail := mocks.NewMockEmailService(ctrl)

    // Set expectations
    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(nil).
        Times(1)

    mockCache.EXPECT().
        DeletePattern(gomock.Any(), "users:*").
        Return(nil).
        Times(1)

    mockEmail.EXPECT().
        Send(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
        Return(nil).
        Times(1)

    // Create use case with mocks
    uc := usecases.NewCreateUserUseCase(mockRepo, mockCache, mockEmail)

    // Execute
    err := uc.Execute(context.Background(), dto.CreateUserRequest{
        Email: "test@example.com",
        Name:  "Test User",
    })

    // Assert
    assert.NoError(t, err)
}
```

### In Dependency Injection

Wire concrete implementations to ports:

```go
// main.go
func main() {
    // Create adapters (concrete implementations)
    userRepo := postgres.NewUserRepository(db)      // PostgreSQL adapter
    cacheService := redis.NewCacheService(redisClient)  // Redis adapter
    emailService := smtp.NewEmailService(smtpConfig)    // SMTP adapter

    // Create use cases (inject as interfaces)
    createUserUC := usecases.NewCreateUserUseCase(
        userRepo,      // ports.UserRepository
        cacheService,  // ports.CacheService
        emailService,  // ports.EmailService
    )

    // Create handlers
    userHandler := handlers.NewUserHandler(createUserUC)
    // ...
}
```

## Adding New Ports

### 1. Define Interface

```go
// internal/ports/notification.go
package ports

type NotificationService interface {
    SendPush(ctx context.Context, userID uuid.UUID, message string) error
}
```

### 2. Add Mock Generation

```go
//go:generate mockgen -source=notification.go -destination=mocks/notification.go -package=mocks
```

### 3. Create Adapter

```go
// adapters/external/notification/firebase.go
package notification

type FirebaseNotificationService struct {
    client *firebase.Client
}

func (s *FirebaseNotificationService) SendPush(ctx context.Context, userID uuid.UUID, message string) error {
    // Implementation
}

// Compile-time verification
var _ ports.NotificationService = (*FirebaseNotificationService)(nil)
```

### 4. Wire in main.go

```go
notificationSvc := notification.NewFirebaseService(firebaseClient)
createUserUC := usecases.NewCreateUserUseCase(
    userRepo,
    cacheService,
    emailService,
    notificationSvc,  // New port!
)
```

### 5. Generate Mocks

```bash
go generate ./internal/ports/...
```

## Benefits

### Testability

```go
// Easy to test - just mock the ports
mockRepo := mocks.NewMockUserRepository(ctrl)
uc := usecases.NewCreateUserUseCase(mockRepo)
```

### Flexibility

```go
// Development: In-memory
userRepo := memory.NewUserRepository()

// Production: PostgreSQL
userRepo := postgres.NewUserRepository(db)

// Same use case works with both!
uc := usecases.NewCreateUserUseCase(userRepo)
```

### Clear Boundaries

```go
// Ports make layer boundaries explicit
Application Layer → uses ports → defined in ports package
Adapters → implement ports → in adapters package
```

## Common Questions

**Q: Should I create a port for everything?**

A: Only for dependencies that cross layer boundaries:
- ✅ YES: Repositories, cache, external services
- ❌ NO: Domain services, entity methods

**Q: Can use cases call other use cases?**

A: Generally avoid. Instead:
- Extract common logic to domain service
- Or create a coordinating use case that orchestrates multiple operations

**Q: How do I handle transactions?**

A: Use the `TransactionManager` port:

```go
err := txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
    // All repository calls use txCtx
    if err := userRepo.Create(txCtx, user); err != nil {
        return err
    }
    return profileRepo.Create(txCtx, profile)
})
```

**Q: Should ports return DTOs or domain entities?**

A: Domain entities! Ports are part of the application/domain boundary, not the HTTP boundary.

## See Also

- `../domain/` - Domain layer (entities, value objects, services)
- `../application/` - Application layer (use cases, DTOs)
- `../../adapters/` - Adapter implementations of ports
