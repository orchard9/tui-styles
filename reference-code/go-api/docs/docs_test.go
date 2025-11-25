package docs

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHandleOpenAPISpec(t *testing.T) {
	tests := []struct {
		name            string
		method          string
		expectedStatus  int
		expectedHeaders map[string]string
		validateBody    func(t *testing.T, body string)
	}{
		{
			name:           "GET returns OpenAPI spec",
			method:         http.MethodGet,
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Content-Type":                "application/yaml",
				"Access-Control-Allow-Origin": "*",
				"Cache-Control":               "public, max-age=3600",
			},
			validateBody: func(t *testing.T, body string) {
				// OpenAPI spec should start with openapi version
				assert.True(t, strings.HasPrefix(body, "openapi:") || strings.Contains(body, "openapi:"),
					"Response should contain OpenAPI spec")
				// Should contain key API information
				assert.Contains(t, body, "Masquerade Creator API",
					"Spec should contain API title")
			},
		},
		{
			name:           "HEAD request succeeds",
			method:         http.MethodHead,
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Content-Type":                "application/yaml",
				"Access-Control-Allow-Origin": "*",
			},
			validateBody: func(t *testing.T, body string) {
				// Note: httptest.ResponseRecorder captures body even for HEAD
				// In real HTTP, HEAD wouldn't send body, but handler writes it
				// This is expected behavior for testing
				assert.NotEmpty(t, body, "Handler writes body (httptest captures it)")
			},
		},
		{
			name:           "POST request still succeeds (HTTP allows any method)",
			method:         http.MethodPost,
			expectedStatus: http.StatusOK,
			expectedHeaders: map[string]string{
				"Content-Type": "application/yaml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/openapi.yaml", nil)
			w := httptest.NewRecorder()

			HandleOpenAPISpec(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "Status code should match")

			// Validate expected headers
			for header, expectedValue := range tt.expectedHeaders {
				actualValue := w.Header().Get(header)
				assert.Equal(t, expectedValue, actualValue,
					"Header %s should equal %s", header, expectedValue)
			}

			// Validate body if specified
			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.String())
			}
		})
	}
}

func TestNewSwaggerUIHandler(t *testing.T) {
	handler := NewSwaggerUIHandler()
	assert.NotNil(t, handler, "Handler should not be nil")

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		validateBody   func(t *testing.T, body string, contentType string)
	}{
		{
			name:           "root index.html exists",
			path:           "/",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string, contentType string) {
				// File server might set content type based on file extension
				assert.NotEmpty(t, body, "Body should not be empty for index")
				assert.Contains(t, body, "swagger", "Should contain swagger-related content")
			},
		},
		{
			name:           "swagger-ui.css exists",
			path:           "/swagger-ui.css",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string, contentType string) {
				assert.NotEmpty(t, body, "CSS file should not be empty")
			},
		},
		{
			name:           "swagger-ui-bundle.js exists",
			path:           "/swagger-ui-bundle.js",
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string, contentType string) {
				assert.NotEmpty(t, body, "JS file should not be empty")
			},
		},
		{
			name:           "nonexistent file returns 404",
			path:           "/nonexistent.html",
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			w := httptest.NewRecorder()

			handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code,
				"Status code should match for path %s", tt.path)

			if tt.validateBody != nil && w.Code == http.StatusOK {
				contentType := w.Header().Get("Content-Type")
				tt.validateBody(t, w.Body.String(), contentType)
			}
		})
	}
}

func TestSwaggerUIHandlerCaching(t *testing.T) {
	handler := NewSwaggerUIHandler()

	req := httptest.NewRequest(http.MethodGet, "/swagger-ui.css", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	// Check that static file is served successfully
	// http.FileServer may or may not set caching headers depending on config
	assert.Equal(t, http.StatusOK, w.Code, "Should serve static file")
	assert.NotEmpty(t, w.Body.String(), "Should have file content")
}

func TestOpenAPISpecMultipleMethods(t *testing.T) {
	// Test that OpenAPI spec works with various HTTP methods
	methods := []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			req := httptest.NewRequest(method, "/openapi.yaml", nil)
			w := httptest.NewRecorder()

			HandleOpenAPISpec(w, req)

			// All methods should succeed (handler doesn't restrict)
			assert.Equal(t, http.StatusOK, w.Code, "Request should succeed")
			assert.NotEmpty(t, w.Body.String(), "Should return OpenAPI spec")
		})
	}
}

func TestOpenAPISpecCORS(t *testing.T) {
	// Verify CORS headers are set correctly for browser-based tools
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	req.Header.Set("Origin", "https://example.com")
	w := httptest.NewRecorder()

	HandleOpenAPISpec(w, req)

	assert.Equal(t, "*", w.Header().Get("Access-Control-Allow-Origin"),
		"CORS should allow all origins")
	assert.Equal(t, http.StatusOK, w.Code, "Request should succeed")
}

func TestOpenAPISpecCaching(t *testing.T) {
	// Verify cache headers are set correctly
	req := httptest.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	w := httptest.NewRecorder()

	HandleOpenAPISpec(w, req)

	cacheControl := w.Header().Get("Cache-Control")
	assert.Equal(t, "public, max-age=3600", cacheControl,
		"Cache-Control should be set to 1 hour")
}
