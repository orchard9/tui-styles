// Package postgres implements repository interfaces for PostgreSQL.
// This is an ADAPTER that implements the port interfaces defined in internal/ports.
//
// Responsibilities:
// - Implement repository interfaces
// - Handle SQL queries and transactions
// - Map database errors to domain errors
// - Manage database connection lifecycle
package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/orchard9/go-core-http-toolkit/db"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
	"github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
	"github.com/orchard9/peach/apps/email-worker/internal/ports"
)

// UserRepository implements ports.UserRepository for PostgreSQL
//
// Thread-Safety: Safe for concurrent use (db.DB handles connection pooling)
//
// Error Handling:
// - Maps SQL errors to domain errors
// - Returns domain.ErrNotFound for sql.ErrNoRows
// - Returns domain.ErrConflict for unique violations
// - Returns domain.ErrInternal for other database errors
type UserRepository struct {
	db *db.DB
}

// NewUserRepository creates a new PostgreSQL repository
//
// Parameters:
// - db: Database connection pool (managed by go-core-http-toolkit)
//
// Returns configured repository ready for use
func NewUserRepository(database *db.DB) ports.UserRepository {
	return &UserRepository{
		db: database,
	}
}

// Create persists a new User to the database
//
// SQL: INSERT INTO users (id, email, name, created_at, updated_at) VALUES (...)
//
// Error Mapping:
// - Unique constraint violation → domain.ErrConflict
// - Other errors → domain.ErrInternal
func (r *UserRepository) Create(ctx context.Context, entity *entities.User) error {
	query := `
		INSERT INTO users (
			id,
			email,
			name,
			created_at,
			updated_at
		) VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.ExecContext(ctx, query,
		entity.ID(),
		entity.Email(),
		entity.Name(),
		entity.CreatedAt(),
		entity.UpdatedAt(),
	)

	if err != nil {
		// Map database errors to domain errors
		if isUniqueViolation(err) {
			return fmt.Errorf("%w: email already exists", domain.ErrConflict)
		}
		return fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return nil
}

// FindByID retrieves a User by unique identifier
//
// SQL: SELECT * FROM users WHERE id = $1
//
// Returns:
// - domain.ErrNotFound if entity doesn't exist
// - domain.ErrInternal for database errors
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	query := `
		SELECT
			id,
			email,
			name,
			created_at,
			updated_at
		FROM users
		WHERE id = $1
	`

	var row struct {
		ID        uuid.UUID
		Email     string
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	err := r.db.GetContext(ctx, &row, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: user with id %s", domain.ErrNotFound, id)
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	// Reconstruct domain entity from database row
	return entities.ReconstituteUser(
		row.ID,
		row.Email,
		row.Name,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

// Update persists changes to an existing User
//
// SQL: UPDATE users SET ... WHERE id = $1
//
// Returns:
// - domain.ErrNotFound if entity doesn't exist
// - domain.ErrConflict for unique violations
// - domain.ErrInternal for database errors
func (r *UserRepository) Update(ctx context.Context, entity *entities.User) error {
	query := `
		UPDATE users
		SET
			email = $2,
			name = $3,
			updated_at = $4
		WHERE id = $1
	`

	result, err := r.db.ExecContext(ctx, query,
		entity.ID(),
		entity.Email(),
		entity.Name(),
		entity.UpdatedAt(),
	)

	if err != nil {
		if isUniqueViolation(err) {
			return fmt.Errorf("%w: email already exists", domain.ErrConflict)
		}
		return fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("%w: user with id %s", domain.ErrNotFound, entity.ID())
	}

	return nil
}

// Delete removes a User from the database
//
// SQL: DELETE FROM users WHERE id = $1
//
// Idempotent: No error if entity already deleted
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := "DELETE FROM users WHERE id = $1"

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	// Idempotent - no error if already deleted
	return nil
}

// List retrieves multiple Users with offset-based pagination
//
// Example usage:
//
//	opts := ports.ListOptions{Limit: 20, Offset: 0}
//	entities, err := repo.List(ctx, opts)
//
// Returns:
// - Slice of entities matching the query
// - Empty slice if no results (not an error)
// - domain.ErrInternal for database errors
func (r *UserRepository) List(ctx context.Context, opts ports.ListOptions) ([]*entities.User, error) {
	// Set defaults
	limit := opts.Limit
	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}

	offset := opts.Offset
	if offset < 0 {
		offset = 0
	}

	// Build simple query with limit and offset
	query := `
		SELECT
			id,
			email,
			name,
			created_at,
			updated_at
		FROM users
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	args := []interface{}{limit, offset}

	var rows []struct {
		ID        uuid.UUID
		Email     string
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	err := r.db.SelectContext(ctx, &rows, query, args...)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	// Convert rows to domain entities
	items := make([]*entities.User, 0, len(rows))
	for _, row := range rows {
		entity, err := entities.ReconstituteUser(
			row.ID,
			row.Email,
			row.Name,
			row.CreatedAt,
			row.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
		}
		items = append(items, entity)
	}

	return items, nil
}

// FindByEmail retrieves a User by email address
//
// SQL: SELECT * FROM users WHERE email = $1
//
// Returns:
// - domain.ErrNotFound if not found
// - domain.ErrInternal for database errors
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	query := `
		SELECT
			id,
			email,
			name,
			created_at,
			updated_at
		FROM users
		WHERE email = $1
	`

	var row struct {
		ID        uuid.UUID
		Email     string
		Name      string
		CreatedAt time.Time
		UpdatedAt time.Time
	}

	err := r.db.GetContext(ctx, &row, query, email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("%w: user with email %s", domain.ErrNotFound, email)
		}
		return nil, fmt.Errorf("%w: %v", domain.ErrInternal, err)
	}

	return entities.ReconstituteUser(
		row.ID,
		row.Email,
		row.Name,
		row.CreatedAt,
		row.UpdatedAt,
	)
}

// Helper functions for error mapping

// isUniqueViolation checks if error is a unique constraint violation
//
// PostgreSQL error code 23505: unique_violation
// Uses pgx driver error types for detection
func isUniqueViolation(err error) bool {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// 23505 = unique_violation
		return pgErr.Code == "23505"
	}
	return false
}

// Compile-time verification that implementation satisfies interface
var _ ports.UserRepository = (*UserRepository)(nil)
