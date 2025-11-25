package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/masquerade/creator-api/internal/config"
)

// TestRequestID verifies that the RequestID middleware generates a unique ID.
func TestRequestID(t *testing.T) {
	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request ID is in context
		requestID := GetRequestID(r)
		assert.NotEmpty(t, requestID, "Request ID should be in context")
		assert.NotEqual(t, "unknown", requestID, "Request ID should not be 'unknown'")
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify response status
	assert.Equal(t, http.StatusOK, w.Code)

	// Verify X-Request-ID header is set in response
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"), "X-Request-ID header should be set")
}

// TestRequestIDForwarding verifies that existing X-Request-ID is preserved.
func TestRequestIDForwarding(t *testing.T) {
	existingID := "test-request-id-12345"

	handler := RequestID(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify existing request ID is preserved in context
		requestID := GetRequestID(r)
		assert.Equal(t, existingID, requestID, "Existing request ID should be preserved")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("X-Request-ID", existingID)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify same ID in response header
	assert.Equal(t, existingID, w.Header().Get("X-Request-ID"))
}

// TestGetRequestIDUnknown verifies GetRequestID returns "unknown" when not set.
func TestGetRequestIDUnknown(t *testing.T) {
	req := httptest.NewRequest("GET", "/test", nil)
	requestID := GetRequestID(req)
	assert.Equal(t, "unknown", requestID, "Should return 'unknown' when not set")
}

// TestRecovery verifies panic recovery and error response.
func TestRecovery(t *testing.T) {
	handler := Recoverer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic - this should be caught")
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Should not panic
	handler.ServeHTTP(w, req)

	// Verify 500 status code
	assert.Equal(t, http.StatusInternalServerError, w.Code)

	// Verify response contains error message (plain text from http.Error)
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}

// TestRecoveryWithRequestID verifies panic recovery works with RequestID middleware.
func TestRecoveryWithRequestID(t *testing.T) {
	// Stack middleware: RequestID -> Recoverer -> Panic handler
	handler := RequestID(Recoverer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		panic("test panic")
	})))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify 500 status code and recovery worked
	assert.Equal(t, http.StatusInternalServerError, w.Code)
	assert.Contains(t, w.Body.String(), "Internal Server Error")
}

// TestRecoveryNoError verifies normal request flow (no panic).
func TestRecoveryNoError(t *testing.T) {
	handler := Recoverer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
}

// TestCORS verifies CORS middleware configuration.
func TestCORS(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedMethods: []string{"GET", "POST"},
			AllowedHeaders: []string{"Content-Type", "Authorization"},
		},
	}

	handler := CORS(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	// Test preflight OPTIONS request
	req := httptest.NewRequest("OPTIONS", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Verify CORS headers are set
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
	assert.Contains(t, w.Header().Get("Access-Control-Allow-Methods"), "POST")
}

// TestCORSActualRequest verifies CORS headers on actual requests.
func TestCORSActualRequest(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedMethods: []string{"GET", "POST"},
			AllowedHeaders: []string{"Content-Type"},
		},
	}

	handler := CORS(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("response"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
}

// TestCORSDisallowedOrigin verifies CORS blocks non-allowed origins.
func TestCORSDisallowedOrigin(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedMethods: []string{"GET"},
			AllowedHeaders: []string{"Content-Type"},
		},
	}

	handler := CORS(cfg)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://evil.com")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// CORS should not set Access-Control-Allow-Origin for disallowed origins
	assert.NotEqual(t, "http://evil.com", w.Header().Get("Access-Control-Allow-Origin"))
}

// TestResponseTime verifies response time measurement.
func TestResponseTime(t *testing.T) {
	handler := ResponseTime(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	// Should complete without error and log timing
	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}

// TestResponseTimeWithRequestID verifies logging includes request ID.
func TestResponseTimeWithRequestID(t *testing.T) {
	// Stack: RequestID -> ResponseTime -> Handler
	handler := RequestID(ResponseTime(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}

// TestResponseTimeStatusCodes verifies different status codes are captured.
func TestResponseTimeStatusCodes(t *testing.T) {
	tests := []struct {
		name       string
		statusCode int
	}{
		{"Success", http.StatusOK},
		{"Created", http.StatusCreated},
		{"Bad Request", http.StatusBadRequest},
		{"Not Found", http.StatusNotFound},
		{"Internal Error", http.StatusInternalServerError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler := ResponseTime(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
			}))

			req := httptest.NewRequest("GET", "/test", nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.statusCode, w.Code)
		})
	}
}

// TestResponseWriterDefaultStatus verifies default 200 status when not explicitly set.
func TestResponseWriterDefaultStatus(t *testing.T) {
	handler := ResponseTime(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Don't call WriteHeader, just write body
		_, _ = w.Write([]byte("test"))
	}))

	req := httptest.NewRequest("GET", "/test", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "test", w.Body.String())
}

// TestMiddlewareStack verifies all middleware work together.
func TestMiddlewareStack(t *testing.T) {
	cfg := &config.Config{
		CORS: config.CORSConfig{
			AllowedOrigins: []string{"http://localhost:3000"},
			AllowedMethods: []string{"GET", "POST"},
			AllowedHeaders: []string{"Content-Type", "X-Request-ID"},
		},
	}

	// Full middleware stack as used in production
	handler := RequestID(
		CORS(cfg)(
			Recoverer(
				ResponseTime(
					http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
						requestID := GetRequestID(r)
						assert.NotEmpty(t, requestID)
						w.WriteHeader(http.StatusOK)
						_, _ = w.Write([]byte("success"))
					}),
				),
			),
		),
	)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	assert.Equal(t, "success", w.Body.String())
	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
	assert.Equal(t, "http://localhost:3000", w.Header().Get("Access-Control-Allow-Origin"))
}
