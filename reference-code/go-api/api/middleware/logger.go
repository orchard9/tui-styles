// Package middleware provides HTTP middleware for the Creator API.
package middleware

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog/log"
)

// Logger is a structured logging middleware that logs HTTP requests with context.
// It logs request start (debug level) and completion (info/warn/error based on status).
// Includes request_id, method, path, status, duration, and remote_addr in all logs.
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Extract request ID from context (injected by Chi's RequestID middleware)
		requestID := middleware.GetReqID(r.Context())

		// Log request start at debug level
		log.Debug().
			Str("request_id", requestID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Str("remote_addr", r.RemoteAddr).
			Str("user_agent", r.UserAgent()).
			Msg("Request started")

		// Call next handler
		next.ServeHTTP(ww, r)

		// Calculate request duration
		duration := time.Since(start)

		// Determine log level based on response status
		// 2xx/3xx = info, 4xx = warn, 5xx = error
		logger := log.Info()
		if ww.Status() >= 400 && ww.Status() < 500 {
			logger = log.Warn()
		} else if ww.Status() >= 500 {
			logger = log.Error()
		}

		// Log request completion with structured fields
		logger.
			Str("request_id", requestID).
			Str("method", r.Method).
			Str("path", r.URL.Path).
			Int("status", ww.Status()).
			Dur("duration_ms", duration).
			Int("bytes", ww.BytesWritten()).
			Str("remote_addr", r.RemoteAddr).
			Msg("Request completed")
	})
}
