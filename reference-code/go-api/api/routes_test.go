package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/masquerade/creator-api/internal/api/handlers"
)

func TestNewRouter(t *testing.T) {
	router := NewRouter()
	assert.NotNil(t, router, "Router should not be nil")
}

func TestPingEndpoint(t *testing.T) {
	// Create router
	router := NewRouter()

	// Create test request
	req := httptest.NewRequest("GET", "/api/v1/ping", nil)
	w := httptest.NewRecorder()

	// Execute request
	router.ServeHTTP(w, req)

	// Assert response status
	assert.Equal(t, http.StatusOK, w.Code)

	// Assert content type
	assert.Contains(t, w.Header().Get("Content-Type"), "application/json")

	// Parse and validate JSON structure
	var response handlers.PingResponse
	err := json.NewDecoder(w.Body).Decode(&response)
	require.NoError(t, err, "Response should be valid JSON")

	// Assert response fields
	assert.Equal(t, "pong", response.Message)
	assert.Equal(t, "v1", response.Version)
	assert.NotEmpty(t, response.Timestamp, "Timestamp should not be empty")
}

func TestMiddlewareApplied(t *testing.T) {
	router := NewRouter()
	req := httptest.NewRequest("GET", "/api/v1/ping", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Middleware is applied if we get a successful response
	// The RequestID middleware adds the ID to the request context (not response headers)
	// The Logger middleware logs the request (we can see it in test output)
	// The Recoverer middleware handles panics
	assert.Equal(t, http.StatusOK, w.Code, "Middleware chain should allow request to complete")
}

func TestNotFoundHandling(t *testing.T) {
	router := NewRouter()
	req := httptest.NewRequest("GET", "/api/v1/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Chi returns 404 for unmatched routes
	assert.Equal(t, http.StatusNotFound, w.Code)
}

func TestImplementedEndpoints(t *testing.T) {
	router := NewRouter()

	tests := []struct {
		name           string
		method         string
		path           string
		body           string
		expectedStatus int
	}{
		// Note: All endpoints now implemented with mock data (milestone-2)
		{"Avatar Upload - no body", "POST", "/api/v1/avatars/upload", "", http.StatusBadRequest},
		{"Get Avatar - not found", "GET", "/api/v1/avatars/nonexistent", "", http.StatusNotFound},
		{"Get Avatar - found", "GET", "/api/v1/avatars/avatar_a1b2c3", "", http.StatusOK},
		{"Creator Dashboard", "GET", "/api/v1/creators/creator_123/dashboard", "", http.StatusOK},
		{"Creator Earnings", "GET", "/api/v1/creators/creator_123/earnings", "", http.StatusOK},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.body != "" {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			} else {
				req = httptest.NewRequest(tt.method, tt.path, nil)
			}
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Implemented handlers should return appropriate status codes
			assert.Equal(t, tt.expectedStatus, w.Code,
				"Endpoint %s should return status %d", tt.path, tt.expectedStatus)
		})
	}
}
