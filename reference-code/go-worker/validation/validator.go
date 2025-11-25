// Package validation provides request validation utilities for the application layer.
//
// Architecture Position:
// - Lives in application layer (not domain, not adapters)
// - Used by DTOs to validate inputs before reaching use cases
// - Provides infrastructure for validation, NOT business rules
//
// Responsibilities:
// - Structural validation (required fields, format checks)
// - Type validation (email format, UUID format, ranges)
// - Custom validation rules that are NOT business logic
//
// NOT Responsible For:
// - Business rules (e.g., "user must be over 18") → domain layer
// - Database constraints (e.g., "email must be unique") → repository layer
package validation

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator is the global validator instance
//
// Singleton Pattern:
// - One instance per application
// - Thread-safe (validator/v10 is concurrent-safe)
// - Registered custom validators are global
//
// Usage:
//
//	validator := validation.New()
//	if err := validator.Validate(dto); err != nil {
//	    // Handle validation error
//	}
var defaultValidator *Validator

// Validator wraps go-playground/validator with custom error formatting
//
// Design:
// - Thin wrapper around validator/v10
// - Provides consistent error messages
// - Extensible with custom validation rules
type Validator struct {
	validate *validator.Validate
}

// New creates a new validator instance
//
// Call this once during application initialization:
//
//	validator := validation.New()
//
// Thread-Safety: Safe for concurrent use
func New() *Validator {
	v := validator.New()

	// Register custom validators
	registerCustomValidators(v)

	return &Validator{
		validate: v,
	}
}

// GetDefaultValidator returns the global validator instance
//
// Lazy initialization pattern:
// - Created on first use
// - Reused for all subsequent calls
func GetDefaultValidator() *Validator {
	if defaultValidator == nil {
		defaultValidator = New()
	}
	return defaultValidator
}

// Validate validates a struct using validation tags
//
// Parameters:
// - data: Any struct with validation tags
//
// Returns:
// - nil if valid
// - ValidationErrors with field-level details if invalid
//
// Example:
//
//	type CreateUserRequest struct {
//	    Email string `validate:"required,email"`
//	    Name  string `validate:"required,min=2,max=100"`
//	}
//
//	req := CreateUserRequest{Email: "invalid", Name: ""}
//	if err := validator.Validate(req); err != nil {
//	    // err.Error() = "email: must be a valid email address; name: is required"
//	}
func (v *Validator) Validate(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		if validationErrs, ok := err.(validator.ValidationErrors); ok {
			return NewValidationErrors(validationErrs)
		}
		return err
	}
	return nil
}

// ValidationErrors contains field-level validation errors
//
// Provides:
// - Human-readable error messages
// - Field names (JSON tag names)
// - Validation rule that failed
//
// Example Output:
//
//	"email: must be a valid email address; name: is required"
type ValidationErrors struct {
	Errors []FieldError
}

// FieldError represents a single field validation failure
//
// Fields:
// - Field: JSON field name (e.g., "email")
// - Tag: Validation rule that failed (e.g., "required", "email")
// - Value: The invalid value (for debugging)
// - Message: Human-readable error message
type FieldError struct {
	Field   string
	Tag     string
	Value   interface{}
	Message string
}

// NewValidationErrors converts validator.ValidationErrors to our format
func NewValidationErrors(errs validator.ValidationErrors) *ValidationErrors {
	var fieldErrors []FieldError

	for _, err := range errs {
		fieldErrors = append(fieldErrors, FieldError{
			Field:   jsonFieldName(err),
			Tag:     err.Tag(),
			Value:   err.Value(),
			Message: humanReadableMessage(err),
		})
	}

	return &ValidationErrors{Errors: fieldErrors}
}

// Error implements the error interface
//
// Format: "field1: error message; field2: error message"
//
// Example:
//
//	"email: must be a valid email address; name: is required"
func (ve *ValidationErrors) Error() string {
	var messages []string
	for _, err := range ve.Errors {
		messages = append(messages, fmt.Sprintf("%s: %s", err.Field, err.Message))
	}
	return strings.Join(messages, "; ")
}

// Fields returns a map of field names to error messages
//
// Useful for API responses that need field-specific errors:
//
//	{
//	  "errors": {
//	    "email": "must be a valid email address",
//	    "name": "is required"
//	  }
//	}
func (ve *ValidationErrors) Fields() map[string]string {
	fields := make(map[string]string)
	for _, err := range ve.Errors {
		fields[err.Field] = err.Message
	}
	return fields
}

// jsonFieldName extracts the JSON tag name from the field
//
// Example:
//
//	type User struct {
//	    Email string `json:"email" validate:"required"`
//	}
//	// Returns "email" instead of "Email"
func jsonFieldName(fe validator.FieldError) string {
	// Try to get JSON tag from struct field
	field := fe.Field()

	// If no JSON tag, use field name but lowercase first letter
	return strings.ToLower(field[:1]) + field[1:]
}

// humanReadableMessage converts validation tag to human-readable message
//
// Maps technical validation tags to user-friendly messages:
// - "required" → "is required"
// - "email" → "must be a valid email address"
// - "min=2" → "must be at least 2 characters"
func humanReadableMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "is required"
	case "email":
		return "must be a valid email address"
	case "min":
		return fmt.Sprintf("must be at least %s characters", fe.Param())
	case "max":
		return fmt.Sprintf("must be at most %s characters", fe.Param())
	case "len":
		return fmt.Sprintf("must be exactly %s characters", fe.Param())
	case "gt":
		return fmt.Sprintf("must be greater than %s", fe.Param())
	case "gte":
		return fmt.Sprintf("must be greater than or equal to %s", fe.Param())
	case "lt":
		return fmt.Sprintf("must be less than %s", fe.Param())
	case "lte":
		return fmt.Sprintf("must be less than or equal to %s", fe.Param())
	case "uuid":
		return "must be a valid UUID"
	case "uuid4":
		return "must be a valid UUID v4"
	case "oneof":
		return fmt.Sprintf("must be one of: %s", fe.Param())
	case "url":
		return "must be a valid URL"
	case "alpha":
		return "must contain only letters"
	case "alphanum":
		return "must contain only letters and numbers"
	case "numeric":
		return "must be numeric"
	default:
		return fmt.Sprintf("failed validation: %s", fe.Tag())
	}
}

// registerCustomValidators registers custom validation rules
//
// Custom validators for domain-specific needs:
// - Business-agnostic structural validation
// - Format checking
// - NOT business rules
//
// Example:
//
//	validate.RegisterValidation("username", validateUsername)
func registerCustomValidators(v *validator.Validate) {
	// Example: Custom username validation (structural, not business logic)
	// v.RegisterValidation("username", func(fl validator.FieldLevel) bool {
	//     username := fl.Field().String()
	//     // Only alphanumeric and underscore, 3-20 chars
	//     return regexp.MustCompile(`^[a-zA-Z0-9_]{3,20}$`).MatchString(username)
	// })

	// Add more custom validators here as needed
}
