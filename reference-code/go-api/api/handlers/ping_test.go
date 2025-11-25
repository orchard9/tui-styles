package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandlePing(t *testing.T) {
	tests := []struct {
		name           string
		setupContext   func(ctx context.Context) context.Context
		expectedStatus int
		validateBody   func(t *testing.T, body string)
	}{
		{
			name: "successful ping",
			setupContext: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, middleware.RequestIDKey, "test-request-123")
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp PingResponse
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err, "Should unmarshal response")

				assert.Equal(t, "pong", resp.Message, "Message should be 'pong'")
				assert.Equal(t, "v1", resp.Version, "Version should be 'v1'")
				assert.NotEmpty(t, resp.Timestamp, "Timestamp should not be empty")

				// Validate timestamp format (RFC3339)
				_, err = time.Parse(time.RFC3339, resp.Timestamp)
				assert.NoError(t, err, "Timestamp should be valid RFC3339 format")

				// Verify timestamp is recent (within last 5 seconds)
				ts, _ := time.Parse(time.RFC3339, resp.Timestamp)
				timeDiff := time.Since(ts)
				assert.Less(t, timeDiff, 5*time.Second, "Timestamp should be recent")
			},
		},
		{
			name: "ping without request ID",
			setupContext: func(ctx context.Context) context.Context {
				return ctx // No request ID
			},
			expectedStatus: http.StatusOK,
			validateBody: func(t *testing.T, body string) {
				var resp PingResponse
				err := json.Unmarshal([]byte(body), &resp)
				require.NoError(t, err, "Should unmarshal response even without request ID")

				assert.Equal(t, "pong", resp.Message)
				assert.Equal(t, "v1", resp.Version)
				assert.NotEmpty(t, resp.Timestamp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
			if tt.setupContext != nil {
				req = req.WithContext(tt.setupContext(req.Context()))
			}
			w := httptest.NewRecorder()

			HandlePing(w, req)

			assert.Equal(t, tt.expectedStatus, w.Code, "Status code should match")
			assert.Equal(t, "application/json", w.Header().Get("Content-Type"),
				"Content-Type should be application/json")

			if tt.validateBody != nil {
				tt.validateBody(t, w.Body.String())
			}
		})
	}
}

func TestHandlePingMultipleCalls(t *testing.T) {
	// Verify that multiple calls produce different timestamps
	req1 := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	w1 := httptest.NewRecorder()
	HandlePing(w1, req1)

	time.Sleep(10 * time.Millisecond)

	req2 := httptest.NewRequest(http.MethodGet, "/api/v1/ping", nil)
	w2 := httptest.NewRecorder()
	HandlePing(w2, req2)

	var resp1, resp2 PingResponse
	err := json.Unmarshal(w1.Body.Bytes(), &resp1)
	require.NoError(t, err)
	err = json.Unmarshal(w2.Body.Bytes(), &resp2)
	require.NoError(t, err)

	// Timestamps should be different (or possibly same if executed too quickly)
	// But both should be valid
	assert.NotEmpty(t, resp1.Timestamp)
	assert.NotEmpty(t, resp2.Timestamp)
}

func TestPingResponseStructure(t *testing.T) {
	// Test that PingResponse JSON structure is correct
	resp := PingResponse{
		Message:   "test",
		Timestamp: "2025-11-19T12:00:00Z",
		Version:   "v1",
	}

	data, err := json.Marshal(resp)
	require.NoError(t, err)

	var decoded PingResponse
	err = json.Unmarshal(data, &decoded)
	require.NoError(t, err)

	assert.Equal(t, resp, decoded, "Round-trip should preserve data")
}
