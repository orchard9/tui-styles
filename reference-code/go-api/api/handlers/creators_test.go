package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/masquerade/creator-api/internal/mock"
)

func TestHandleCreatorSignup_Success(t *testing.T) {
	// Create test request with valid JSON
	reqBody := `{"creator_name":"Jane Doe","email":"jane@example.com","portfolio":"https://artstation.com/janedoe"}`
	req := httptest.NewRequest("POST", "/api/v1/creators/signup", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	HandleCreatorSignup(w, req)

	// Check status code
	if w.Code != http.StatusCreated {
		t.Errorf("expected status %d, got %d", http.StatusCreated, w.Code)
	}

	// Check content type
	contentType := w.Header().Get("Content-Type")
	if contentType != "application/json" {
		t.Errorf("expected Content-Type application/json, got %s", contentType)
	}

	// Decode response
	var response mock.SignupResponse
	if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	// Validate response fields
	if response.CreatorID == "" {
		t.Error("expected non-empty creator_id")
	}

	if response.Status != "trial" {
		t.Errorf("expected status 'trial', got '%s'", response.Status)
	}

	if !strings.Contains(response.Message, "Jane Doe") {
		t.Errorf("expected message to contain 'Jane Doe', got '%s'", response.Message)
	}

	if response.TrialEndsAt.IsZero() {
		t.Error("expected non-zero trial_ends_at")
	}

	if len(response.NextSteps) != 3 {
		t.Errorf("expected 3 next steps, got %d", len(response.NextSteps))
	}
}

func TestHandleCreatorSignup_InvalidJSON(t *testing.T) {
	// Create test request with invalid JSON
	reqBody := `{"creator_name":"Jane Doe","email":"`
	req := httptest.NewRequest("POST", "/api/v1/creators/signup", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	HandleCreatorSignup(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleCreatorSignup_MissingCreatorName(t *testing.T) {
	// Create test request with missing creator_name
	reqBody := `{"email":"jane@example.com"}`
	req := httptest.NewRequest("POST", "/api/v1/creators/signup", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	HandleCreatorSignup(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleCreatorSignup_MissingEmail(t *testing.T) {
	// Create test request with missing email
	reqBody := `{"creator_name":"Jane Doe"}`
	req := httptest.NewRequest("POST", "/api/v1/creators/signup", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	HandleCreatorSignup(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestHandleCreatorSignup_EmptyBody(t *testing.T) {
	// Create test request with empty body
	req := httptest.NewRequest("POST", "/api/v1/creators/signup", bytes.NewReader([]byte{}))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Call handler
	HandleCreatorSignup(w, req)

	// Check status code
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected status %d, got %d", http.StatusBadRequest, w.Code)
	}
}
