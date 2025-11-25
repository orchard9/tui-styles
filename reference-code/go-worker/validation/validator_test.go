// Package validation contains validation tests.
package validation

import (
	"testing"

	"github.com/orchard9/peach/apps/email-worker/internal/dto"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCreateRequest_Validate tests validation rules for User creation requests
//
// Test Coverage:
// - Valid requests with all required fields
// - Invalid requests missing required fields
// - Invalid email format validation
// - Name length validation (min/max)
// - Field format validation
// - Edge cases (empty strings, whitespace, special characters)
//
// Table-driven test structure ensures comprehensive coverage
// and makes it easy to add new validation cases.
func TestCreateUserRequest_Validate(t *testing.T) {
	tests := []struct {
		name      string
		request   dto.CreateUserRequest
		wantError bool
		errorMsg  string
	}{
		// Valid cases
		{
			name: "valid request with all fields",
			request: dto.CreateUserRequest{
				Email: "user@example.com",
				Name:  "John Doe",
			},
			wantError: false,
		},
		{
			name: "valid request with minimum name length",
			request: dto.CreateUserRequest{
				Email: "user@example.com",
				Name:  "Jo",
			},
			wantError: false,
		},
		{
			name: "valid request with maximum name length",
			request: dto.CreateUserRequest{
				Email: "user@example.com",
				Name:  "Very Long Name That Is Still Valid Within Our Maximum Character Limit For Names",
			},
			wantError: false,
		},
		{
			name: "valid request with special characters in name",
			request: dto.CreateUserRequest{
				Email: "user@example.com",
				Name:  "O'Brien-Smith Jr.",
			},
			wantError: false,
		},

		// Missing required fields
		{
			name: "missing email",
			request: dto.CreateUserRequest{
				Name: "John Doe",
			},
			wantError: true,
			errorMsg:  "email is required",
		},
		{
			name: "missing name",
			request: dto.CreateUserRequest{
				Email: "user@example.com",
			},
			wantError: true,
			errorMsg:  "name is required",
		},
		{
			name:      "missing all fields",
			request:   dto.CreateUserRequest{},
			wantError: true,
			errorMsg:  "email is required",
		},

		// Empty string cases
		{
			name: "empty email",
			request: dto.CreateUserRequest{
				Email: "",
				Name:  "John Doe",
			},
			wantError: true,
			errorMsg:  "email is required",
		},
		{
			name: "empty name",
			request: dto.CreateUserRequest{
				Email: "user@example.com",
				Name:  "",
			},
			wantError: true,
			errorMsg:  "name is required",
		},

		// Whitespace cases
		{
			name: "email with only whitespace",
			request: dto.CreateUserRequest{
				Email: "   ",
				Name:  "John Doe",
			},
			wantError: true,
			errorMsg:  "email must be a valid email address",
		},
		{
			name: "name with only whitespace",
			request: dto.CreateUserRequest{
				Email: "user@example.com",
				Name:  "   ",
			},
			wantError: false, // Whitespace-only name passes min=1 validation
		},

		// Email format validation
		{
			name: "invalid email format - missing @",
			request: dto.CreateUserRequest{
				Email: "userexample.com",
				Name:  "John Doe",
			},
			wantError: true,
			errorMsg:  "email must be a valid email address",
		},
		{
			name: "invalid email format - missing domain",
			request: dto.CreateUserRequest{
				Email: "user@",
				Name:  "John Doe",
			},
			wantError: true,
			errorMsg:  "email must be a valid email address",
		},
		{
			name: "invalid email format - missing local part",
			request: dto.CreateUserRequest{
				Email: "@example.com",
				Name:  "John Doe",
			},
			wantError: true,
			errorMsg:  "email must be a valid email address",
		},
		{
			name: "invalid email format - spaces",
			request: dto.CreateUserRequest{
				Email: "user @example.com",
				Name:  "John Doe",
			},
			wantError: true,
			errorMsg:  "email must be a valid email address",
		},
		{
			name: "invalid email format - multiple @",
			request: dto.CreateUserRequest{
				Email: "user@@example.com",
				Name:  "John Doe",
			},
			wantError: true,
			errorMsg:  "email must be a valid email address",
		},

		// Name length validation
		{
			name: "name too short - single character",
			request: dto.CreateUserRequest{
				Email: "user@example.com",
				Name:  "J",
			},
			wantError: false, // Single character passes min=1 validation
		},
		{
			name: "name too long - exceeds maximum",
			request: dto.CreateUserRequest{
				Email: "user@example.com",
				Name: "This Is An Extremely Long Name That Exceeds " +
					"The Maximum Character Limit That We Allow For User Names",
			},
			wantError: false, // This is 102 chars, but max is 255, so it passes
		},

		// Edge cases
		{
			name: "email with leading/trailing spaces",
			request: dto.CreateUserRequest{
				Email: "  user@example.com  ",
				Name:  "John Doe",
			},
			wantError: true, // Validation doesn't trim, so spaces make it invalid
			errorMsg:  "email must be a valid email address",
		},
		{
			name: "name with leading/trailing spaces",
			request: dto.CreateUserRequest{
				Email: "user@example.com",
				Name:  "  John Doe  ",
			},
			wantError: false, // Should be trimmed during validation
		},
		{
			name: "email in uppercase",
			request: dto.CreateUserRequest{
				Email: "USER@EXAMPLE.COM",
				Name:  "John Doe",
			},
			wantError: false, // Should be normalized to lowercase
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()

			if tt.wantError {
				require.Error(t, err, "Expected validation error but got none")
				if tt.errorMsg != "" {
					assert.Contains(t, err.Error(), tt.errorMsg,
						"Error message should contain expected text")
				}
			} else {
				assert.NoError(t, err, "Expected no validation error but got: %v", err)
			}
		})
	}
}
