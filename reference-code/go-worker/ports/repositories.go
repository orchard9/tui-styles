package ports

import (
	"context"

	"github.com/google/uuid"

	"github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
)

//go:generate mockgen -source=repositories.go -destination=mocks/repositories.go -package=mocks

// ListOptions configures pagination and filtering for List queries.
//
// Pagination Strategy:
// - Offset-based pagination (limit/offset) is provided for simplicity
// - For high-scale production systems, consider cursor-based pagination
// - Cursor pagination avoids consistency issues when data changes between pages
//
// Design Considerations:
// - Limit and Offset have sensible defaults (10, 0)
// - Implementations should enforce maximum limits to prevent resource exhaustion
// - Consider adding SortBy/SortOrder for flexible ordering
// - Consider adding filter criteria for more complex queries
//
// Migration Path to Cursor Pagination:
// If you need cursor-based pagination for high-scale systems:
//  1. Add Cursor field (e.g., base64-encoded bookmark)
//  2. Return NextCursor in response
//  3. Deprecate Offset field gradually
//  4. Update documentation with cursor usage examples
//
// Example Usage:
//
//	opts := ListOptions{Limit: 20, Offset: 0}
//	users, err := repo.List(ctx, opts)
type ListOptions struct {
	// Limit is the maximum number of results to return.
	// Default: 10, Maximum: 100 (enforced by implementations)
	Limit int

	// Offset is the number of results to skip.
	// Default: 0
	// Note: Large offsets can be inefficient; consider cursor pagination for deep paging
	Offset int

	// Future extensions:
	// SortBy    string // Field to sort by
	// SortOrder string // "asc" or "desc"
	// Filters   map[string]interface{} // Dynamic filtering
	// Cursor    string // For cursor-based pagination
}

// UserRepository defines the contract for User persistence.
//
// This is a PORT (interface) in hexagonal architecture.
// Implementations are ADAPTERS (e.g., PostgreSQL, in-memory).
//
// Design Principles:
// - Interface depends only on domain types
// - No infrastructure dependencies
// - Returns domain errors, not infrastructure errors
// - Context-first for cancellation/timeout
//
// Aggregate Root Pattern:
// - User is the aggregate root for its consistency boundary
// - Repository operations maintain aggregate invariants atomically
// - Child entities are accessed through the aggregate root
// - Cross-aggregate references use IDs, not direct entity references
// - Transaction boundaries align with aggregate boundaries
//
// Thread Safety:
// - Implementations MUST be safe for concurrent use
// - Repository handles its own synchronization
//
// Example Implementations:
// - adapters/repository/postgres/user.go
// - adapters/repository/memory/user.go
type UserRepository interface {
	// Create persists a new User.
	//
	// Business Rules:
	// - User must be valid (validated in domain layer)
	// - ID must be unique
	// - Email must be unique
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - entity: Domain entity to persist
	//
	// Returns:
	// - error: domain.ErrConflict if duplicate
	// - error: domain.ErrInternal for persistence failures
	//
	// Example:
	//   err := repo.Create(ctx, user)
	//   if errors.Is(err, domain.ErrConflict) {
	//       // Handle duplicate
	//   }
	Create(ctx context.Context, entity *entities.User) error

	// FindByID retrieves a User by unique identifier.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - id: Unique identifier
	//
	// Returns:
	// - *entities.User: Found entity
	// - error: domain.ErrNotFound if not found
	// - error: domain.ErrInternal for query failures
	//
	// Example:
	//   user, err := repo.FindByID(ctx, userID)
	//   if errors.Is(err, domain.ErrNotFound) {
	//       // Handle not found
	//   }
	FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error)

	// FindByEmail retrieves a User by email address.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - email: Email address to search
	//
	// Returns:
	// - *entities.User: Found entity
	// - error: domain.ErrNotFound if not found
	// - error: domain.ErrInternal for query failures
	//
	// Example:
	//   user, err := repo.FindByEmail(ctx, "user@example.com")
	FindByEmail(ctx context.Context, email string) (*entities.User, error)

	// List retrieves multiple Users with pagination using ListOptions.
	//
	// Pagination Design:
	// - Uses offset-based pagination (limit/offset) for simplicity
	// - Implementations should enforce maximum limits (typically 100)
	// - For high-scale systems with deep pagination, consider cursor-based approaches
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - opts: ListOptions configuring pagination behavior
	//
	// Returns:
	// - []*entities.User: List of entities (may be empty, never nil)
	// - error: domain.ErrInternal for query failures
	// - error: domain.ErrValidation if opts are invalid (e.g., negative limit)
	//
	// Ordering:
	// - Results ordered by created_at DESC (newest first)
	// - Consistent ordering ensures predictable pagination
	//
	// Performance Considerations:
	// - Large offsets can be slow; warn clients about deep pagination
	// - Consider adding indexes on sort columns
	// - Consider caching for frequently accessed pages
	//
	// Example:
	//   opts := ListOptions{Limit: 20, Offset: 0}
	//   users, err := repo.List(ctx, opts) // First page
	//   if err != nil {
	//       return err
	//   }
	//
	//   // Next page
	//   opts.Offset = 20
	//   nextUsers, err := repo.List(ctx, opts)
	List(ctx context.Context, opts ListOptions) ([]*entities.User, error)

	// Update persists changes to an existing User.
	//
	// Business Rules:
	// - User must exist
	// - Email must remain unique
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - entity: Domain entity with updated values
	//
	// Returns:
	// - error: domain.ErrNotFound if entity doesn't exist
	// - error: domain.ErrConflict if email conflict
	// - error: domain.ErrInternal for persistence failures
	//
	// Example:
	//   err := repo.Update(ctx, user)
	Update(ctx context.Context, entity *entities.User) error

	// Delete removes a User from persistence.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - id: Unique identifier of entity to delete
	//
	// Returns:
	// - error: domain.ErrNotFound if entity doesn't exist
	// - error: domain.ErrInternal for persistence failures
	//
	// Note: Consider soft delete for audit trails
	//
	// Example:
	//   err := repo.Delete(ctx, userID)
	Delete(ctx context.Context, id uuid.UUID) error
}
