# Repository Adapters

## Purpose

Repository adapters implement the `ports.*Repository` interfaces defined in `internal/ports`. They handle all data persistence concerns while keeping the domain layer pure.

## Architecture

```
internal/ports/*.go           ← Interface definitions
        ↑
        │ implements
        │
internal/repository/
├── postgres/                 ← Production PostgreSQL implementations
│   ├── user.go
│   └── transaction.go
│
└── memory/                   ← Testing in-memory implementations
    └── user.go
```

## Implementations

### PostgreSQL (`postgres/`)

**Purpose**: Production implementation for data persistence

**Features**:
- Uses `go-core-http-toolkit/db` for connection management
- Handles transactions and connection pooling
- Maps SQL errors to domain errors
- Thread-safe (connection pool handles concurrency)

**Error Mapping**:
| SQL Error | Domain Error | Description |
|-----------|-------------|-------------|
| `sql.ErrNoRows` | `domain.ErrNotFound` | Entity doesn't exist |
| Unique violation (23505) | `domain.ErrConflict` | Duplicate key |
| Foreign key violation (23503) | `domain.ErrConflict` | Referential integrity |
| Connection timeout | `domain.ErrInternal` | Database unavailable |
| Other SQL errors | `domain.ErrInternal` | Unexpected failures |

**Example**:
```go
// internal/repository/postgres/user.go
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
    query := "INSERT INTO users (id, email, name) VALUES ($1, $2, $3)"
    _, err := r.db.ExecContext(ctx, query, user.ID(), user.Email(), user.Name())
    if err != nil {
        if isUniqueViolation(err) {
            return domain.ErrConflict
        }
        return domain.ErrInternal
    }
    return nil
}
```

### In-Memory (`memory/`)

**Purpose**: Testing implementation without external dependencies

**Features**:
- Fast (no I/O operations)
- Thread-safe (uses sync.RWMutex)
- Deterministic behavior
- `Clear()` method for test cleanup
- Implements same sorting/pagination as PostgreSQL

**Not for Production**: Data is not persisted to disk

**Example**:
```go
// internal/repository/memory/user.go
func (r *UserRepository) Create(ctx context.Context, user *entities.User) error {
    r.mu.Lock()
    defer r.mu.Unlock()

    if _, exists := r.email[user.Email()]; exists {
        return domain.ErrConflict
    }

    r.data[user.ID()] = user
    r.email[user.Email()] = user.ID()
    return nil
}
```

## Usage

### In Production (cmd/server/main.go)

```go
import (
    "github.com/orchard9/go-core-http-toolkit/db"
    "github.com/orchard9/peach/apps/email-worker/internal/repository/postgres"
    "github.com/orchard9/peach/apps/email-worker/internal/service"
)

func main() {
    // Initialize database
    database, err := db.NewPostgresDB(db.Config{
        Host:     cfg.DBHost,
        Port:     cfg.DBPort,
        User:     cfg.DBUser,
        Password: cfg.DBPassword,
        DBName:   cfg.DBName,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer database.Close()

    // Create repository
    userRepo := postgres.NewUserRepository(database)

    // Inject into service
    userService := service.NewUserService(userRepo, cacheService, txMgr)

    // Use in handlers
    userHandler := handlers.NewUserHandler(userService)
}
```

### In Tests

```go
import (
    "github.com/orchard9/peach/apps/email-worker/internal/repository/memory"
    "github.com/orchard9/peach/apps/email-worker/internal/service"
)

func TestUserService_CreateUser(t *testing.T) {
    // Use in-memory repository (no database needed)
    repo := memory.NewInMemoryUserRepository()

    // Create service with mock dependencies
    userService := service.NewUserService(repo, mockCache, mockTxMgr)

    // Test without database
    user, err := userService.CreateUser(ctx, "test@example.com", "Test User")
    assert.NoError(t, err)
    assert.NotNil(t, user)

    // Verify it was stored
    found, err := repo.FindByID(ctx, user.ID())
    assert.NoError(t, err)
    assert.Equal(t, user.Email(), found.Email())

    // Cleanup for next test
    repo.Clear()
}
```

### Integration Tests (PostgreSQL)

```go
import (
    "github.com/ory/dockertest/v3"
    "github.com/orchard9/peach/apps/email-worker/internal/repository/postgres"
)

func setupTestDB(t *testing.T) (*db.DB, func()) {
    pool, _ := dockertest.NewPool("")

    resource, _ := pool.RunWithOptions(&dockertest.RunOptions{
        Repository: "postgres",
        Tag:        "15-alpine",
        Env: []string{
            "POSTGRES_PASSWORD=test",
            "POSTGRES_DB=testdb",
        },
    })

    // Connect and run migrations...

    cleanup := func() {
        pool.Purge(resource)
    }

    return database, cleanup
}

func TestPostgresUserRepository_Create(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := postgres.NewUserRepository(db)
    user, _ := entities.NewUser("test@example.com", "Test")

    err := repo.Create(context.Background(), user)
    assert.NoError(t, err)

    // Verify persistence
    found, err := repo.FindByID(context.Background(), user.ID())
    assert.NoError(t, err)
    assert.Equal(t, user.Email(), found.Email())
}
```

## Adding New Repositories

Follow these steps to add a new repository:

### 1. Define Port Interface

**File**: `internal/ports/repositories.go`

```go
type PostRepository interface {
    Create(ctx context.Context, post *entities.Post) error
    FindByID(ctx context.Context, id uuid.UUID) (*entities.Post, error)
    // ... other methods
}
```

### 2. Implement PostgreSQL Adapter

**File**: `internal/repository/postgres/post.go`

```go
type PostRepository struct {
    db *db.DB
}

func NewPostRepository(database *db.DB) ports.PostRepository {
    return &PostRepository{db: database}
}

func (r *PostRepository) Create(ctx context.Context, post *entities.Post) error {
    // SQL implementation with error mapping
}

// Compile-time verification
var _ ports.PostRepository = (*PostRepository)(nil)
```

### 3. Implement In-Memory Adapter

**File**: `internal/repository/memory/post.go`

```go
type PostRepository struct {
    mu   sync.RWMutex
    data map[uuid.UUID]*entities.Post
}

func NewInMemoryPostRepository() ports.PostRepository {
    return &PostRepository{
        data: make(map[uuid.UUID]*entities.Post),
    }
}

func (r *PostRepository) Create(ctx context.Context, post *entities.Post) error {
    // In-memory implementation
}

var _ ports.PostRepository = (*PostRepository)(nil)
```

### 4. Write Tests

**File**: `internal/repository/postgres/post_test.go`

```go
func TestPostgresPostRepository_Create(t *testing.T) {
    db, cleanup := setupTestDB(t)
    defer cleanup()

    repo := postgres.NewPostRepository(db)
    // Test implementation
}
```

### 5. Wire Up in main.go

```go
postRepo := postgres.NewPostRepository(database)
postService := service.NewPostService(postRepo)
```

## Error Handling

Repositories **must** map infrastructure errors to domain errors:

### PostgreSQL Error Mapping

```go
func (r *Repo) Create(ctx context.Context, entity *entities.Entity) error {
    _, err := r.db.ExecContext(ctx, query, args...)
    if err != nil {
        // Map to domain errors
        if isUniqueViolation(err) {
            return domain.ErrConflict
        }
        if isForeignKeyViolation(err) {
            return domain.ErrConflict
        }
        return domain.ErrInternal
    }
    return nil
}

func (r *Repo) FindByID(ctx context.Context, id uuid.UUID) (*entities.Entity, error) {
    err := r.db.GetContext(ctx, &row, query, id)
    if err != nil {
        if errors.Is(err, sql.ErrNoRows) {
            return nil, domain.ErrNotFound
        }
        return nil, domain.ErrInternal
    }
    return entity, nil
}
```

### Why Error Mapping?

1. **Domain independence**: Domain and application layers don't import SQL libraries
2. **Consistent errors**: All repositories return same error types
3. **Swappable implementations**: Can replace PostgreSQL with MongoDB without changing callers
4. **Testability**: Easy to test error handling with mocks

## Testing

### Unit Tests (In-Memory)

```bash
# Fast, no dependencies
go test ./internal/repository/memory/... -v
```

### Integration Tests (PostgreSQL)

```bash
# Requires Docker for test database
docker-compose up -d postgres
go test ./internal/repository/postgres/... -v
```

### Watch Mode

```bash
# Auto-rerun tests on file changes
make test-watch
```

## Best Practices

### 1. Always Map Errors

❌ **Bad**: Returning SQL errors directly
```go
return err  // Returns sql.ErrNoRows to caller
```

✅ **Good**: Map to domain errors
```go
if errors.Is(err, sql.ErrNoRows) {
    return domain.ErrNotFound
}
return domain.ErrInternal
```

### 2. Use Context for Cancellation

All repository methods accept `context.Context`:

```go
func (r *Repo) FindByID(ctx context.Context, id uuid.UUID) (*Entity, error) {
    // Context is passed to database operations
    err := r.db.GetContext(ctx, &row, query, id)
    // If context is cancelled, query is aborted
}
```

### 3. Verify Interface Implementation

Add compile-time verification:

```go
var _ ports.UserRepository = (*PostgresUserRepository)(nil)
```

This ensures your implementation satisfies the interface.

### 4. Idempotent Deletes

Deletes should not error if entity doesn't exist:

```go
func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
    _, err := r.db.ExecContext(ctx, "DELETE FROM table WHERE id = $1", id)
    // Don't check rows affected - idempotent
    return err
}
```

### 5. Empty Slices for List

Return empty slice, not error, when no results:

```go
func (r *Repo) List(ctx context.Context, limit, offset int) ([]*Entity, error) {
    var rows []Row
    err := r.db.SelectContext(ctx, &rows, query, limit, offset)
    if err != nil {
        return nil, err
    }
    // If no rows, returns empty slice (not nil)
    return rows, nil
}
```

## Common Patterns

### Pagination

```go
func (r *Repo) List(ctx context.Context, limit, offset int) ([]*Entity, error) {
    query := `
        SELECT * FROM table
        ORDER BY created_at DESC
        LIMIT $1 OFFSET $2
    `
    // ...
}
```

### Soft Delete

```go
func (r *Repo) Delete(ctx context.Context, id uuid.UUID) error {
    query := "UPDATE table SET deleted_at = NOW() WHERE id = $1"
    _, err := r.db.ExecContext(ctx, query, id)
    return err
}
```

### Optimistic Locking

```go
func (r *Repo) Update(ctx context.Context, entity *Entity) error {
    query := `
        UPDATE table SET data = $1, version = version + 1
        WHERE id = $2 AND version = $3
    `
    result, err := r.db.ExecContext(ctx, query, entity.Data(), entity.ID(), entity.Version())
    if err != nil {
        return err
    }
    if rowsAffected, _ := result.RowsAffected(); rowsAffected == 0 {
        return domain.ErrConflict  // Concurrent modification
    }
    return nil
}
```

## See Also

- [Ports Documentation](../ports/README.md) - Interface definitions
- [Domain Documentation](../domain/README.md) - Entity definitions
- [ARCHITECTURE.md](../../docs/ARCHITECTURE.md) - Overall architecture
- [TESTING.md](../../docs/TESTING.md) - Testing strategies
