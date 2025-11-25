// Package handlers contains HTTP request handlers for the Creator API.
package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// HealthResponse represents the response structure for the health check endpoint.
// Used by infrastructure for liveness/readiness probes.
type HealthResponse struct {
	Status    string    `json:"status"`
	Version   string    `json:"version"`
	Timestamp time.Time `json:"timestamp"`
	Service   string    `json:"service"`
}

// HandleHealth handles GET /health requests.
// Returns service status, version, current timestamp, and service name.
// This endpoint is used by Kubernetes/Docker for health checks.
func HandleHealth(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Status:    "healthy",
		Version:   "0.1.0",
		Timestamp: time.Now().UTC(),
		Service:   "creator-api",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
