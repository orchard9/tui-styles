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

func TestRecoverer(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

	// Create handler that panics
	handler := Recoverer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	}))

	// Create test request with request ID
	req := httptest.NewRequest("GET", "/test", nil)
	ctx := context.WithValue(req.Context(), middleware.RequestIDKey, "panic-test-123")
	req = req.WithContext(ctx)

	w := httptest.NewRecorder()

	// Execute handler (should not crash)
	handler.ServeHTTP(w, req)

	// Verify response is 500
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal Server Error")

	// Verify panic was logged
	logOutput := buf.String()
	require.NotEmpty(t, logOutput, "Expected log output from panic")

	assert.Contains(t, logOutput, "error", "Should log at error level")
	assert.Contains(t, logOutput, "Panic recovered", "Should log panic message")
	assert.Contains(t, logOutput, "panic-test-123", "Should include request ID")
	assert.Contains(t, logOutput, "test panic", "Should include panic value")
	assert.Contains(t, logOutput, "stack_trace", "Should include stack trace field")
}

func TestRecovererWithNoPanic(t *testing.T) {
	// Capture log output
	var buf bytes.Buffer
	log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

	// Create handler that doesn't panic
	handler := Recoverer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Execute handler
	handler.ServeHTTP(w, req)

	// Verify normal response
	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())

	// Verify no panic logs
	logOutput := buf.String()
	assert.NotContains(t, logOutput, "Panic recovered", "Should not log panic for normal request")
}

func TestRecovererWithDifferentPanicTypes(t *testing.T) {
	tests := []struct {
		name       string
		panicValue interface{}
	}{
		{
			name:       "panic with string",
			panicValue: "string panic",
		},
		{
			name:       "panic with error",
			panicValue: assert.AnError,
		},
		{
			name:       "panic with int",
			panicValue: 42,
		},
		{
			name:       "panic with struct",
			panicValue: struct{ msg string }{msg: "struct panic"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Capture log output
			var buf bytes.Buffer
			log.Logger = zerolog.New(&buf).With().Timestamp().Logger()

			// Create handler that panics with specific value
			handler := Recoverer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(tt.panicValue)
			}))

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			// Execute handler (should not crash)
			handler.ServeHTTP(w, req)

			// Verify response is 500
			assert.Equal(t, http.StatusInternalServerError, w.Code)

			// Verify panic was logged
			logOutput := buf.String()
			assert.Contains(t, logOutput, "Panic recovered")
			assert.NotEmpty(t, logOutput)
		})
	}
}
