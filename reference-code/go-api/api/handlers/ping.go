// Package handlers contains HTTP request handlers for the Creator API.
package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

// PingResponse represents the response structure for the ping endpoint.
type PingResponse struct {
	Message   string `json:"message"`
	Timestamp string `json:"timestamp"`
	Version   string `json:"version"`
}

// HandlePing handles GET /api/v1/ping requests.
// Returns a simple health check response with current timestamp.
func HandlePing(w http.ResponseWriter, r *http.Request) {
	response := PingResponse{
		Message:   "pong",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Version:   "v1",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
