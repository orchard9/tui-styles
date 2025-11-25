// Package memory provides in-memory repository implementations for testing.
//
// Purpose:
// - Fast unit/integration testing without database
// - Deterministic test behavior
// - No external dependencies
//
// NOT for production use - data is not persisted
package memory

import (
	"context"
	"fmt"
	"sort"
	"sync"

	"github.com/google/uuid"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
	"github.com/orchard9/peach/apps/email-worker/internal/domain/entities"
	"github.com/orchard9/peach/apps/email-worker/internal/ports"
)

// UserRepository implements ports.UserRepository in memory
//
// Thread-Safety: Safe for concurrent use (uses mutex)
//
// Use Cases:
// - Unit testing use cases
// - Integration testing without database
// - Development without PostgreSQL
type UserRepository struct {
	mu    sync.RWMutex
	data  map[uuid.UUID]*entities.User
	email map[string]uuid.UUID // Index by email for fast lookup
}

// NewInMemoryUserRepository creates a new in-memory repository
//
// Returns an empty repository ready for use
func NewInMemoryUserRepository() ports.UserRepository {
	return &UserRepository{
		data:  make(map[uuid.UUID]*entities.User),
		email: make(map[string]uuid.UUID),
	}
}

// Create stores a User in memory
//
// Validates:
// - Email uniqueness
//
// Returns:
// - domain.ErrConflict if email already exists
func (r *UserRepository) Create(ctx context.Context, entity *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check for duplicate email
	if _, exists := r.email[entity.Email()]; exists {
		return fmt.Errorf("%w: email already exists", domain.ErrConflict)
	}

	// Store entity
	r.data[entity.ID()] = entity
	r.email[entity.Email()] = entity.ID()

	return nil
}

// FindByID retrieves a User by ID
//
// Returns:
// - domain.ErrNotFound if entity doesn't exist
func (r *UserRepository) FindByID(ctx context.Context, id uuid.UUID) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	entity, exists := r.data[id]
	if !exists {
		return nil, fmt.Errorf("%w: user with id %s", domain.ErrNotFound, id)
	}

	return entity, nil
}

// Update modifies an existing User
//
// Validates:
// - Entity exists
// - Email uniqueness (if changed)
//
// Returns:
// - domain.ErrNotFound if entity doesn't exist
// - domain.ErrConflict if email already exists
func (r *UserRepository) Update(ctx context.Context, entity *entities.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	existing, exists := r.data[entity.ID()]
	if !exists {
		return fmt.Errorf("%w: user with id %s", domain.ErrNotFound, entity.ID())
	}

	// Check if email changed and is duplicate
	if existing.Email() != entity.Email() {
		if _, emailExists := r.email[entity.Email()]; emailExists {
			return fmt.Errorf("%w: email already exists", domain.ErrConflict)
		}
		// Update email index
		delete(r.email, existing.Email())
		r.email[entity.Email()] = entity.ID()
	}

	r.data[entity.ID()] = entity
	return nil
}

// Delete removes a User
//
// Idempotent: No error if entity doesn't exist
func (r *UserRepository) Delete(ctx context.Context, id uuid.UUID) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	entity, exists := r.data[id]
	if exists {
		delete(r.email, entity.Email())
		delete(r.data, id)
	}

	// Idempotent - no error if not exists
	return nil
}

// List returns all Users with pagination
//
// Sorting: Ordered by created_at DESC (newest first)
//
// Returns:
// - Empty slice if no results (not an error)
// - Applies pagination after sorting
func (r *UserRepository) List(ctx context.Context, opts ports.ListOptions) ([]*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	// Extract limit and offset from options
	limit := opts.Limit
	if limit <= 0 {
		limit = 20 // Default
	}
	offset := opts.Offset
	if offset < 0 {
		offset = 0
	}

	// Collect all entities
	all := make([]*entities.User, 0, len(r.data))
	for _, entity := range r.data {
		all = append(all, entity)
	}

	// Sort by created_at DESC (newest first, matching PostgreSQL ORDER BY created_at DESC)
	sort.Slice(all, func(i, j int) bool {
		return all[i].CreatedAt().After(all[j].CreatedAt())
	})

	// Apply pagination
	start := offset
	if start >= len(all) {
		return []*entities.User{}, nil
	}

	end := start + limit
	if end > len(all) {
		end = len(all)
	}

	return all[start:end], nil
}

// FindByEmail retrieves a User by email
//
// Returns:
// - domain.ErrNotFound if no entity with that email
func (r *UserRepository) FindByEmail(ctx context.Context, email string) (*entities.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	id, exists := r.email[email]
	if !exists {
		return nil, fmt.Errorf("%w: user with email %s", domain.ErrNotFound, email)
	}

	return r.data[id], nil
}

// Clear removes all data (useful for test cleanup between test cases)
//
// Not part of ports.UserRepository interface - test helper only
func (r *UserRepository) Clear() {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.data = make(map[uuid.UUID]*entities.User)
	r.email = make(map[string]uuid.UUID)
}

// Compile-time verification that implementation satisfies interface
var _ ports.UserRepository = (*UserRepository)(nil)
