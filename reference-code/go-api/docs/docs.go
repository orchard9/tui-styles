// Package docs provides embedded API documentation files and handlers.
package docs

import (
	"embed"
	"io/fs"
	"net/http"

	"github.com/rs/zerolog/log"
)

// Embed the OpenAPI specification file
//
//go:embed openapi.yaml
var openAPISpec embed.FS

// Embed Swagger UI static files for interactive documentation
//
//go:embed swagger-ui
var swaggerUIFiles embed.FS

// HandleOpenAPISpec serves the OpenAPI YAML specification file.
// This endpoint allows clients to download the raw spec for code generation,
// testing tools, or other OpenAPI consumers.
//
// Response: application/yaml
// CORS: Enabled for all origins to support browser-based tools.
func HandleOpenAPISpec(w http.ResponseWriter, r *http.Request) {
	// Read the OpenAPI spec from the embedded filesystem
	specData, err := openAPISpec.ReadFile("openapi.yaml")
	if err != nil {
		log.Error().Err(err).Msg("Failed to read OpenAPI spec")
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/yaml")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "public, max-age=3600") // Cache for 1 hour
	w.WriteHeader(http.StatusOK)

	if _, err := w.Write(specData); err != nil {
		log.Error().Err(err).Msg("Failed to write OpenAPI spec")
	}
}

// NewSwaggerUIHandler creates an HTTP handler that serves the Swagger UI
// static files from the embedded filesystem.
//
// The handler strips the "swagger-ui" prefix from the embedded FS
// and serves files directly under the /docs path.
//
// Returns: http.Handler for serving Swagger UI assets.
func NewSwaggerUIHandler() http.Handler {
	// Strip the "swagger-ui" prefix from embedded FS
	swaggerFS, err := fs.Sub(swaggerUIFiles, "swagger-ui")
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to create sub-filesystem for Swagger UI")
	}

	// Create file server for the embedded filesystem
	return http.FileServer(http.FS(swaggerFS))
}
