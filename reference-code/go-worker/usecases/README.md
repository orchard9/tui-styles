# Use Cases - Application Layer

## Overview

This package contains **use cases** (also called **application services** or **interactors**) that represent specific business operations in the system.

Use cases follow the **Vertical Slice Architecture** pattern, where each use case encapsulates a complete business operation from input to output.

## Architecture Position

```
┌─────────────────────────────────────────┐
│         Delivery Layer (HTTP/gRPC)      │  ← Handlers call use cases
└─────────────────────────────────────────┘
                    ↓
┌─────────────────────────────────────────┐
│      Application Layer (Use Cases)      │  ← YOU ARE HERE
│  - Orchestrate domain logic             │
│  - Coordinate between ports             │
│  - Define transaction boundaries        │
└─────────────────────────────────────────┘
                    ↓
┌──────────────────┬──────────────────────┐
│   Domain Layer   │   Ports (Interfaces) │  ← Use cases depend on these
│  - Entities      │   - Repositories     │
│  - Value Objects │   - Cache            │
│  - Services      │   - External APIs    │
└──────────────────┴──────────────────────┘
```

## Vertical Slice Architecture

### What is it?

Instead of organizing code by technical layers (all services together), we organize by **business operations** (vertical slices).

**Traditional Entity-Centric Approach (❌ Avoid):**
```go
// UserService has many responsibilities
type UserService struct {
    repo ports.UserRepository
}

func (s *UserService) CreateUser(...) error { }
func (s *UserService) UpdateUserEmail(...) error { }
func (s *UserService) DeleteUser(...) error { }
func (s *UserService) GetUser(...) error { }
```

**Problems:**
- Violates Single Responsibility Principle
- Hard to test individual operations
- Changes to one operation affect others
- Difficult to scale team (everyone touches UserService)

**Vertical Slice Approach (✅ Recommended):**
```go
// Each use case has ONE responsibility
type CreateUserUseCase struct {
    repo ports.UserRepository
}

type UpdateUserEmailUseCase struct {
    repo ports.UserRepository
}

type DeleteUserUseCase struct {
    repo ports.UserRepository
}
```

**Benefits:**
- Single Responsibility Principle
- Easy to test in isolation
- Changes are localized
- Team can work in parallel
- Clear boundaries for transactions

## Use Case Structure

Every use case follows this pattern:

```go
// 1. Input DTO (Data Transfer Object)
type CreateUserInput struct {
    Email string
    Name  string
}

// 2. Output DTO
type CreateUserOutput struct {
    ID    string
    Email string
    Name  string
}

// 3. Use Case struct
type CreateUserUseCase struct {
    repo  ports.UserRepository
    cache ports.CacheService  // Optional
}

// 4. Constructor (dependency injection)
func NewCreateUserUseCase(repo ports.UserRepository, cache ports.CacheService) *CreateUserUseCase {
    return &CreateUserUseCase{
        repo:  repo,
        cache: cache,
    }
}

// 5. Execute method
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
    // Business logic orchestration
}
```

## Responsibilities

### ✅ Use Cases SHOULD:

1. **Orchestrate** domain logic by calling entities/services
2. **Coordinate** between ports (repository, cache, external APIs)
3. **Define** transaction boundaries
4. **Validate** inputs (basic validation)
5. **Map** between DTOs and domain entities
6. **Handle** cross-cutting concerns (logging, metrics, tracing)
7. **Return** domain errors

### ❌ Use Cases SHOULD NOT:

1. **Contain** business rules (belongs in domain entities/services)
2. **Know** about HTTP, gRPC, or any delivery mechanism
3. **Depend** on concrete implementations (only interfaces)
4. **Perform** complex calculations (belongs in domain)
5. **Handle** infrastructure errors directly (map to domain errors)

## Common Patterns

### 1. Basic CRUD Use Case

```go
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) (*CreateUserOutput, error) {
    // Step 1: Validate input
    if input.Email == "" {
        return nil, domain.ErrInvalidInput("email is required")
    }

    // Step 2: Create domain entity (enforces business rules)
    entity, err := entities.NewUser(input.Email, input.Name)
    if err != nil {
        return nil, err
    }

    // Step 3: Persist via repository
    if err := uc.repo.Create(ctx, entity); err != nil {
        return nil, err
    }

    // Step 4: Return output DTO
    return &CreateUserOutput{
        ID:    entity.ID().String(),
        Email: entity.Email(),
        Name:  entity.Name(),
    }, nil
}
```

### 2. Use Case with Caching

```go
func (uc *GetUserUseCase) Execute(ctx context.Context, input GetUserInput) (*GetUserOutput, error) {
    id, err := uuid.Parse(input.ID)
    if err != nil {
        return nil, domain.ErrInvalidInput("invalid ID format")
    }

    // Try cache first
    if uc.cache != nil {
        var cached entities.User
        if exists, _ := uc.cache.Get(ctx, fmt.Sprintf("user:%s", id), &cached); exists {
            return uc.entityToOutput(&cached), nil
        }
    }

    // Fetch from repository
    entity, err := uc.repo.FindByID(ctx, id)
    if err != nil {
        return nil, err
    }

    // Update cache
    if uc.cache != nil {
        _ = uc.cache.Set(ctx, fmt.Sprintf("user:%s", id), entity, 3600)
    }

    return uc.entityToOutput(entity), nil
}
```

### 3. Use Case with Transaction

```go
type TransferMoneyUseCase struct {
    accountRepo ports.AccountRepository
    txManager   ports.TransactionManager
}

func (uc *TransferMoneyUseCase) Execute(ctx context.Context, input TransferMoneyInput) error {
    // Use transaction manager for atomic operations
    return uc.txManager.WithTransaction(ctx, func(txCtx context.Context) error {
        // All operations use txCtx (transactional context)

        // 1. Debit source account
        sourceAccount, err := uc.accountRepo.FindByID(txCtx, input.SourceID)
        if err != nil {
            return err
        }
        if err := sourceAccount.Debit(input.Amount); err != nil {
            return err
        }
        if err := uc.accountRepo.Update(txCtx, sourceAccount); err != nil {
            return err
        }

        // 2. Credit destination account
        destAccount, err := uc.accountRepo.FindByID(txCtx, input.DestinationID)
        if err != nil {
            return err
        }
        destAccount.Credit(input.Amount)
        if err := uc.accountRepo.Update(txCtx, destAccount); err != nil {
            return err
        }

        return nil // Commit on success
    })
}
```

### 4. Use Case with Multiple Dependencies

```go
type SendWelcomeEmailUseCase struct {
    userRepo     ports.UserRepository
    emailService ports.EmailService
    cache        ports.CacheService
}

func (uc *SendWelcomeEmailUseCase) Execute(ctx context.Context, input SendWelcomeEmailInput) error {
    // Fetch user
    user, err := uc.userRepo.FindByID(ctx, input.UserID)
    if err != nil {
        return err
    }

    // Send email via external service
    if err := uc.emailService.Send(ctx, user.Email(), "Welcome!", "..."); err != nil {
        return err
    }

    // Mark email sent in cache (idempotency)
    _ = uc.cache.Set(ctx, fmt.Sprintf("welcome_sent:%s", input.UserID), true, 86400)

    return nil
}
```

## Testing Use Cases

### Unit Testing

Use mocks for all dependencies:

```go
func TestCreateUserUseCase_Execute(t *testing.T) {
    // Setup
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)
    useCase := usecases.NewCreateUserUseCase(mockRepo, nil)

    // Test successful creation
    input := usecases.CreateUserInput{
        Email: "test@example.com",
        Name:  "Test User",
    }

    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(nil)

    output, err := useCase.Execute(context.Background(), input)
    assert.NoError(t, err)
    assert.NotEmpty(t, output.ID)
    assert.Equal(t, input.Email, output.Email)
}
```

### Integration Testing

Use real dependencies (testcontainers):

```go
func TestCreateUserUseCase_Integration(t *testing.T) {
    // Setup real database with testcontainers
    container := test.StartPostgresContainer(t)
    defer container.Terminate(context.Background())

    db := test.ConnectToDatabase(t, container)
    repo := postgres.NewUserRepository(db)
    useCase := usecases.NewCreateUserUseCase(repo, nil)

    // Test with real database
    input := usecases.CreateUserInput{
        Email: "test@example.com",
        Name:  "Test User",
    }

    output, err := useCase.Execute(context.Background(), input)
    assert.NoError(t, err)

    // Verify in database
    user, err := repo.FindByID(context.Background(), uuid.MustParse(output.ID))
    assert.NoError(t, err)
    assert.Equal(t, input.Email, user.Email())
}
```

## Naming Conventions

### Use Case Names

- **Command:** `CreateUserUseCase`, `UpdateUserEmailUseCase`, `DeleteUserUseCase`
- **Query:** `GetUserUseCase`, `ListUsersUseCase`, `SearchUsersByEmailUseCase`

### Input/Output Names

- **Input:** `CreateUserInput`, `UpdateUserEmailInput`
- **Output:** `CreateUserOutput`, `UpdateUserEmailOutput`

### File Names

- One use case per file
- `create_user.go`, `update_user_email.go`, `get_user.go`

## When to Create a New Use Case

Create a new use case when:

1. **Different Business Operation:** Creating vs Updating vs Deleting
2. **Different Authorization:** Admin-only operations vs user operations
3. **Different Transaction Boundary:** Operations that require different atomicity guarantees
4. **Different Side Effects:** Operations with different external integrations
5. **Different Validation Rules:** Operations with significantly different input validation

## Anti-Patterns to Avoid

### ❌ God Use Case

```go
// Too many responsibilities
type UserUseCase struct {
    // 20 dependencies
}

func (uc *UserUseCase) CreateUser() { }
func (uc *UserUseCase) UpdateUser() { }
func (uc *UserUseCase) DeleteUser() { }
// ... 15 more methods
```

### ❌ Anemic Use Case

```go
// Just a pass-through, no value
func (uc *CreateUserUseCase) Execute(ctx context.Context, input CreateUserInput) error {
    return uc.repo.Create(ctx, &entities.User{...})
}
```

Use cases should orchestrate, not just pass through.

### ❌ Business Logic in Use Case

```go
// Business logic should be in domain, not use case
func (uc *TransferMoneyUseCase) Execute(ctx context.Context, input TransferMoneyInput) error {
    if input.Amount < 0 {  // ❌ Business rule in use case
        return errors.New("amount must be positive")
    }
    // ...
}

// ✅ Business logic in domain entity
func (a *Account) Debit(amount Money) error {
    if amount < 0 {  // ✅ Business rule in entity
        return domain.ErrInvalidAmount("amount must be positive")
    }
    // ...
}
```

## Best Practices

1. **Keep use cases thin** - Orchestrate, don't implement
2. **One responsibility** - Each use case does one thing well
3. **Accept interfaces** - Depend on ports, not concrete implementations
4. **Return DTOs** - Don't expose domain entities to delivery layer
5. **Use context** - Always accept context.Context as first parameter
6. **Handle errors** - Map infrastructure errors to domain errors
7. **Document clearly** - Each use case should have clear documentation
8. **Test thoroughly** - Both unit tests (mocks) and integration tests (real deps)

## Further Reading

- [Hexagonal Architecture](../ARCHITECTURE.md)
- [Domain Layer](../domain/README.md)
- [Ports](../ports/README.md)
- [Testing Guide](../../TESTING.md)
