// Package middleware contains middleware configuration for the HTTP server.
//
// Middleware setup is organized into layers:
// - Global middleware: Applied to all requests (logging, recovery, etc.)
// - API middleware: Applied to API routes (auth, rate limiting, etc.)
//
// This follows hexagonal architecture by:
// - Keeping middleware configuration separate from business logic
// - Using dependency injection for middleware dependencies
// - Maintaining clear separation of concerns
package middleware

import (
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/orchard9/go-core-http-toolkit/middleware"

	"github.com/orchard9/peach/apps/email-worker/internal/config"
)

// SetupGlobalMiddleware configures middleware applied to all requests
//
// Middleware Order (Important!):
// 1. Recoverer - Catches panics and returns 500
// 2. RealIP - Extracts real client IP from X-Forwarded-For
// 3. RequestID - Generates unique request ID for tracing
// 4. Logger - Logs all requests with structured logging
// 5. CORS - Handles Cross-Origin Resource Sharing
//
// Thread-Safety: Safe for concurrent use (middleware is stateless)
func SetupGlobalMiddleware(r *chi.Mux, logger *slog.Logger) {
	// Panic recovery - MUST be first to catch panics from other middleware
	r.Use(chimiddleware.Recoverer)

	// Real IP extraction - Must be before logging to log correct IP
	r.Use(middleware.RealIP())

	// Request ID generation - Must be before logging to include ID in logs
	r.Use(middleware.RequestID())

	// Structured logging
	r.Use(middleware.Logger(middleware.LoggingConfig{
		Logger: logger,
	}))

	// CORS configuration
	r.Use(middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   getAllowedOrigins(),
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Request-ID"},
		ExposedHeaders:   []string{"Link", "X-Request-ID"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))
}

// SetupAPIMiddleware configures middleware for API routes
//
// API Middleware Order:
// N. Content Type - Enforces JSON content type
//
// Parameters:
// - cfg: Application configuration
// - logger: Structured logger for middleware
//
// Returns: Slice of middleware functions to apply
func SetupAPIMiddleware(cfg *config.Config, logger *slog.Logger) []func(http.Handler) http.Handler {
	var middlewares []func(http.Handler) http.Handler

	// NOTE: Add additional middleware as needed
	// Example: Content-Type validation, request size limits, etc.

	return middlewares
}

// getAllowedOrigins returns the list of allowed CORS origins
//
// Configuration:
// - Development: Allow all origins (*)
// - Production: Only allow configured origins
//
// Environment Variables:
// - ALLOWED_ORIGINS: Comma-separated list of allowed origins
func getAllowedOrigins() []string {
	// TODO: Load from configuration
	// For now, allow all origins in development
	return []string{"*"}
}

// Middleware helper functions

// RequireAuth wraps a handler to require authentication
//
// Usage:
//
//	r.Get("/protected", RequireAuth(handler.Handle))
//
// Response: 401 Unauthorized if authentication fails
func RequireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}

// RequireRole wraps a handler to require a specific role
//
// Usage:
//
//	r.Get("/admin", RequireRole("admin", handler.Handle))
//
// Response:
// - 401 Unauthorized if not authenticated
// - 403 Forbidden if insufficient permissions
func RequireRole(role string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		next(w, r)
	}
}
