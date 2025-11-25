// Package api provides HTTP routing and request handling for the Creator API.
package api

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/masquerade/creator-api/internal/api/handlers"
	"github.com/masquerade/creator-api/internal/docs"
)

// NewRouter creates and configures a Chi router with all API routes and middleware.
// The router includes global middleware (RequestID, RealIP, Logger, Recoverer) and
// versioned API routes under /api/v1.
func NewRouter() *chi.Mux {
	r := chi.NewRouter()

	// Global middleware (applied to all routes)
	r.Use(middleware.RequestID) // Inject request ID into context
	r.Use(middleware.RealIP)    // Set RemoteAddr to real IP
	r.Use(middleware.Logger)    // Log all requests
	r.Use(middleware.Recoverer) // Recover from panics

	// Health check endpoint (unversioned for infrastructure)
	r.Get("/health", handlers.HandleHealth)

	// Documentation endpoints (no auth required)
	r.Get("/docs/openapi.yaml", docs.HandleOpenAPISpec)
	r.Handle("/docs/*", http.StripPrefix("/docs", docs.NewSwaggerUIHandler()))

	// API v1 routes (versioned for future compatibility)
	r.Route("/api/v1", func(r chi.Router) {
		// Health/test endpoint
		r.Get("/ping", handlers.HandlePing)

		// Creator endpoints (stubs for task 004+)
		r.Post("/creators/signup", handlers.HandleCreatorSignup)

		// Avatar endpoints (stubs for task 006)
		r.Post("/avatars/upload", handlers.HandleAvatarUpload)
		r.Get("/avatars/{id}", handlers.HandleGetAvatar)

		// Dashboard/analytics (stubs for task 006)
		r.Get("/creators/{id}/dashboard", handlers.HandleCreatorDashboard)
		r.Get("/creators/{id}/earnings", handlers.HandleCreatorEarnings)
	})

	return r
}
