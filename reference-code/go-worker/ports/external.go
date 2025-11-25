package ports

import (
	"context"
	"time"
)

//go:generate mockgen -source=external.go -destination=mocks/external.go -package=mocks

// EmailService defines the contract for sending emails.
//
// This is a PORT (interface) in hexagonal architecture.
// Implementations are ADAPTERS (e.g., SMTP, SendGrid, AWS SES).
//
// Design Principles:
// - No assumptions about email provider
// - Simple interface for common operations
// - Returns domain errors
//
// Example Implementations:
// - adapters/external/email/smtp.go
// - adapters/external/email/sendgrid.go
// - adapters/external/email/mock.go (for testing)
type EmailService interface {
	// Send sends an email.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - to: Recipient email address
	// - subject: Email subject
	// - body: Email body (plain text or HTML)
	//
	// Returns:
	// - error: domain.ErrValidation for invalid email
	// - error: domain.ErrInternal for send failures
	//
	// Example:
	//   err := emailService.Send(ctx, "user@example.com", "Welcome!", "Welcome to our app!")
	Send(ctx context.Context, to, subject, body string) error

	// SendTemplate sends an email using a template.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - to: Recipient email address
	// - templateID: Template identifier
	// - data: Template data
	//
	// Returns:
	// - error: domain.ErrValidation for invalid email/template
	// - error: domain.ErrInternal for send failures
	//
	// Example:
	//   data := map[string]string{"name": "John", "link": "https://..."}
	//   err := emailService.SendTemplate(ctx, "user@example.com", "welcome", data)
	SendTemplate(ctx context.Context, to, templateID string, data map[string]string) error
}

// StorageService defines the contract for file storage.
//
// This is a PORT (interface) in hexagonal architecture.
// Implementations are ADAPTERS (e.g., S3, GCS, local filesystem).
//
// Design Principles:
// - Provider-agnostic interface
// - Support for common storage operations
// - Returns domain errors
//
// Example Implementations:
// - adapters/external/storage/s3.go
// - adapters/external/storage/gcs.go
// - adapters/external/storage/local.go
type StorageService interface {
	// Upload stores a file.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - key: File identifier/path
	// - content: File content
	//
	// Returns:
	// - error: domain.ErrValidation for invalid key
	// - error: domain.ErrInternal for upload failures
	//
	// Example:
	//   err := storage.Upload(ctx, "avatars/user123.jpg", imageData)
	Upload(ctx context.Context, key string, content []byte) error

	// Download retrieves a file.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - key: File identifier/path
	//
	// Returns:
	// - []byte: File content
	// - error: domain.ErrNotFound if file doesn't exist
	// - error: domain.ErrInternal for download failures
	//
	// Example:
	//   content, err := storage.Download(ctx, "avatars/user123.jpg")
	Download(ctx context.Context, key string) ([]byte, error)

	// Delete removes a file.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - key: File identifier/path
	//
	// Returns:
	// - error: domain.ErrInternal for deletion failures
	//
	// Note: Deleting non-existent file is not an error
	//
	// Example:
	//   err := storage.Delete(ctx, "avatars/user123.jpg")
	Delete(ctx context.Context, key string) error

	// GenerateSignedURL creates a temporary URL for accessing a file.
	//
	// Parameters:
	// - ctx: Context for cancellation/timeout
	// - key: File identifier/path
	// - expiration: URL validity duration
	//
	// Returns:
	// - string: Signed URL valid for specified duration
	// - error: domain.ErrNotFound if file doesn't exist
	// - error: domain.ErrInternal for URL generation failures
	//
	// Example:
	//   url, err := storage.GenerateSignedURL(ctx, "private/doc.pdf", 1*time.Hour)
	GenerateSignedURL(ctx context.Context, key string, expiration time.Duration) (string, error)
}
