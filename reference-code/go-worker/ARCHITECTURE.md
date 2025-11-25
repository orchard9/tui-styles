# Architecture Overview

## Hexagonal Architecture (Ports & Adapters)

This project follows **Hexagonal Architecture** to maintain clean separation between business logic and infrastructure concerns.

### Core Principle

> **Business logic should not depend on infrastructure details**

The domain layer contains pure business logic with zero external dependencies. All interactions with external systems (databases, APIs, caching) happen through **interfaces (ports)** implemented by **adapters**.

## Layer Structure

```
┌──────────────────────────────────────────────────────────────┐
│                     HTTP/gRPC Handlers                        │
│                    (Inbound Adapters)                         │
│   internal/handlers/*.go                                      │
└──────────────────────┬───────────────────────────────────────┘
                       │ calls
                       ▼
┌──────────────────────────────────────────────────────────────┐
│                  Application Layer                            │
│                (Use Cases / Orchestration)                    │
│   internal/service/*.go                                       │
│   - Orchestrates domain operations                            │
│   - Manages transactions                                      │
│   - Coordinates multiple repositories                         │
└──────────────────────┬───────────────────────────────────────┘
                       │ uses
                       ▼
┌──────────────────────────────────────────────────────────────┐
│                     Domain Layer                              │
│                  (Business Logic)                             │
│   internal/domain/entities/*.go                               │
│   internal/domain/services/*.go                               │
│   - Pure business rules                                       │
│   - NO external dependencies                                  │
│   - NO I/O operations                                         │
└──────────────────────┬───────────────────────────────────────┘
                       │ depends on
                       ▼
┌──────────────────────────────────────────────────────────────┐
│                    Ports (Interfaces)                         │
│   internal/ports/*.go                                         │
│   - Repository interfaces                                     │
│   - Service interfaces                                        │
│   - External system contracts                                 │
└──────────────────────┬───────────────────────────────────────┘
                       │ implemented by
                       ▼
┌──────────────────────────────────────────────────────────────┐
│                 Adapters (Outbound)                           │
│   internal/repository/postgres/*.go                           │
│   - Database implementations                                  │
│   - Cache implementations                                     │
│   - External service clients                                  │
└──────────────────────────────────────────────────────────────┘
```

## Directory Structure

```
/var/folders/70/y0td7vjs3hz520207tgz_y0m0000gn/T/bootstrap-email-worker-3343678467/email-worker/
├── cmd/
│   └── server/
│       └── main.go              # Dependency injection & server setup
│
├── internal/
│   ├── domain/                  # Business logic (pure Go)
│   │   ├── entities/            # Domain entities
│   │   │   ├── user.go
│   │   │   └── email.go        # Value objects
│   │   ├── services/            # Domain services
│   │   │   └── user_service.go
│   │   └── errors.go            # Domain errors
│   │
│   ├── ports/                   # Interfaces (contracts)
│   │   ├── repositories.go      # Repository interfaces
│   │   ├── cache.go            # Cache interfaces
│   │   ├── transaction.go      # Transaction manager
│   │   ├── external.go         # External service interfaces
│   │   └── mocks/              # Generated mocks for testing
│   │
│   ├── service/                 # Use cases (orchestration)
│   │   └── user_service.go     # Business use cases
│   │
│   ├── handlers/                # HTTP handlers (inbound adapters)
│   │   ├── user_handler.go
│   │   └── base.go
│   │
│   ├── repository/              # Data persistence (outbound adapters)
│   │   ├── postgres/
│   │   │   └── user.go
│   │   └── memory/             # In-memory for testing
│   │       └── user.go
│   │
│   ├── config/                  # Configuration
│   └── errors/                  # Error handling utilities
│
├── migrations/                  # Database migrations
│   └── postgres/
├── test/
│   └── integration/            # Integration tests
└── docs/
    ├── ARCHITECTURE.md         # This file
    ├── DEVELOPMENT.md          # Development guide
    └── TESTING.md             # Testing guide
```

## Layer Responsibilities

### Domain Layer (`internal/domain/`)

**What it does:**
- Defines business entities and their behavior
- Enforces business rules and invariants
- Contains domain services for multi-entity operations
- Defines domain errors

**What it CANNOT do:**
- Import database libraries (sql, gorm, etc.)
- Import HTTP libraries (chi, gin, etc.)
- Perform I/O operations
- Depend on infrastructure

**Example:**
```go
// internal/domain/entities/user.go
type User struct {
    id        uuid.UUID
    email     Email  // Value object
    name      string
    createdAt time.Time
    updatedAt time.Time
}

// Business logic method
func (u *User) ChangeEmail(newEmail string) error {
    email, err := NewEmail(newEmail)  // Validation in value object
    if err != nil {
        return err
    }
    u.email = email
    u.updatedAt = time.Now()
    return nil
}

// Validation business rule
func (u *User) Validate() error {
    if u.name == "" {
        return ErrValidation("name", "name cannot be empty")
    }
    if len(u.name) < 2 {
        return ErrValidation("name", "name must be at least 2 characters")
    }
    return nil
}
```

### Domain-Driven Design: Aggregate Roots

**What is an Aggregate?**

An aggregate is a cluster of domain objects (entities and value objects) that are treated as a single unit for data changes. The aggregate root is the entry point entity that enforces invariants across the entire aggregate.

**Why Aggregates Matter:**
- **Consistency boundaries**: Define what must be consistent in a single transaction
- **Encapsulation**: Hide internal complexity behind the root entity
- **Concurrency control**: Version/lock the root to detect conflicts
- **Transaction boundaries**: Each aggregate root is a unit of transactional consistency

**Decision Framework: When is something an Aggregate Root?**

Ask these questions:
1. **Does it have its own lifecycle?** Can it be created, modified, and deleted independently?
2. **Does it enforce invariants?** Does it need to validate rules across multiple objects?
3. **Is it a natural transactional boundary?** Should changes to it be atomic?
4. **Can external systems reference it?** Do other parts of the system need to look it up directly?

If you answer YES to most of these, it's likely an aggregate root.

**Examples:**

**Clear Aggregate Root: User**
```go
// User is an aggregate root
type User struct {
    id        uuid.UUID
    email     Email          // Value object
    profile   UserProfile    // Entity within aggregate
    settings  UserSettings   // Entity within aggregate
    createdAt time.Time
}

// ✅ Enforce invariants across the aggregate
func (u *User) UpdateProfile(name, bio string) error {
    if u.IsDeactivated() {
        return ErrInvalidOperation("cannot update deactivated user profile")
    }
    return u.profile.Update(name, bio)
}

// ✅ Control access to internal entities
func (u *User) ChangeEmailPreference(enabled bool) error {
    if !u.email.IsVerified() {
        return ErrInvalidOperation("verify email before changing preferences")
    }
    u.settings.SetEmailEnabled(enabled)
    return nil
}
```

**Clear Aggregate Root: Order**
```go
// Order is an aggregate root
type Order struct {
    id          uuid.UUID
    customerID  uuid.UUID    // Reference to Customer aggregate
    items       []OrderItem  // Entities owned by Order
    status      OrderStatus
    total       Money        // Value object
    placedAt    time.Time
}

// ✅ OrderItem only accessed through Order
type OrderItem struct {
    productID uuid.UUID     // Reference to Product aggregate
    quantity  int
    price     Money
}

// ✅ Enforce business rules across items
func (o *Order) AddItem(productID uuid.UUID, quantity int, price Money) error {
    if o.status != OrderStatusDraft {
        return ErrInvalidOperation("cannot modify submitted order")
    }
    if quantity <= 0 {
        return ErrValidation("quantity", "must be positive")
    }

    // Enforce aggregate invariant: max 10 items per order
    if len(o.items) >= 10 {
        return ErrInvalidOperation("order cannot have more than 10 items")
    }

    o.items = append(o.items, OrderItem{
        productID: productID,
        quantity:  quantity,
        price:     price,
    })
    o.recalculateTotal()
    return nil
}

// ✅ Maintain consistency when transitioning state
func (o *Order) Submit() error {
    if len(o.items) == 0 {
        return ErrValidation("items", "order must have at least one item")
    }
    if o.status != OrderStatusDraft {
        return ErrInvalidOperation("order already submitted")
    }
    o.status = OrderStatusPending
    o.placedAt = time.Now()
    return nil
}
```

**Debatable: Blog Post and Comments**

**Option A: Post is aggregate root, Comments are separate**
```go
// ✅ If comments have independent lifecycle and no invariants with Post
type Post struct {
    id        uuid.UUID
    authorID  uuid.UUID
    title     string
    content   string
    published bool
}

// Comments are separate aggregate root
type Comment struct {
    id       uuid.UUID
    postID   uuid.UUID  // Reference to Post aggregate
    authorID uuid.UUID
    content  string
}

// Use when:
// - Comments can be added/deleted without modifying Post
// - No invariant like "post must have >0 comments to publish"
// - Different services/teams may own Posts vs Comments
```

**Option B: Post is aggregate root, Comments are owned**
```go
// ✅ If comments must respect Post invariants
type Post struct {
    id        uuid.UUID
    authorID  uuid.UUID
    title     string
    content   string
    comments  []Comment  // Owned by aggregate
    published bool
}

type Comment struct {
    id       uuid.UUID
    authorID uuid.UUID
    content  string
}

// Enforce invariants through root
func (p *Post) AddComment(authorID uuid.UUID, content string) error {
    if !p.published {
        return ErrInvalidOperation("cannot comment on unpublished post")
    }
    if p.IsLocked() {
        return ErrInvalidOperation("post is locked for comments")
    }
    p.comments = append(p.comments, Comment{
        id:       uuid.New(),
        authorID: authorID,
        content:  content,
    })
    return nil
}

// Use when:
// - Post has business rules about comments (max count, approval, locking)
// - Comments only make sense in context of their Post
// - You need transactional consistency (e.g., publish post + first comment)
```

**The Four Critical Aggregate Rules:**

**1. Reference other aggregates by ID only**
```go
// ❌ BAD: Holding reference to another aggregate
type Order struct {
    customer *Customer  // Don't do this!
    items    []OrderItem
}

// ✅ GOOD: Reference by ID
type Order struct {
    customerID uuid.UUID  // Reference by identifier
    items      []OrderItem
}
```

**2. One repository per aggregate root**
```go
// ✅ Correct: Repository for aggregate root only
type OrderRepository interface {
    Save(ctx context.Context, order *Order) error
    FindByID(ctx context.Context, id uuid.UUID) (*Order, error)
    // Load entire aggregate (including owned entities)
}

// ❌ WRONG: No separate repository for owned entities
// OrderItemRepository should NOT exist if OrderItem is owned by Order
```

**3. Transactions should not span aggregate roots**
```go
// ❌ BAD: Modifying multiple aggregates in one transaction
func (s *Service) PlaceOrder(ctx context.Context, orderID, customerID uuid.UUID) error {
    return s.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
        // Don't modify Customer AND Order in same transaction
        customer := s.customerRepo.FindByID(txCtx, customerID)
        customer.IncrementOrderCount()  // ❌ Modifying multiple roots
        s.customerRepo.Save(txCtx, customer)

        order := s.orderRepo.FindByID(txCtx, orderID)
        order.Submit()
        s.orderRepo.Save(txCtx, order)
        return nil
    })
}

// ✅ GOOD: One aggregate per transaction, eventual consistency between aggregates
func (s *Service) PlaceOrder(ctx context.Context, orderID uuid.UUID) error {
    // Transaction only touches Order aggregate
    err := s.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
        order := s.orderRepo.FindByID(txCtx, orderID)
        if err := order.Submit(); err != nil {
            return err
        }
        return s.orderRepo.Save(txCtx, order)
    })

    if err != nil {
        return err
    }

    // Update Customer aggregate separately (eventual consistency)
    // Use domain events or separate transaction
    go s.eventBus.Publish(OrderPlacedEvent{OrderID: orderID})
    return nil
}
```

**4. Keep aggregates small**
```go
// ❌ BAD: Too large, too many owned entities
type Organization struct {
    id       uuid.UUID
    users    []User      // Could be thousands!
    projects []Project   // Could be hundreds!
    invoices []Invoice   // Could be years of data!
}

// ✅ GOOD: Small aggregate, reference others by ID
type Organization struct {
    id      uuid.UUID
    name    string
    ownerId uuid.UUID  // Reference to User aggregate
    settings OrganizationSettings
}

// Users, Projects, Invoices are separate aggregate roots
// They reference Organization by organizationID
```

**Repository Design Implications:**

```go
// Repository loads/saves entire aggregate atomically
type OrderRepository interface {
    // ✅ Save persists Order AND all OrderItems together
    Save(ctx context.Context, order *Order) error

    // ✅ Load retrieves Order AND all OrderItems together
    FindByID(ctx context.Context, id uuid.UUID) (*Order, error)

    // ✅ Query methods on aggregate root only
    FindByCustomerID(ctx context.Context, customerID uuid.UUID) ([]*Order, error)
    FindByStatus(ctx context.Context, status OrderStatus) ([]*Order, error)
}

// Implementation ensures atomicity
func (r *PostgresOrderRepository) Save(ctx context.Context, order *Order) error {
    // Use transaction to save order + items together
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return err
    }
    defer tx.Rollback()

    // Save root
    if err := r.saveOrder(ctx, tx, order); err != nil {
        return err
    }

    // Save owned entities
    if err := r.saveOrderItems(ctx, tx, order.ID(), order.Items()); err != nil {
        return err
    }

    return tx.Commit()
}
```

**Transaction Boundaries - Quick Reference:**

```
✅ Single Transaction:
   Order + OrderItems (same aggregate)
   User + UserProfile + UserSettings (same aggregate)

❌ Separate Transactions (use events for consistency):
   Order + Customer (different aggregates)
   Post + Author (different aggregates)
   Invoice + Payment (different aggregates)
```

**Quick Decision Checklist:**

When designing aggregates, ask:
- [ ] Does this entity have its own identity and lifecycle?
- [ ] Does it enforce invariants across owned entities?
- [ ] Can I keep it small (< 5 owned entities)?
- [ ] Does it represent a natural transaction boundary?
- [ ] Do other aggregates reference it by ID only?
- [ ] Does it have exactly one repository?

If you answer YES to all, you have a well-designed aggregate root.

### Ports Layer (`internal/ports/`)

**What it does:**
- Defines interfaces for external systems
- Specifies contracts between layers
- Enables dependency inversion
- Provides mock generation directives

**Key Principles:**
- Interfaces depend only on domain types
- No infrastructure types in signatures
- Context-first for all operations
- Returns domain errors, not infrastructure errors

**Example:**
```go
// internal/ports/repositories.go
type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
    FindByEmail(ctx context.Context, email string) (*entities.User, error)
    List(ctx context.Context, limit, offset int) ([]*entities.User, error)
    Update(ctx context.Context, user *entities.User) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### Application/Service Layer (`internal/service/`)

**What it does:**
- Implements use cases (business scenarios)
- Orchestrates domain logic
- Manages transactions
- Coordinates multiple repositories
- Handles caching strategies

**What it CANNOT do:**
- Contain business logic (that goes in domain)
- Depend on concrete implementations (use ports)
- Handle HTTP/gRPC details (that's handlers' job)

**Example:**
```go
// internal/service/user_service.go
type UserService struct {
    userRepo ports.UserRepository    // Interface!
    cache    ports.CacheService      // Interface!
    txMgr    ports.TransactionManager // Interface!
}

func (s *UserService) CreateUser(ctx context.Context, email, name string) (*entities.User, error) {
    // 1. Create domain entity (with validation)
    user, err := entities.NewUser(email, name)
    if err != nil {
        return nil, err
    }

    // 2. Check uniqueness via repository
    existing, err := s.userRepo.FindByEmail(ctx, email)
    if err != nil && !errors.Is(err, domain.ErrNotFound) {
        return nil, err
    }
    if existing != nil {
        return nil, domain.ErrConflict
    }

    // 3. Persist via repository port
    if err := s.userRepo.Create(ctx, user); err != nil {
        return nil, err
    }

    // 4. Invalidate cache (best effort)
    if err := s.cache.DeletePattern(ctx, "users:*"); err != nil {
        // Log but don't fail
        log.Printf("cache invalidation failed: %v", err)
    }

    return user, nil
}
```

### Handlers Layer (`internal/handlers/`)

**What it does:**
- Receives HTTP/gRPC requests
- Validates input
- Calls use cases
- Maps responses
- Handles HTTP-specific concerns (status codes, headers)

**Example:**
```go
// internal/handlers/user_handler.go
type UserHandler struct {
    userService *service.UserService
}

func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    var req CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.BadRequest(w, "invalid request body")
        return
    }

    // 2. Call use case
    user, err := h.userService.CreateUser(r.Context(), req.Email, req.Name)
    if err != nil {
        // Map domain errors to HTTP status codes
        switch {
        case errors.Is(err, domain.ErrConflict):
            response.Conflict(w, "email already exists")
        case errors.Is(err, domain.ErrValidation):
            response.BadRequest(w, err.Error())
        default:
            response.InternalError(w, "failed to create user")
        }
        return
    }

    // 3. Map to response DTO
    resp := UserResponse{
        ID:        user.ID().String(),
        Email:     user.Email().String(),
        Name:      user.Name(),
        CreatedAt: user.CreatedAt(),
    }

    response.Created(w, resp)
}
```

### Repository Adapters (`internal/repository/`)

**What it does:**
- Implements port interfaces
- Handles database connections
- Maps domain entities to/from database models
- Maps infrastructure errors to domain errors

**Example:**
```go
// internal/repository/postgres/user.go
type PostgresUserRepository struct {
    db *sql.DB
}

func (r *PostgresUserRepository) Create(ctx context.Context, user *entities.User) error {
    query := `
        INSERT INTO users (id, email, name, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5)
    `

    _, err := r.db.ExecContext(ctx, query,
        user.ID(),
        user.Email().String(),
        user.Name(),
        user.CreatedAt(),
        user.UpdatedAt(),
    )

    if err != nil {
        // Map infrastructure errors to domain errors
        if isUniqueViolation(err) {
            return domain.ErrConflict
        }
        return fmt.Errorf("failed to create user: %w", domain.ErrInternal)
    }

    return nil
}

func (r *PostgresUserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
    query := `
        SELECT id, email, name, created_at, updated_at
        FROM users
        WHERE id = $1
    `

    var (
        userID    uuid.UUID
        email     string
        name      string
        createdAt time.Time
        updatedAt time.Time
    )

    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &userID, &email, &name, &createdAt, &updatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrNotFound
        }
        return nil, fmt.Errorf("failed to find user: %w", domain.ErrInternal)
    }

    // Reconstitute domain entity from database
    return entities.ReconstitateUser(userID, email, name, createdAt, updatedAt)
}
```

## Dependency Flow

```
HTTP Request
    ↓
HTTP Handler (internal/handlers)
    ↓ calls
Service/Use Case (internal/service)
    ↓ uses
Domain Entity/Service (internal/domain)
    ↓ through
Repository Interface (internal/ports)
    ↓ implemented by
PostgreSQL Repository (internal/repository/postgres)
    ↓ queries
Database
```

## Key Principles

### 1. Dependency Inversion

**Bad** (concrete dependency):
```go
type UserService struct {
    db *sql.DB  // ❌ Depends on concrete implementation
}
```

**Good** (interface dependency):
```go
type UserService struct {
    userRepo ports.UserRepository  // ✅ Depends on interface
}
```

### 2. Single Responsibility

Each layer has ONE responsibility:
- **Domain**: Business rules
- **Service**: Orchestration
- **Ports**: Contracts
- **Handlers**: HTTP transport
- **Repository**: Data persistence

### 3. Open/Closed

Add features by:
- Creating new use cases (open for extension)
- NOT modifying existing use cases (closed for modification)

### 4. Interface Segregation

Small, focused interfaces:
```go
// ✅ Good: Focused interface
type UserRepository interface {
    Create(ctx context.Context, user *entities.User) error
    FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)
}

// ❌ Bad: God interface with 20+ methods
type Repository interface {
    CreateUser(...)
    CreatePost(...)
    CreateComment(...)
    // Too much!
}
```

### 5. Testability

Mock interfaces for unit testing:
```go
func TestCreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    defer ctrl.Finish()

    mockRepo := mocks.NewMockUserRepository(ctrl)
    mockRepo.EXPECT().
        Create(gomock.Any(), gomock.Any()).
        Return(nil)

    service := NewUserService(mockRepo, nil, nil)
    user, err := service.CreateUser(ctx, "test@example.com", "Test User")

    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

## Common Patterns

### Pattern 1: Value Objects

Encapsulate validation in immutable value objects:

```go
// internal/domain/entities/email.go
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if !emailRegex.MatchString(value) {
        return Email{}, ErrInvalidEmail("invalid email format")
    }
    return Email{value: value}, nil
}

func (e Email) String() string {
    return e.value
}
```

### Pattern 2: Repository with Context

All repository methods accept `context.Context` for:
- Cancellation
- Timeout
- Transaction propagation

```go
func (r *Repo) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
    // Context used for query timeout and transaction participation
    row := r.db.QueryRowContext(ctx, query, id)
    // ...
}
```

### Pattern 3: Error Mapping

Adapters map infrastructure errors to domain errors:

```go
func (r *Repo) Create(ctx context.Context, user *User) error {
    _, err := r.db.ExecContext(ctx, query, args...)
    if err != nil {
        if isUniqueViolation(err) {
            return domain.ErrConflict  // Map to domain error
        }
        if err == sql.ErrNoRows {
            return domain.ErrNotFound
        }
        return domain.ErrInternal
    }
    return nil
}
```

### Pattern 4: Transactional Use Cases

Use TransactionManager for atomic operations:

```go
func (s *Service) TransferFunds(ctx context.Context, fromID, toID uuid.UUID, amount int) error {
    return s.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
        // All operations use txCtx for transaction participation
        if err := s.accountRepo.Withdraw(txCtx, fromID, amount); err != nil {
            return err  // Automatic rollback
        }
        if err := s.accountRepo.Deposit(txCtx, toID, amount); err != nil {
            return err  // Automatic rollback
        }
        return nil  // Automatic commit
    })
}
```

## Testing Strategy

### Unit Tests (Domain Layer)
- Test business logic without mocks
- No database required
- Fast execution (<1ms per test)

```go
func TestUser_ChangeEmail(t *testing.T) {
    user, _ := entities.NewUser("old@example.com", "Test")

    err := user.ChangeEmail("new@example.com")

    assert.NoError(t, err)
    assert.Equal(t, "new@example.com", user.Email().String())
}
```

### Unit Tests (Service Layer)
- Mock all dependencies
- Test orchestration logic
- Verify interactions

```go
func TestUserService_CreateUser(t *testing.T) {
    ctrl := gomock.NewController(t)
    mockRepo := mocks.NewMockUserRepository(ctrl)
    mockCache := mocks.NewMockCacheService(ctrl)

    mockRepo.EXPECT().FindByEmail(gomock.Any(), "test@example.com").Return(nil, domain.ErrNotFound)
    mockRepo.EXPECT().Create(gomock.Any(), gomock.Any()).Return(nil)
    mockCache.EXPECT().DeletePattern(gomock.Any(), "users:*").Return(nil)

    service := NewUserService(mockRepo, mockCache, nil)
    user, err := service.CreateUser(ctx, "test@example.com", "Test")

    assert.NoError(t, err)
    assert.NotNil(t, user)
}
```

### Integration Tests (Repository)
- Test repository implementations
- Use test database (dockertest or testcontainers)
- Verify SQL queries work correctly

```go
func TestPostgresUserRepository_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := postgres.NewUserRepository(db)

    user, _ := entities.NewUser("test@example.com", "Test")
    err := repo.Create(context.Background(), user)

    assert.NoError(t, err)

    // Verify it was persisted
    found, err := repo.FindByID(context.Background(), user.ID())
    assert.NoError(t, err)
    assert.Equal(t, user.Email(), found.Email())
}
```

### End-to-End Tests
- Test complete user flows
- Use test HTTP server
- Verify full request/response cycle

## Further Reading

- [Domain Layer README](../internal/domain/README.md)
- [Ports Layer README](../internal/ports/README.md)
- [Development Guide](./DEVELOPMENT.md)
- [Testing Guide](./TESTING.md)

## Questions?

- **Q: Why not just use MVC?**
  - A: Hexagonal architecture enforces clean separation. In MVC, models often contain database logic. Here, domain is pure business logic.

- **Q: Isn't this over-engineered?**
  - A: For CRUD apps, maybe. For complex business logic, the separation pays off in maintainability and testability.

- **Q: Can I skip the domain layer?**
  - A: Not recommended. Even simple validation logic belongs in the domain, not in handlers or services.

- **Q: How do I handle authentication?**
  - A: Authentication is infrastructure concern (middleware). Authorization rules go in domain/use cases.

- **Q: Where do DTOs go?**
  - A: Request/response DTOs can live alongside handlers. They're transport-specific, not domain concepts.
