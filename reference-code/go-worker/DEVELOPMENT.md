# Development Guide

This guide provides step-by-step instructions for common development tasks in this hexagonal architecture project.

## Table of Contents

- [Adding a New Feature](#adding-a-new-feature)
- [Adding a New Entity](#adding-a-new-entity)
- [Adding a New Use Case](#adding-a-new-use-case)
- [Adding a New Repository](#adding-a-new-repository)
- [Adding a New HTTP Endpoint](#adding-a-new-http-endpoint)
- [Working with Transactions](#working-with-transactions)
- [Common Pitfalls](#common-pitfalls)
- [Best Practices](#best-practices)

## Adding a New Feature

Follow these steps to add a complete feature (e.g., "Blog Post Management"):

### 1. Define Domain Entity

**File**: `internal/domain/entities/post.go`

```go
package entities

import (
    "time"
    "github.com/google/uuid"
    "github.com/orchard9/peach/apps/email-worker/internal/domain"
)

type Post struct {
    id        uuid.UUID
    title     string
    content   string
    authorID  uuid.UUID
    published bool
    createdAt time.Time
    updatedAt time.Time
}

// NewPost creates a new post with validation
func NewPost(title, content string, authorID uuid.UUID) (*Post, error) {
    if title == "" {
        return nil, domain.ErrValidation("title", "title cannot be empty")
    }
    if len(title) > 200 {
        return nil, domain.ErrValidation("title", "title too long")
    }
    if content == "" {
        return nil, domain.ErrValidation("content", "content cannot be empty")
    }

    now := time.Now()
    return &Post{
        id:        uuid.New(),
        title:     title,
        content:   content,
        authorID:  authorID,
        published: false,
        createdAt: now,
        updatedAt: now,
    }, nil
}

// Getters
func (p *Post) ID() uuid.UUID        { return p.id }
func (p *Post) Title() string        { return p.title }
func (p *Post) Content() string      { return p.content }
func (p *Post) AuthorID() uuid.UUID  { return p.authorID }
func (p *Post) Published() bool      { return p.published }
func (p *Post) CreatedAt() time.Time { return p.createdAt }
func (p *Post) UpdatedAt() time.Time { return p.updatedAt }

// Business methods
func (p *Post) Publish() error {
    if p.published {
        return domain.ErrConflict
    }
    p.published = true
    p.updatedAt = time.Now()
    return nil
}

func (p *Post) UpdateContent(title, content string) error {
    if title == "" || content == "" {
        return domain.ErrValidation("content", "title and content required")
    }
    p.title = title
    p.content = content
    p.updatedAt = time.Now()
    return nil
}

// ReconstitatePost recreates post from database
func ReconstitatePost(id uuid.UUID, title, content string, authorID uuid.UUID, published bool, createdAt, updatedAt time.Time) (*Post, error) {
    return &Post{
        id:        id,
        title:     title,
        content:   content,
        authorID:  authorID,
        published: published,
        createdAt: createdAt,
        updatedAt: updatedAt,
    }, nil
}
```

### 2. Define Repository Port

**File**: `internal/ports/repositories.go` (add to existing file)

```go
//go:generate mockgen -source=repositories.go -destination=mocks/repositories.go -package=mocks

type PostRepository interface {
    Create(ctx context.Context, post *entities.Post) error
    FindByID(ctx context.Context, id uuid.UUID) (*entities.Post, error)
    FindByAuthor(ctx context.Context, authorID uuid.UUID) ([]*entities.Post, error)
    List(ctx context.Context, limit, offset int) ([]*entities.Post, error)
    Update(ctx context.Context, post *entities.Post) error
    Delete(ctx context.Context, id uuid.UUID) error
}
```

### 3. Create Service/Use Case

**File**: `internal/service/post_service.go`

```go
package service

import (
    "context"
    "errors"
    "github.com/google/uuid"

    "github.com/orchard9/peach/apps/email-worker/internal/domain"
    "github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
    "github.com/orchard9/peach/apps/email-worker/internal/ports"
)

type PostService struct {
    postRepo ports.PostRepository
    userRepo ports.UserRepository
    cache    ports.CacheService
}

func NewPostService(postRepo ports.PostRepository, userRepo ports.UserRepository, cache ports.CacheService) *PostService {
    return &PostService{
        postRepo: postRepo,
        userRepo: userRepo,
        cache:    cache,
    }
}

func (s *PostService) CreatePost(ctx context.Context, title, content string, authorID uuid.UUID) (*entities.Post, error) {
    // Verify author exists
    author, err := s.userRepo.FindByID(ctx, authorID)
    if err != nil {
        if errors.Is(err, domain.ErrNotFound) {
            return nil, domain.ErrValidation("author_id", "author not found")
        }
        return nil, err
    }

    // Create domain entity
    post, err := entities.NewPost(title, content, author.ID())
    if err != nil {
        return nil, err
    }

    // Persist
    if err := s.postRepo.Create(ctx, post); err != nil {
        return nil, err
    }

    // Invalidate cache
    s.cache.DeletePattern(ctx, "posts:*")

    return post, nil
}

func (s *PostService) PublishPost(ctx context.Context, postID uuid.UUID) error {
    // Fetch post
    post, err := s.postRepo.FindByID(ctx, postID)
    if err != nil {
        return err
    }

    // Business logic
    if err := post.Publish(); err != nil {
        return err
    }

    // Persist changes
    if err := s.postRepo.Update(ctx, post); err != nil {
        return err
    }

    // Invalidate cache
    s.cache.DeletePattern(ctx, "posts:*")

    return nil
}
```

### 4. Implement Repository

**File**: `internal/repository/postgres/post.go`

```go
package postgres

import (
    "context"
    "database/sql"
    "fmt"

    "github.com/google/uuid"
    "github.com/lib/pq"

    "github.com/orchard9/peach/apps/email-worker/internal/domain"
    "github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
    "github.com/orchard9/peach/apps/email-worker/internal/ports"
)

type PostgresPostRepository struct {
    db *sql.DB
}

func NewPostRepository(db *sql.DB) ports.PostRepository {
    return &PostgresPostRepository{db: db}
}

func (r *PostgresPostRepository) Create(ctx context.Context, post *entities.Post) error {
    query := `
        INSERT INTO posts (id, title, content, author_id, published, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `

    _, err := r.db.ExecContext(ctx, query,
        post.ID(),
        post.Title(),
        post.Content(),
        post.AuthorID(),
        post.Published(),
        post.CreatedAt(),
        post.UpdatedAt(),
    )

    if err != nil {
        if pqErr, ok := err.(*pq.Error); ok {
            if pqErr.Code == "23505" { // unique_violation
                return domain.ErrConflict
            }
        }
        return fmt.Errorf("failed to create post: %w", domain.ErrInternal)
    }

    return nil
}

func (r *PostgresPostRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.Post, error) {
    query := `
        SELECT id, title, content, author_id, published, created_at, updated_at
        FROM posts
        WHERE id = $1
    `

    var (
        postID    uuid.UUID
        title     string
        content   string
        authorID  uuid.UUID
        published bool
        createdAt time.Time
        updatedAt time.Time
    )

    err := r.db.QueryRowContext(ctx, query, id).Scan(
        &postID, &title, &content, &authorID, &published, &createdAt, &updatedAt,
    )

    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrNotFound
        }
        return nil, fmt.Errorf("failed to find post: %w", domain.ErrInternal)
    }

    return entities.ReconstitatePost(postID, title, content, authorID, published, createdAt, updatedAt)
}
```

### 5. Create Migration

**File**: `migrations/postgres/003_create_posts_table.up.sql`

```sql
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY,
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    author_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    published BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_posts_author_id ON posts(author_id);
CREATE INDEX idx_posts_published ON posts(published);
CREATE INDEX idx_posts_created_at ON posts(created_at DESC);
```

**File**: `migrations/postgres/003_create_posts_table.down.sql`

```sql
DROP TABLE IF EXISTS posts;
```

### 6. Add HTTP Handler

**File**: `internal/handlers/post_handler.go`

```go
package handlers

import (
    "encoding/json"
    "net/http"

    "github.com/go-chi/chi/v5"
    "github.com/google/uuid"
    "github.com/orchard9/go-core-http-toolkit/response"

    "github.com/orchard9/peach/apps/email-worker/internal/service"
)

type PostHandler struct {
    postService *service.PostService
}

func NewPostHandler(postService *service.PostService) *PostHandler {
    return &PostHandler{postService: postService}
}

type CreatePostRequest struct {
    Title    string `json:"title"`
    Content  string `json:"content"`
    AuthorID string `json:"author_id"`
}

type PostResponse struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    Content   string `json:"content"`
    AuthorID  string `json:"author_id"`
    Published bool   `json:"published"`
    CreatedAt string `json:"created_at"`
}

func (h *PostHandler) CreatePost(w http.ResponseWriter, r *http.Request) {
    var req CreatePostRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        response.BadRequest(w, "invalid request body")
        return
    }

    authorID, err := uuid.Parse(req.AuthorID)
    if err != nil {
        response.BadRequest(w, "invalid author_id")
        return
    }

    post, err := h.postService.CreatePost(r.Context(), req.Title, req.Content, authorID)
    if err != nil {
        handleError(w, err)
        return
    }

    resp := PostResponse{
        ID:        post.ID().String(),
        Title:     post.Title(),
        Content:   post.Content(),
        AuthorID:  post.AuthorID().String(),
        Published: post.Published(),
        CreatedAt: post.CreatedAt().Format(time.RFC3339),
    }

    response.Created(w, resp)
}

func (h *PostHandler) PublishPost(w http.ResponseWriter, r *http.Request) {
    idStr := chi.URLParam(r, "id")
    id, err := uuid.Parse(idStr)
    if err != nil {
        response.BadRequest(w, "invalid id")
        return
    }

    if err := h.postService.PublishPost(r.Context(), id); err != nil {
        handleError(w, err)
        return
    }

    response.OK(w, map[string]string{"message": "post published"})
}
```

### 7. Register Routes

**File**: `cmd/server/main.go` (update)

```go
// Wire up post handlers
postRepo := postgres.NewPostRepository(db)
postService := service.NewPostService(postRepo, userRepo, cacheService)
postHandler := handlers.NewPostHandler(postService)

// Register routes
r.Route("/posts", func(r chi.Router) {
    r.Post("/", postHandler.CreatePost)
    r.Get("/{id}", postHandler.GetPost)
    r.Put("/{id}/publish", postHandler.PublishPost)
})
```

### 8. Generate Mocks

```bash
go generate ./internal/ports/...
```

### 9. Write Tests

See [TESTING.md](./TESTING.md) for testing strategies.

## Adding a New Entity

When adding a new domain entity:

1. **Create entity file**: `internal/domain/entities/entity_name.go`
2. **Follow patterns**:
   - Private fields
   - Public getters
   - Constructor with validation (`NewEntityName`)
   - Reconstitute function for database loading
   - Business methods (not CRUD!)

3. **Add to ports**: Define repository interface in `internal/ports/repositories.go`
4. **Generate mocks**: Run `go generate ./internal/ports/...`

## Adding a New Use Case

Use cases go in `internal/service/`:

1. **Create service struct** with port dependencies:
```go
type MyService struct {
    repo1 ports.Repository1
    repo2 ports.Repository2
    cache ports.CacheService
}
```

2. **Implement use case method**:
   - Accept context as first parameter
   - Return domain entities or errors
   - Orchestrate domain logic
   - Manage transactions if needed

3. **Write unit tests** with mocked ports

## Adding a New Repository

1. **Define port** in `internal/ports/repositories.go`
2. **Implement adapter** in `internal/repository/postgres/entity_name.go`
3. **Follow patterns**:
   - Accept context first
   - Return domain entities
   - Map errors to domain errors
   - Use prepared statements or query builders

4. **Write integration tests** with test database

## Adding a New HTTP Endpoint

1. **Add handler method** in `internal/handlers/`
2. **Define request/response structs** (DTOs)
3. **Call use case** from handler
4. **Map errors** to HTTP status codes
5. **Register route** in `cmd/server/main.go`

## Working with Transactions

Use `ports.TransactionManager` for atomic operations:

```go
func (s *Service) ComplexOperation(ctx context.Context) error {
    return s.txMgr.WithTransaction(ctx, func(txCtx context.Context) error {
        // All repository calls use txCtx
        if err := s.repo1.Create(txCtx, entity1); err != nil {
            return err  // Automatic rollback
        }
        if err := s.repo2.Update(txCtx, entity2); err != nil {
            return err  // Automatic rollback
        }
        return nil  // Automatic commit
    })
}
```

**Key points:**
- Pass transaction context to all repository operations
- Return error to rollback
- Return nil to commit
- Transaction manager handles begin/commit/rollback

## Common Pitfalls

### ❌ Don't put business logic in handlers

```go
// BAD
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // Validation and business logic in handler ❌
    if len(req.Name) < 2 {
        response.BadRequest(w, "name too short")
        return
    }
}
```

```go
// GOOD
func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) {
    // Delegate to service, which uses domain ✅
    user, err := h.service.CreateUser(ctx, req.Email, req.Name)
    if err != nil {
        handleError(w, err)
        return
    }
    response.Created(w, toResponse(user))
}
```

### ❌ Don't depend on concrete types

```go
// BAD
type UserService struct {
    db *sql.DB  // ❌ Concrete dependency
}
```

```go
// GOOD
type UserService struct {
    userRepo ports.UserRepository  // ✅ Interface dependency
}
```

### ❌ Don't import infrastructure in domain

```go
// BAD - internal/domain/entities/user.go
import "database/sql"  // ❌ Infrastructure import in domain

func (u *User) Save(db *sql.DB) error {  // ❌ Domain doing I/O
    // ...
}
```

```go
// GOOD - Domain is pure
func (u *User) Validate() error {  // ✅ Pure business logic
    // ...
}
```

### ❌ Don't return infrastructure errors

```go
// BAD
func (r *Repo) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
    err := r.db.QueryRow(...)
    return nil, err  // ❌ Returns sql.ErrNoRows
}
```

```go
// GOOD
func (r *Repo) FindByID(ctx context.Context, id uuid.UUID) (*User, error) {
    err := r.db.QueryRow(...)
    if err == sql.ErrNoRows {
        return nil, domain.ErrNotFound  // ✅ Maps to domain error
    }
    return nil, err
}
```

## Best Practices

### 1. Keep Domain Pure

- No imports outside stdlib (except uuid)
- No I/O operations
- No infrastructure dependencies
- Pure functions when possible

### 2. Use Value Objects

Encapsulate validation:

```go
type Email struct {
    value string
}

func NewEmail(value string) (Email, error) {
    if !isValidEmail(value) {
        return Email{}, ErrInvalidEmail
    }
    return Email{value: value}, nil
}
```

### 3. Validate in Domain

All business rules belong in domain:

```go
func NewUser(email, name string) (*User, error) {
    if len(name) < 2 {
        return nil, ErrValidation("name", "name too short")
    }
    // ...
}
```

### 4. Small, Focused Interfaces

```go
// Good
type UserFinder interface {
    FindByID(ctx context.Context, id uuid.UUID) (*User, error)
}

// Bad - too broad
type UserRepository interface {
    FindByID(...)
    FindByEmail(...)
    FindByAge(...)
    // 20 more methods...
}
```

### 5. Context First

Always pass context as first parameter:

```go
func (s *Service) Method(ctx context.Context, arg string) error
```

### 6. Meaningful Errors

Use domain errors, not strings:

```go
// Bad
return errors.New("user not found")

// Good
return domain.ErrNotFound
```

### 7. Test at Right Level

- **Domain**: Pure unit tests, no mocks
- **Service**: Unit tests with mocked ports
- **Repository**: Integration tests with real database
- **Handler**: End-to-end tests with test server

### 8. Use Constructors

Always provide constructors that enforce invariants:

```go
func NewUser(email, name string) (*User, error) {
    // Validation
    // Initialization
    return user, nil
}
```

## Quick Reference

### File Locations

- **Domain entities**: `internal/domain/entities/`
- **Domain services**: `internal/domain/services/`
- **Domain errors**: `internal/domain/errors.go`
- **Port interfaces**: `internal/ports/`
- **Use cases**: `internal/service/`
- **HTTP handlers**: `internal/handlers/`
- **Repository implementations**: `internal/repository/postgres/`
- **Migrations**: `migrations/postgres/`
- **Tests**: `test/integration/`, `*_test.go` next to code

### Commands

```bash
# Generate mocks
go generate ./internal/ports/...

# Run tests
make test

# Run tests with coverage
make test-coverage

# Run linter
make lint

# Run migrations
make migrate-up

# Create new migration
migrate create -ext sql -dir migrations/postgres -seq migration_name
```

## Need Help?

- Read [ARCHITECTURE.md](./ARCHITECTURE.md) for layer concepts
- Read [TESTING.md](./TESTING.md) for testing strategies
- Check existing code for patterns
- Review domain/ports/adapters for examples
