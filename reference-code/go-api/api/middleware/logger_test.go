package middleware

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	tests := []struct {
		name           string
		statusCode     int
		expectedLevel  string
		expectedFields []string
	}{
		{
			name:          "2xx success",
			statusCode:    http.StatusOK,
			expectedLevel: "info",
			expectedFields: []string{
				"Request completed",
				"request_id",
				"method",
				"path",
				"status",
				"duration_ms",
				"bytes",
				"remote_addr",
			},
		},
		{
			name:          "4xx client error",
			statusCode:    http.StatusBadRequest,
			expectedLevel: "warn",
			expectedFields: []string{
				"Request completed",
				"request_id",
				"method",
				"path",
				"status",
				"400",
			},
		},
		{
			name:          "5xx server error",
			statusCode:    http.StatusInternalServerError,
			expectedLevel: "error",
			expectedFields: []string{
				"Request completed",
				"request_id",
				"method",
				"path",
				"status",
				"500",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

			// Create test handler that returns the specified status code
			handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				_, _ = w.Write([]byte("test response"))
			}))

			// Create test request with request ID in context
			req := httptest.NewRequest("GET", "/test", nil)
			ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "test-request-123")
			req = req.WithContext(ctx)

			w := httptest.NewRecorder()

			// Execute handler
			handler.ServeHTTP(w, req)

			// Verify response
			assert.Equal(t, tt.statusCode, w.Code)

			// Verify logs
			logOutput := buf.String()
			require.NotEmpty(t, logOutput, "Expected log output")

			// Check that expected fields are present in log output
			for _, field := range tt.expectedFields {
				assert.Contains(t, logOutput, field, "Log should contain field: %s", field)
			}

			// Verify log level
			assert.Contains(t, logOutput, tt.expectedLevel, "Log should be at level: %s", tt.expectedLevel)

			// Verify request ID is included
			assert.Contains(t, logOutput, "test-request-123", "Log should include request ID")
		})
	}
}

func TestLoggerWithoutRequestID(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

	// Create test handler
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Create test request WITHOUT request ID
	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute handler
	handler.ServeHTTP(w, req)

	// Verify response
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify logs still work (request_id will be empty string)
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Request completed")
}

func TestLoggerIncludesRequestStart(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).With().Timestamp().Logger()
	zerolog.SetGlobalLevel(zerolog.DebugLevel) // Enable debug logs

	// Create test handler
	handler := Logger(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Create test request with request ID
	req := httptest.NewRequest("POST", "/api/v1/test", nil)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "start-test-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Execute handler
	handler.ServeHTTP(w, req)

	// Verify "Request started" debug log is present
	logOutput := buf.String()
	assert.Contains(t, logOutput, "Request started", "Should log request start")
	assert.Contains(t, logOutput, "start-test-123", "Should include request ID in start log")
	assert.Contains(t, logOutput, "POST", "Should include method in start log")
	assert.Contains(t, logOutput, "/api/v1/test", "Should include path in start log")
}
