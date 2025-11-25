// Package handlers contains HTTP request handlers for the service.
// It includes base handler functionality, validation helpers, and error handling.
package handlers

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/orchard9/go-core-http-toolkit/middleware/auth"
	"github.com/orchard9/go-core-http-toolkit/response"
	"github.com/orchard9/peach/apps/email-worker/internal/errors"
)

// BaseHandler provides common functionality for all HTTP handlers
type BaseHandler struct {
	logger *slog.Logger
}

// NewBaseHandler creates a new base handler
func NewBaseHandler(logger *slog.Logger) *BaseHandler {
	return &BaseHandler{
		logger: logger,
	}
}

// Logger returns the handler's logger
func (h *BaseHandler) Logger() *slog.Logger {
	return h.logger
}

// GetAccountIDFromContext extracts the authenticated account ID from context
func (h *BaseHandler) GetAccountIDFromContext(r *http.Request) (uuid.UUID, error) {
	claims, ok := auth.GetClaims(r)
	if !ok {
		return uuid.Nil, errors.Authentication("no authentication claims found")
	}

	accountID, err := uuid.Parse(claims.Subject)
	if err != nil {
		return uuid.Nil, errors.Authentication("invalid account ID format")
	}

	return accountID, nil
}

// RequireAuthentication extracts account ID and returns error if not authenticated
func (h *BaseHandler) RequireAuthentication(w http.ResponseWriter, r *http.Request) (uuid.UUID, bool) {
	accountID, err := h.GetAccountIDFromContext(r)
	if err != nil {
		h.HandleError(w, r, err)
		return uuid.Nil, false
	}
	return accountID, true
}

// HandleError processes domain errors and writes appropriate HTTP responses
func (h *BaseHandler) HandleError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}

	// Create error handler and use centralized error mapping
	errorHandler := errors.NewHTTPErrorHandler(h.logger)
	errorHandler.HandleError(w, r, err)
}

// ValidateUUID validates and parses a UUID string
func (h *BaseHandler) ValidateUUID(id, fieldName string) (uuid.UUID, error) {
	if id == "" {
		return uuid.Nil, errors.MissingRequired(fieldName)
	}

	parsed, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil, errors.InvalidFormat(fieldName, "UUID")
	}

	return parsed, nil
}

// ValidateRequiredString validates that a string field is not empty
func (h *BaseHandler) ValidateRequiredString(value, fieldName string) error {
	if value == "" {
		return errors.MissingRequired(fieldName)
	}
	return nil
}

// ValidateRequiredSlice validates that a slice field is not empty
func (h *BaseHandler) ValidateRequiredSlice(value []string, fieldName string) error {
	if len(value) == 0 {
		return errors.MissingRequired(fieldName)
	}
	return nil
}

// ParseInt safely parses an integer from string with validation
func (h *BaseHandler) ParseInt(value, fieldName string) (int, error) {
	if value == "" {
		return 0, errors.MissingRequired(fieldName)
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return 0, errors.InvalidFormat(fieldName, "integer")
	}

	return parsed, nil
}

// ParseOptionalInt safely parses an optional integer from string
func (h *BaseHandler) ParseOptionalInt(value string, defaultValue int) int {
	if value == "" {
		return defaultValue
	}

	parsed, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}

	return parsed
}

// DefaultPageSize returns a safe page size, applying defaults and limits
func (h *BaseHandler) DefaultPageSize(requested, defaultSize, maxSize int) int {
	if requested <= 0 {
		return defaultSize
	}
	if requested > maxSize {
		return maxSize
	}
	return requested
}

// RespondJSON writes a JSON response with proper error handling
func (h *BaseHandler) RespondJSON(w http.ResponseWriter, r *http.Request, data interface{}) {
	resp := response.JSON(data)
	_ = resp.WriteTo(w)
}

// RespondJSONWithStatus writes a JSON response with a specific status code
func (h *BaseHandler) RespondJSONWithStatus(w http.ResponseWriter, r *http.Request, status int, data interface{}) {
	resp := response.JSON(data).WithStatus(status)
	_ = resp.WriteTo(w)
}

// RespondError writes an error response
func (h *BaseHandler) RespondError(w http.ResponseWriter, r *http.Request, message string, status int) {
	resp := response.Error(message).WithStatus(status)
	_ = resp.WriteTo(w)
}
