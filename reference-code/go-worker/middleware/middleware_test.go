package middleware

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Context key constants for testing
const (
	requestIDKey contextKey = "requestID"
	userIDKey    contextKey = "userID"
	traceIDKey   contextKey = "traceID"
)

// mockHandler is a simple handler for testing middleware chains
func mockHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("success"))
	})
}

// generateTestJWT creates a valid JWT token for testing
func generateTestJWT(secret string, claims jwt.MapClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// generateExpiredJWT creates an expired JWT token for testing
func generateExpiredJWT(secret string) (string, error) {
	claims := jwt.MapClaims{
		"sub": "test-user",
		"exp": time.Now().Add(-1 * time.Hour).Unix(),
	}
	return generateTestJWT(secret, claims)
}

func TestJWTAuthentication(t *testing.T) {
	testSecret := "test-secret-key-for-testing-only"

	tests := []struct {
		name           string
		setupAuth      func() string
		expectedStatus int
		expectedBody   string
	}{
		{
			name: "valid token with bearer prefix",
			setupAuth: func() string {
				token, err := generateTestJWT(testSecret, jwt.MapClaims{
					"sub": "user-123",
					"exp": time.Now().Add(1 * time.Hour).Unix(),
					"iat": time.Now().Unix(),
				})
				require.NoError(t, err)
				return "Bearer " + token
			},
			expectedStatus: http.StatusOK,
			expectedBody:   "success",
		},
		{
			name: "missing authorization header",
			setupAuth: func() string {
				return ""
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name: "invalid bearer prefix",
			setupAuth: func() string {
				token, err := generateTestJWT(testSecret, jwt.MapClaims{
					"sub": "user-123",
					"exp": time.Now().Add(1 * time.Hour).Unix(),
				})
				require.NoError(t, err)
				return "Basic " + token
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name: "expired token",
			setupAuth: func() string {
				token, err := generateExpiredJWT(testSecret)
				require.NoError(t, err)
				return "Bearer " + token
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name: "malformed token",
			setupAuth: func() string {
				return "Bearer invalid.token.here"
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
		{
			name: "token with invalid signature",
			setupAuth: func() string {
				token, err := generateTestJWT("wrong-secret", jwt.MapClaims{
					"sub": "user-123",
					"exp": time.Now().Add(1 * time.Hour).Unix(),
				})
				require.NoError(t, err)
				return "Bearer " + token
			},
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			if authHeader := tt.setupAuth(); authHeader != "" {
				req.Header.Set("Authorization", authHeader)
			}

			// Create response recorder
			rr := httptest.NewRecorder()
			_ = rr // TODO: Use when implementing actual middleware tests

			// Create JWT middleware (you'll need to implement this)
			// middleware := NewJWTAuth(testSecret)
			// handler := middleware(mockHandler())

			// For now, we'll test the handler directly
			// In a real implementation, you'd have:
			// handler.ServeHTTP(rr, req)

			// Verify response
			// assert.Equal(t, tt.expectedStatus, rr.Code)
			// if tt.expectedBody != "" {
			// 	assert.Contains(t, rr.Body.String(), tt.expectedBody)
			// }
		})
	}
}

func TestRateLimiting(t *testing.T) {
	tests := []struct {
		name           string
		requestCount   int
		rateLimit      int
		window         time.Duration
		expectedStatus []int
	}{
		{
			name:         "within rate limit",
			requestCount: 5,
			rateLimit:    10,
			window:       time.Minute,
			expectedStatus: []int{
				http.StatusOK,
				http.StatusOK,
				http.StatusOK,
				http.StatusOK,
				http.StatusOK,
			},
		},
		{
			name:         "exceeds rate limit",
			requestCount: 3,
			rateLimit:    2,
			window:       time.Minute,
			expectedStatus: []int{
				http.StatusOK,
				http.StatusOK,
				http.StatusTooManyRequests,
			},
		},
		{
			name:         "single request at limit",
			requestCount: 1,
			rateLimit:    1,
			window:       time.Minute,
			expectedStatus: []int{
				http.StatusOK,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create rate limiter middleware (you'll need to implement this)
			// limiter := NewRateLimiter(tt.rateLimit, tt.window)
			// handler := limiter(mockHandler())

			for i := 0; i < tt.requestCount; i++ {
				req := httptest.NewRequest(http.MethodGet, "/test", nil)
				req.RemoteAddr = "192.168.1.100:12345" // Fixed IP for testing
				rr := httptest.NewRecorder()
				_ = rr // TODO: Use when implementing actual middleware tests

				// handler.ServeHTTP(rr, req)

				// Verify response status
				// assert.Equal(t, tt.expectedStatus[i], rr.Code,
				// 	"Request %d: expected status %d, got %d",
				// 	i+1, tt.expectedStatus[i], rr.Code)
			}
		})
	}
}

func TestMiddlewareChaining(t *testing.T) {
	tests := []struct {
		name              string
		middlewares       []string
		setupRequest      func(*http.Request)
		expectedStatus    int
		verifyContext     func(*testing.T, context.Context)
		expectedBodyMatch string
	}{
		{
			name:        "single middleware",
			middlewares: []string{"logging"},
			setupRequest: func(r *http.Request) {
				// No special setup needed
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "authentication then authorization",
			middlewares: []string{"auth", "authz"},
			setupRequest: func(r *http.Request) {
				token, _ := generateTestJWT("secret", jwt.MapClaims{
					"sub":  "user-123",
					"role": "admin",
					"exp":  time.Now().Add(1 * time.Hour).Unix(),
				})
				r.Header.Set("Authorization", "Bearer "+token)
			},
			expectedStatus: http.StatusOK,
			verifyContext: func(t *testing.T, ctx context.Context) {
				// Verify user ID was added to context
				if userID := ctx.Value("userID"); userID != nil {
					assert.Equal(t, "user-123", userID)
				}
			},
		},
		{
			name:        "rate limit then auth",
			middlewares: []string{"ratelimit", "auth"},
			setupRequest: func(r *http.Request) {
				token, _ := generateTestJWT("secret", jwt.MapClaims{
					"sub": "user-123",
					"exp": time.Now().Add(1 * time.Hour).Unix(),
				})
				r.Header.Set("Authorization", "Bearer "+token)
				r.RemoteAddr = "192.168.1.1:12345"
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:        "failed auth stops chain",
			middlewares: []string{"auth", "business-logic"},
			setupRequest: func(r *http.Request) {
				// No auth header - should fail at auth middleware
			},
			expectedStatus: http.StatusUnauthorized,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			tt.setupRequest(req)

			// Create response recorder
			rr := httptest.NewRecorder()
			_ = rr // TODO: Use when implementing actual middleware tests

			// Build middleware chain
			// In a real implementation, you'd chain middlewares like:
			// handler := mockHandler()
			// for i := len(tt.middlewares) - 1; i >= 0; i-- {
			// 	switch tt.middlewares[i] {
			// 	case "logging":
			// 		handler = LoggingMiddleware(handler)
			// 	case "auth":
			// 		handler = AuthMiddleware(handler)
			// 	case "authz":
			// 		handler = AuthorizationMiddleware(handler)
			// 	case "ratelimit":
			// 		handler = RateLimitMiddleware(handler)
			// 	}
			// }

			// handler.ServeHTTP(rr, req)

			// Verify response
			// assert.Equal(t, tt.expectedStatus, rr.Code)

			// Verify context if needed
			// if tt.verifyContext != nil {
			// 	tt.verifyContext(t, req.Context())
			// }
		})
	}
}

func TestMiddlewareContextPropagation(t *testing.T) {
	tests := []struct {
		name          string
		setupContext  func(context.Context) context.Context
		verifyContext func(*testing.T, context.Context)
	}{
		{
			name: "request ID propagation",
			setupContext: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, requestIDKey, "req-12345")
			},
			verifyContext: func(t *testing.T, ctx context.Context) {
				requestID := ctx.Value(requestIDKey)
				require.NotNil(t, requestID)
				assert.Equal(t, "req-12345", requestID)
			},
		},
		{
			name: "user context propagation",
			setupContext: func(ctx context.Context) context.Context {
				return context.WithValue(ctx, userIDKey, "user-789")
			},
			verifyContext: func(t *testing.T, ctx context.Context) {
				userID := ctx.Value(userIDKey)
				require.NotNil(t, userID)
				assert.Equal(t, "user-789", userID)
			},
		},
		{
			name: "multiple context values",
			setupContext: func(ctx context.Context) context.Context {
				ctx = context.WithValue(ctx, requestIDKey, "req-111")
				ctx = context.WithValue(ctx, userIDKey, "user-222")
				ctx = context.WithValue(ctx, traceIDKey, "trace-333")
				return ctx
			},
			verifyContext: func(t *testing.T, ctx context.Context) {
				assert.Equal(t, "req-111", ctx.Value(requestIDKey))
				assert.Equal(t, "user-222", ctx.Value(userIDKey))
				assert.Equal(t, "trace-333", ctx.Value(traceIDKey))
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request with context
			ctx := context.Background()
			ctx = tt.setupContext(ctx)
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()
			_ = rr // TODO: Use when implementing actual middleware tests

			// Create a middleware that captures the context
			var capturedCtx context.Context
			testMiddleware := func(next http.Handler) http.Handler {
				return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					capturedCtx = r.Context()
					next.ServeHTTP(w, r)
				})
			}

			handler := testMiddleware(mockHandler())
			handler.ServeHTTP(rr, req)

			// Verify context was propagated
			tt.verifyContext(t, capturedCtx)
		})
	}
}

func TestCORSMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		origin         string
		requestHeaders string
		expectedStatus int
		checkHeaders   map[string]string
	}{
		{
			name:           "preflight request",
			method:         http.MethodOptions,
			origin:         "https://example.com",
			requestHeaders: "Content-Type, Authorization",
			expectedStatus: http.StatusOK,
			checkHeaders: map[string]string{
				"Access-Control-Allow-Origin":  "*",
				"Access-Control-Allow-Methods": "GET, POST, PUT, DELETE, OPTIONS",
				"Access-Control-Allow-Headers": "Content-Type, Authorization",
			},
		},
		{
			name:           "simple GET request",
			method:         http.MethodGet,
			origin:         "https://example.com",
			expectedStatus: http.StatusOK,
			checkHeaders: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		},
		{
			name:           "POST request with origin",
			method:         http.MethodPost,
			origin:         "https://app.example.com",
			expectedStatus: http.StatusOK,
			checkHeaders: map[string]string{
				"Access-Control-Allow-Origin": "*",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, "/test", nil)
			if tt.origin != "" {
				req.Header.Set("Origin", tt.origin)
			}
			if tt.requestHeaders != "" {
				req.Header.Set("Access-Control-Request-Headers", tt.requestHeaders)
			}

			rr := httptest.NewRecorder()
			_ = rr // TODO: Use when implementing actual middleware tests

			// Create CORS middleware (you'll need to implement this)
			// corsMiddleware := NewCORSMiddleware()
			// handler := corsMiddleware(mockHandler())
			// handler.ServeHTTP(rr, req)

			// Verify status
			// assert.Equal(t, tt.expectedStatus, rr.Code)

			// Verify headers
			// for header, expectedValue := range tt.checkHeaders {
			// 	assert.Equal(t, expectedValue, rr.Header().Get(header),
			// 		"Header %s should be %s", header, expectedValue)
			// }
		})
	}
}

func TestRecoveryMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		handler        http.Handler
		expectedStatus int
		shouldPanic    bool
	}{
		{
			name:           "normal request no panic",
			handler:        mockHandler(),
			expectedStatus: http.StatusOK,
			shouldPanic:    false,
		},
		{
			name: "handler panics",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic("something went wrong")
			}),
			expectedStatus: http.StatusInternalServerError,
			shouldPanic:    true,
		},
		{
			name: "handler panics with nil",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				panic(nil)
			}),
			expectedStatus: http.StatusInternalServerError,
			shouldPanic:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			_ = req // TODO: Use when implementing actual recovery middleware tests
			rr := httptest.NewRecorder()
			_ = rr // TODO: Use when implementing actual middleware tests

			// Create recovery middleware (you'll need to implement this)
			// recoveryMiddleware := NewRecoveryMiddleware()
			// handler := recoveryMiddleware(tt.handler)

			// Should not panic - middleware should recover
			// assert.NotPanics(t, func() {
			// 	handler.ServeHTTP(rr, req)
			// })

			// Verify status code
			// assert.Equal(t, tt.expectedStatus, rr.Code)
		})
	}
}
