package handlers_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/orchard9/peach/apps/email-worker/internal/handlers"
)

func TestHealth(t *testing.T) {
	// Create request
	req, err := http.NewRequest("GET", "/health", nil)
	require.NoError(t, err)

	// Create response recorder
	rr := httptest.NewRecorder()

	// Create handler
	handler := handlers.Health()

	// Execute request
	handler.ServeHTTP(rr, req)

	// Check status code
	assert.Equal(t, http.StatusOK, rr.Code)

	// Check response body - verify service and status (uptime varies)
	var response map[string]interface{}
	err = json.Unmarshal(rr.Body.Bytes(), &response)
	require.NoError(t, err)
	assert.Equal(t, "email-worker", response["service"])
	assert.Equal(t, "healthy", response["status"])
	assert.NotNil(t, response["uptime"])
}
