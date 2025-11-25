// Package entities contains domain entities with business logic.
//
// Domain entities represent core business concepts and encapsulate business rules.
// They have:
// - Private fields (encapsulation)
// - Public getters (read access)
// - Business logic methods (behavior)
// - No external dependencies (pure domain)
//
// Rules:
// - Entities are ALWAYS created through constructor functions
// - All business rules enforced at creation and modification time
// - No setters - use named methods that express business operations
// - Methods return domain errors, not infrastructure errors
package entities

import (
	"time"

	"github.com/google/uuid"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
)

// User represents a user in the system.
//
// Business Rules:
// - User must have a valid email address
// - User name cannot be empty
// - User is identified by UUID
//
// Invariants:
// - ID is immutable after creation
// - Email must remain valid throughout lifecycle
// - Timestamps track creation and updates
type User struct {
	id        uuid.UUID
	email     Email // Value object for type safety
	name      string
	createdAt time.Time
	updatedAt time.Time
}

// NewUser creates a new User with validation.
//
// Business Rules Enforced:
// - Email must be valid format
// - Name must not be empty
// - Timestamps set to current time
//
// Returns:
// - *User: Valid entity ready for persistence
// - error: domain.ErrInvalidEmail or domain.ErrInvalidInput
//
// Example:
//
//	user, err := entities.NewUser("john@example.com", "John Doe")
//	if err != nil {
//	    // Handle validation error
//	}
func NewUser(emailStr, name string) (*User, error) {
	// Validate email using value object
	email, err := NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	// Validate name
	if name == "" {
		return nil, domain.ErrInvalidInput("name cannot be empty")
	}

	now := time.Now()
	return &User{
		id:        uuid.New(),
		email:     email,
		name:      name,
		createdAt: now,
		updatedAt: now,
	}, nil
}

// ReconstituteUser recreates a User from persistence.
//
// Use this when loading entities from database.
// Skips validation since data is already validated.
//
// Parameters:
// - id: Entity UUID from database
// - email: Email address from database
// - name: Name from database
// - createdAt: Creation timestamp
// - updatedAt: Last update timestamp
//
// Returns: Fully reconstituted entity
func ReconstituteUser(id uuid.UUID, emailStr, name string, createdAt, updatedAt time.Time) (*User, error) {
	email, err := NewEmail(emailStr)
	if err != nil {
		return nil, err
	}

	return &User{
		id:        id,
		email:     email,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}, nil
}

// Getters - Provide read access to private fields

// ID returns the entity's unique identifier.
func (e *User) ID() uuid.UUID {
	return e.id
}

// Email returns the entity's email address.
func (e *User) Email() string {
	return e.email.String()
}

// Name returns the entity's name.
func (e *User) Name() string {
	return e.name
}

// CreatedAt returns when the entity was created.
func (e *User) CreatedAt() time.Time {
	return e.createdAt
}

// UpdatedAt returns when the entity was last updated.
func (e *User) UpdatedAt() time.Time {
	return e.updatedAt
}

// Business Logic Methods - Express domain operations

// ChangeEmail updates the entity's email address.
//
// Business Rule: New email must be valid format.
//
// Parameters:
// - newEmail: New email address to set
//
// Returns:
// - error: domain.ErrInvalidEmail if invalid format
//
// Side Effects:
// - Updates updatedAt timestamp
func (e *User) ChangeEmail(newEmail string) error {
	email, err := NewEmail(newEmail)
	if err != nil {
		return err
	}

	e.email = email
	e.updatedAt = time.Now()
	return nil
}

// ChangeName updates the entity's name.
//
// Business Rule: Name cannot be empty.
//
// Parameters:
// - newName: New name to set
//
// Returns:
// - error: domain.ErrInvalidInput if name is empty
//
// Side Effects:
// - Updates updatedAt timestamp
func (e *User) ChangeName(newName string) error {
	if newName == "" {
		return domain.ErrInvalidInput("name cannot be empty")
	}

	e.name = newName
	e.updatedAt = time.Now()
	return nil
}

// Update marks the entity as updated.
//
// Use this after any business operation that modifies state.
func (e *User) Update() {
	e.updatedAt = time.Now()
}
