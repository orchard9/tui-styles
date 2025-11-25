package middleware

import (
	"log"
	"net/http"
	"time"
)

// responseWriter wraps http.ResponseWriter to capture the HTTP status code.
// This is necessary because the standard ResponseWriter doesn't expose the
// status code after WriteHeader is called.
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    bool
}

// newResponseWriter creates a new response writer wrapper with default status 200 OK.
func newResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK,
		written:        false,
	}
}

// WriteHeader captures the status code before passing to the underlying writer.
func (rw *responseWriter) WriteHeader(code int) {
	if !rw.written {
		rw.statusCode = code
		rw.written = true
		rw.ResponseWriter.WriteHeader(code)
	}
}

// Write ensures WriteHeader is called with default status if not called explicitly.
func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.written {
		rw.WriteHeader(http.StatusOK)
	}
	return rw.ResponseWriter.Write(b)
}

// ResponseTime measures and logs the response time for each HTTP request.
// It wraps the response writer to capture the status code and logs:
//   - Request ID
//   - HTTP method
//   - URL path
//   - Response status code
//   - Duration (request start to completion)
//
// Log format: [request-id] METHOD /path - STATUS - DURATION
// Example: [abc-123] GET /api/v1/ping - 200 - 2.345ms
//
// This middleware should run after RequestID to include the ID in logs.
// In production, these logs will be replaced with structured logging (task 008)
// and metrics (Prometheus).
func ResponseTime(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		wrapped := newResponseWriter(w)

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Calculate duration
		duration := time.Since(start)

		// Log request completion with timing
		log.Printf("[%s] %s %s - %d - %v",
			GetRequestID(r),
			r.Method,
			r.URL.Path,
			wrapped.statusCode,
			duration,
		)
	})
}
