// Package middleware provides custom HTTP middleware for the Creator API.
package middleware

import (
	"net/http"

	"github.com/rs/cors"

	"github.com/masquerade/creator-api/internal/config"
)

// CORS creates a CORS middleware handler configured from the application config.
// It allows cross-origin requests from configured origins (typically localhost:3000/3001 in dev).
//
// The middleware handles:
//   - Preflight OPTIONS requests
//   - CORS headers for actual requests
//   - Credential support (cookies, authorization headers)
//
// Configuration comes from config.CORS (AllowedOrigins, AllowedMethods, AllowedHeaders).
func CORS(cfg *config.Config) func(http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.CORS.AllowedOrigins,
		AllowedMethods:   cfg.CORS.AllowedMethods,
		AllowedHeaders:   cfg.CORS.AllowedHeaders,
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes - cache preflight responses
	})

	return c.Handler
}
