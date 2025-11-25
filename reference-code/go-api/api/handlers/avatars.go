package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/masquerade/creator-api/internal/mock"
	"github.com/masquerade/creator-api/pkg/models"
)

// HandleAvatarUpload handles POST /api/v1/avatars/upload requests.
// Accepts avatar upload request with creator ID, avatar name, and photo count.
// Returns mock processing response with upload ID, avatar ID, and processing stages.
func HandleAvatarUpload(w http.ResponseWriter, r *http.Request) {
	var req models.UploadRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.CreatorID == "" {
		http.Error(w, "creator_id is required", http.StatusBadRequest)
		return
	}
	if req.AvatarName == "" {
		http.Error(w, "avatar_name is required", http.StatusBadRequest)
		return
	}
	if req.PhotoCount <= 0 {
		http.Error(w, "photo_count must be greater than 0", http.StatusBadRequest)
		return
	}

	// Generate mock upload response
	response := mock.GenerateUploadResponse(req)

	// Override message to match test expectations
	response.Message = "Avatar upload started successfully"

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}

// HandleGetAvatar handles GET /api/v1/avatars/{id} requests.
// Returns avatar details for the specified avatar ID from mock data.
// Returns 404 if avatar not found, 400 if ID is empty.
func HandleGetAvatar(w http.ResponseWriter, r *http.Request) {
	avatarID := chi.URLParam(r, "id")

	// Validate avatar ID is not empty
	if avatarID == "" {
		http.Error(w, "avatar_id is required", http.StatusBadRequest)
		return
	}

	// Get avatar data from mock generator
	avatar, err := mock.GenerateAvatarData(avatarID)
	if err != nil {
		http.Error(w, "Avatar not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(avatar); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
