// Package entities contains value objects for type safety and validation.
//
// Value Objects:
// - Immutable after creation
// - Validated at construction time
// - Compared by value, not identity
// - No setters - create new instance instead
//
// Example: Email is a value object ensuring valid email format
package entities

import (
	"regexp"
	"strings"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
)

// Email represents a validated email address.
//
// Value Object Properties:
// - Immutable: Cannot be changed after creation
// - Always valid: Invalid emails cannot be constructed
// - Value equality: Two emails with same string are equal
//
// Business Rules:
// - Must match standard email format
// - Stored in lowercase for consistency
//
// Usage:
//
//	email, err := entities.NewEmail("user@example.com")
//	if err != nil {
//	    // Handle invalid email
//	}
//	fmt.Println(email.String()) // "user@example.com"
type Email struct {
	value string
}

// emailRegex validates email format.
//
// Pattern matches standard email addresses:
// - Local part: letters, numbers, dots, underscores, hyphens, plus signs
// - @ symbol
// - Domain: letters, numbers, dots, hyphens
// - TLD: at least 2 letters
var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

// NewEmail creates a validated Email value object.
//
// Validation Rules:
// - Must not be empty
// - Must match email regex pattern
// - Normalized to lowercase
//
// Parameters:
// - value: Email address string to validate
//
// Returns:
// - Email: Valid email value object
// - error: domain.ErrInvalidEmail if validation fails
//
// Example:
//
//	email, err := NewEmail("User@Example.COM")
//	// email.String() returns "user@example.com"
func NewEmail(value string) (Email, error) {
	// Trim whitespace
	value = strings.TrimSpace(value)

	// Check empty
	if value == "" {
		return Email{}, domain.ErrInvalidEmail("email cannot be empty")
	}

	// Normalize to lowercase
	value = strings.ToLower(value)

	// Validate format
	if !emailRegex.MatchString(value) {
		return Email{}, domain.ErrInvalidEmail("invalid email format")
	}

	return Email{value: value}, nil
}

// String returns the email address as a string.
//
// Returns: Lowercase email address
func (e Email) String() string {
	return e.value
}

// Equals checks if two emails are equal.
//
// Value objects are compared by value, not reference.
//
// Parameters:
// - other: Email to compare with
//
// Returns: true if email addresses are identical
func (e Email) Equals(other Email) bool {
	return e.value == other.value
}
