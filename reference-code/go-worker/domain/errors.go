// Package domain defines domain-level errors.
//
// Domain errors represent business rule violations and domain-specific failures.
// They are:
// - Independent of infrastructure
// - Meaningful to business logic
// - Wrapped with context when propagated
//
// Error Categories:
// - Validation errors: Invalid input, business rule violations
// - Not found errors: Entity doesn't exist
// - Conflict errors: Duplicate, constraint violation
// - Internal errors: Unexpected domain-level failures
//
// Usage:
//
//	if email == "" {
//	    return domain.ErrInvalidEmail("email cannot be empty")
//	}
package domain

import (
	"errors"
	"fmt"
)

// Common domain errors.
//
// These sentinel errors can be checked with errors.Is():
//
//	if errors.Is(err, domain.ErrNotFound) {
//	    // Handle not found
//	}
var (
	// ErrNotFound indicates the requested entity does not exist.
	//
	// Example: User with given ID not found in repository
	ErrNotFound = errors.New("entity not found")

	// ErrConflict indicates a constraint violation or duplicate.
	//
	// Example: User with email already exists
	ErrConflict = errors.New("entity conflict")

	// ErrInternal indicates an unexpected domain-level error.
	//
	// Example: Invariant violation, unexpected state
	ErrInternal = errors.New("internal domain error")

	// ErrValidation indicates input validation failed.
	//
	// Example: Invalid email format, empty required field
	ErrValidation = errors.New("validation failed")
)

// ValidationError represents a validation failure with context.
//
// Provides detailed information about what failed validation.
//
// Example:
//
//	return &ValidationError{
//	    Field:   "email",
//	    Message: "invalid email format",
//	}
type ValidationError struct {
	Field   string
	Message string
}

// Error implements the error interface.
func (e *ValidationError) Error() string {
	if e.Field != "" {
		return fmt.Sprintf("validation error on field '%s': %s", e.Field, e.Message)
	}
	return fmt.Sprintf("validation error: %s", e.Message)
}

// Is allows errors.Is() to match ErrValidation.
func (e *ValidationError) Is(target error) bool {
	return target == ErrValidation
}

// ErrInvalidInput creates a validation error for invalid input.
//
// Use for general input validation failures.
//
// Parameters:
// - message: Description of validation failure
//
// Returns: ValidationError wrapped with context
//
// Example:
//
//	if name == "" {
//	    return ErrInvalidInput("name cannot be empty")
//	}
func ErrInvalidInput(message string) error {
	return &ValidationError{
		Message: message,
	}
}

// ErrInvalidField creates a validation error for specific field.
//
// Use when validation fails on a named field.
//
// Parameters:
// - field: Name of the field that failed validation
// - message: Description of validation failure
//
// Returns: ValidationError with field context
//
// Example:
//
//	return ErrInvalidField("age", "must be positive")
func ErrInvalidField(field, message string) error {
	return &ValidationError{
		Field:   field,
		Message: message,
	}
}

// ErrInvalidEmail creates a validation error for invalid email.
//
// Specific error type for email validation failures.
//
// Parameters:
// - message: Description of why email is invalid
//
// Returns: ValidationError for email field
//
// Example:
//
//	if !isValid(email) {
//	    return ErrInvalidEmail("invalid email format")
//	}
func ErrInvalidEmail(message string) error {
	return &ValidationError{
		Field:   "email",
		Message: message,
	}
}

// NotFoundError represents an entity not found error with context.
//
// Provides information about what entity was not found.
//
// Example:
//
//	return &NotFoundError{
//	    Entity: "User",
//	    ID:     userID.String(),
//	}
type NotFoundError struct {
	Entity string
	ID     string
}

// Error implements the error interface.
func (e *NotFoundError) Error() string {
	if e.ID != "" {
		return fmt.Sprintf("%s with id %s not found", e.Entity, e.ID)
	}
	return fmt.Sprintf("%s not found", e.Entity)
}

// Is allows errors.Is() to match ErrNotFound.
func (e *NotFoundError) Is(target error) bool {
	return target == ErrNotFound
}

// ErrEntityNotFound creates a not found error with context.
//
// Use when an entity lookup fails.
//
// Parameters:
// - entity: Name of the entity type (e.g., "User", "Post")
// - id: Identifier that was not found
//
// Returns: NotFoundError with entity context
//
// Example:
//
//	return ErrEntityNotFound("User", userID.String())
func ErrEntityNotFound(entity, id string) error {
	return &NotFoundError{
		Entity: entity,
		ID:     id,
	}
}

// ConflictError represents a constraint violation or duplicate.
//
// Provides information about what caused the conflict.
//
// Example:
//
//	return &ConflictError{
//	    Message: "email already exists",
//	}
type ConflictError struct {
	Message string
}

// Error implements the error interface.
func (e *ConflictError) Error() string {
	return fmt.Sprintf("conflict: %s", e.Message)
}

// Is allows errors.Is() to match ErrConflict.
func (e *ConflictError) Is(target error) bool {
	return target == ErrConflict
}

// ErrDuplicate creates a conflict error for duplicates.
//
// Use when attempting to create entity that already exists.
//
// Parameters:
// - message: Description of what is duplicated
//
// Returns: ConflictError with context
//
// Example:
//
//	return ErrDuplicate("user with this email already exists")
func ErrDuplicate(message string) error {
	return &ConflictError{
		Message: message,
	}
}
