package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/masquerade/creator-api/internal/mock"
)

// SignupRequest represents the request body for creator signup.
type SignupRequest struct {
	CreatorName string `json:"creator_name"`
	Email       string `json:"email"`
	Portfolio   string `json:"portfolio,omitempty"`
}

// HandleCreatorSignup handles POST /api/v1/creators/signup requests.
// Creates a new creator account with 14-day trial period.
// Returns mock response for development/testing.
func HandleCreatorSignup(w http.ResponseWriter, r *http.Request) {
	var req SignupRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Basic validation
	if req.CreatorName == "" || req.Email == "" {
		http.Error(w, "creator_name and email are required", http.StatusBadRequest)
		return
	}

	// Generate mock response
	response := mock.GenerateSignupResponse(req.CreatorName, req.Email)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// HandleCreatorDashboard handles GET /api/v1/creators/{id}/dashboard requests.
// Returns mock dashboard data with metrics, top avatars, and recent activity.
func HandleCreatorDashboard(w http.ResponseWriter, r *http.Request) {
	creatorID := chi.URLParam(r, "id")

	// Validate creator ID is not empty
	if creatorID == "" {
		http.Error(w, "creator_id is required", http.StatusBadRequest)
		return
	}

	// Generate mock dashboard data
	response := mock.GenerateDashboardData(creatorID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// HandleCreatorEarnings handles GET /api/v1/creators/{id}/earnings requests.
// Returns mock earnings data with balance, breakdown, and payout history.
func HandleCreatorEarnings(w http.ResponseWriter, r *http.Request) {
	creatorID := chi.URLParam(r, "id")

	// Validate creator ID is not empty
	if creatorID == "" {
		http.Error(w, "creator_id is required", http.StatusBadRequest)
		return
	}

	// Generate mock earnings data
	response := mock.GenerateEarningsData(creatorID)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
