// Package errors provides comprehensive error handling for email-worker.
// It includes error types, constructors, and HTTP/gRPC error mapping.
package errors

import (
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/orchard9/go-core-http-toolkit/response"
)

// Standard error types for consistent error handling
var (
	// Authentication errors
	ErrAuthentication = errors.New("authentication error")
	ErrInvalidToken   = errors.New("invalid authentication token")
	ErrTokenExpired   = errors.New("authentication token expired")

	// Authorization errors
	ErrUnauthorized     = errors.New("unauthorized")
	ErrPermissionDenied = errors.New("permission denied")

	// Validation errors
	ErrValidation      = errors.New("validation error")
	ErrMissingRequired = errors.New("missing required field")
	ErrInvalidFormat   = errors.New("invalid format")

	// Resource errors
	ErrNotFound      = errors.New("resource not found")
	ErrAlreadyExists = errors.New("resource already exists")
	ErrConflict      = errors.New("resource conflict")

	// System errors
	ErrInternal    = errors.New("internal server error")
	ErrUnavailable = errors.New("service unavailable")
	ErrTimeout     = errors.New("operation timeout")

	// Database errors
	ErrDatabaseConnection = errors.New("database connection error")
	ErrDatabaseQuery      = errors.New("database query error")
)

// DomainError represents an application domain error with context
type DomainError struct {
	Type    error
	Message string
	Details map[string]interface{}
}

// Error implements the error interface
func (e *DomainError) Error() string {
	if e.Message != "" {
		return e.Message
	}
	return e.Type.Error()
}

// Unwrap returns the underlying error type
func (e *DomainError) Unwrap() error {
	return e.Type
}

// Is checks if the error matches the target
func (e *DomainError) Is(target error) bool {
	return errors.Is(e.Type, target)
}

// Wrap creates a new domain error with additional context
func (e *DomainError) Wrap(message string) error {
	return &DomainError{
		Type:    e.Type,
		Message: message,
		Details: e.Details,
	}
}

// WithDetails adds details to the error
func (e *DomainError) WithDetails(details map[string]interface{}) *DomainError {
	e.Details = details
	return e
}

// Error type constructors
var (
	ErrAuth     = &DomainError{Type: ErrAuthentication}
	ErrAuthz    = &DomainError{Type: ErrUnauthorized}
	ErrValid    = &DomainError{Type: ErrValidation}
	ErrResource = &DomainError{Type: ErrNotFound}
	ErrSys      = &DomainError{Type: ErrInternal}
	ErrDB       = &DomainError{Type: ErrDatabaseConnection}
)

// HTTPErrorHandler handles HTTP error responses
type HTTPErrorHandler struct {
	logger *slog.Logger
}

// NewHTTPErrorHandler creates a new HTTP error handler
func NewHTTPErrorHandler(logger *slog.Logger) *HTTPErrorHandler {
	return &HTTPErrorHandler{logger: logger}
}

// HandleError processes domain errors and writes appropriate HTTP responses
func (h *HTTPErrorHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	// Extract request ID for logging
	requestID := r.Header.Get("X-Request-ID")
	if requestID == "" {
		requestID = "unknown"
	}

	var domainErr *DomainError
	if !errors.As(err, &domainErr) {
		// Unknown error - log and return generic internal error
		h.logger.Error("unhandled error",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
		)

		_ = response.Error("internal server error").
			WithStatus(http.StatusInternalServerError).
			WriteTo(w)
		return
	}

	// Map domain errors to HTTP status codes and responses
	var statusCode int
	var message string

	switch {
	case errors.Is(err, ErrAuthentication):
		statusCode = http.StatusUnauthorized
		message = "Authentication required"
	case errors.Is(err, ErrInvalidToken):
		statusCode = http.StatusUnauthorized
		message = "Invalid authentication token"
	case errors.Is(err, ErrTokenExpired):
		statusCode = http.StatusUnauthorized
		message = "Authentication token expired"
	case errors.Is(err, ErrUnauthorized):
		statusCode = http.StatusForbidden
		message = "Access denied"
	case errors.Is(err, ErrPermissionDenied):
		statusCode = http.StatusForbidden
		message = "Permission denied"
	case errors.Is(err, ErrValidation):
		statusCode = http.StatusBadRequest
		message = domainErr.Error() // Use specific validation message
	case errors.Is(err, ErrMissingRequired):
		statusCode = http.StatusBadRequest
		message = domainErr.Error()
	case errors.Is(err, ErrInvalidFormat):
		statusCode = http.StatusBadRequest
		message = domainErr.Error()
	case errors.Is(err, ErrNotFound):
		statusCode = http.StatusNotFound
		message = "Resource not found"
	case errors.Is(err, ErrAlreadyExists):
		statusCode = http.StatusConflict
		message = "Resource already exists"
	case errors.Is(err, ErrConflict):
		statusCode = http.StatusConflict
		message = "Resource conflict"
	case errors.Is(err, ErrUnavailable):
		statusCode = http.StatusServiceUnavailable
		message = "Service unavailable"
	case errors.Is(err, ErrTimeout):
		statusCode = http.StatusGatewayTimeout
		message = "Operation timeout"
	default:
		// Log unexpected domain errors
		h.logger.Error("unexpected domain error",
			slog.String("request_id", requestID),
			slog.String("error", err.Error()),
			slog.String("path", r.URL.Path),
			slog.String("method", r.Method),
		)
		statusCode = http.StatusInternalServerError
		message = "Internal server error"
	}

	// Create error response
	errorResp := map[string]interface{}{
		"error": map[string]interface{}{
			"message":    message,
			"request_id": requestID,
		},
	}

	// Add details if available and not a sensitive error
	if domainErr.Details != nil && statusCode < 500 {
		errorResp["error"].(map[string]interface{})["details"] = domainErr.Details
	}

	_ = response.JSON(errorResp).WithStatus(statusCode).WriteTo(w)
}

// Error constructor functions for elegant error creation

// Authentication error constructors
func Authentication(message string) error {
	return ErrAuth.Wrap(message)
}

func InvalidToken(message string) error {
	return (&DomainError{Type: ErrInvalidToken}).Wrap(message)
}

func TokenExpired(message string) error {
	return (&DomainError{Type: ErrTokenExpired}).Wrap(message)
}

// Authorization error constructors
func Unauthorized(message string) error {
	return ErrAuthz.Wrap(message)
}

func PermissionDenied(message string) error {
	return (&DomainError{Type: ErrPermissionDenied}).Wrap(message)
}

// Validation error constructors
func Validation(message string) error {
	return ErrValid.Wrap(message)
}

func MissingRequired(field string) error {
	return ErrValid.Wrap(fmt.Sprintf("field '%s' is required", field))
}

func InvalidFormat(field, expected string) error {
	return ErrValid.Wrap(fmt.Sprintf("field '%s' has invalid format: expected %s", field, expected))
}

// Resource error constructors
func NotFound(resource string) error {
	return ErrResource.Wrap(fmt.Sprintf("%s not found", resource))
}

func AlreadyExists(resource string) error {
	return (&DomainError{Type: ErrAlreadyExists}).Wrap(fmt.Sprintf("%s already exists", resource))
}

func Conflict(message string) error {
	return (&DomainError{Type: ErrConflict}).Wrap(message)
}

// System error constructors
func Internal(message string) error {
	return ErrSys.Wrap(message)
}

func Unavailable(service string) error {
	return (&DomainError{Type: ErrUnavailable}).Wrap(fmt.Sprintf("%s unavailable", service))
}

func Timeout(operation string) error {
	return (&DomainError{Type: ErrTimeout}).Wrap(fmt.Sprintf("%s timeout", operation))
}

// Database error constructors
func DatabaseConnection(details string) error {
	return ErrDB.Wrap(details)
}

func DatabaseQuery(query string) error {
	return (&DomainError{Type: ErrDatabaseQuery}).Wrap(fmt.Sprintf("query failed: %s", query))
}

// WrapError wraps an error with additional context
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", message, err)
}

// IsRetryable checks if an error is retryable
func IsRetryable(err error) bool {
	if err == nil {
		return false
	}

	switch {
	case errors.Is(err, ErrUnavailable):
		return true
	case errors.Is(err, ErrTimeout):
		return true
	case errors.Is(err, ErrDatabaseConnection):
		return true
	default:
		return false
	}
}
