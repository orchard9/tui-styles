package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

// ContextKey is a custom type for context keys to avoid collisions.
type ContextKey string

const (
	// RequestIDKey is the context key for request IDs.
	RequestIDKey ContextKey = "request_id"
)

// RequestID generates or forwards a unique request ID for each HTTP request.
// It checks for an existing X-Request-ID header and uses it if present,
// otherwise generates a new UUID.
//
// The request ID is:
//   - Added to the request context (accessible via GetRequestID)
//   - Added to response headers (X-Request-ID)
//   - Used by other middleware for logging and error tracking
//
// This middleware should be first in the chain so all subsequent middleware
// and handlers have access to the request ID.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if request already has an ID (from upstream proxy or client)
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			// Generate new UUID if not provided
			requestID = uuid.New().String()
		}

		// Add to context for handlers to access
		ctx := context.WithValue(r.Context(), RequestIDKey, requestID)

		// Add to response headers so clients can track requests
		w.Header().Set("X-Request-ID", requestID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetRequestID extracts the request ID from the request context.
// Returns "unknown" if the request ID is not found (should never happen if
// RequestID middleware is used correctly).
func GetRequestID(r *http.Request) string {
	if id, ok := r.Context().Value(RequestIDKey).(string); ok {
		return id
	}
	return "unknown"
}
