// Package middleware contains custom HTTP middleware (adapters) for cross-cutting concerns.
//
// Middleware are ADAPTERS in hexagonal architecture:
// - Wrap HTTP handlers with additional functionality
// - Handle infrastructure concerns (metrics, logging, auth)
// - Remain independent of business logic
// - Compose cleanly with standard http.Handler interface
//
// Responsibilities:
// - Request/Response metrics collection
// - Audit logging for sensitive operations
// - Role-based access control enforcement
// - Query parameter parsing and validation
//
// NOT Responsible For:
// - Business logic (belongs in domain/use cases)
// - Data persistence (belongs in repositories)
// - Domain validation (belongs in domain layer)
package middleware

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

	"github.com/orchard9/peach/apps/email-worker/internal/domain"
)

// Context keys for middleware values
type contextKey string

const (
	// ContextKeyPagination stores parsed pagination parameters
	ContextKeyPagination contextKey = "pagination"

	// ContextKeyUserID stores authenticated user ID
	ContextKeyUserID contextKey = "user_id"

	// ContextKeyUserRoles stores authenticated user roles
	ContextKeyUserRoles contextKey = "user_roles"

	// ContextKeyRequestID stores request ID for tracing
	ContextKeyRequestID contextKey = "request_id"
)

// ============================================================================
// METRICS MIDDLEWARE
// ============================================================================

// MetricsConfig configures metrics collection
type MetricsConfig struct {
	// ServiceName is used as namespace for metrics
	ServiceName string

	// RecordDuration enables request duration tracking
	RecordDuration bool

	// RecordSize enables request/response size tracking
	RecordSize bool

	// ExcludePaths are paths to exclude from metrics (e.g., /health, /metrics)
	ExcludePaths []string
}

// DefaultMetricsConfig returns sensible defaults
func DefaultMetricsConfig(serviceName string) MetricsConfig {
	return MetricsConfig{
		ServiceName:    serviceName,
		RecordDuration: true,
		RecordSize:     true,
		ExcludePaths:   []string{"/health", "/metrics", "/ready"},
	}
}

// MetricsMiddleware collects HTTP request/response metrics
//
// Metrics Collected:
// - http_requests_total: Counter of total requests by method, path, status
// - http_request_duration_seconds: Histogram of request duration
// - http_request_size_bytes: Histogram of request body size
// - http_response_size_bytes: Histogram of response body size
//
// Thread-Safety: Safe for concurrent use (Prometheus handles synchronization)
type MetricsMiddleware struct {
	config MetricsConfig

	requestsTotal   *prometheus.CounterVec
	requestDuration *prometheus.HistogramVec
	requestSize     *prometheus.HistogramVec
	responseSize    *prometheus.HistogramVec
}

// NewMetricsMiddleware creates a new metrics middleware
//
// Automatically registers metrics with Prometheus default registry.
func NewMetricsMiddleware(config MetricsConfig) *MetricsMiddleware {
	m := &MetricsMiddleware{
		config: config,
		requestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: config.ServiceName,
				Name:      "http_requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
	}

	if config.RecordDuration {
		m.requestDuration = promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: config.ServiceName,
				Name:      "http_request_duration_seconds",
				Help:      "HTTP request duration in seconds",
				Buckets:   []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10},
			},
			[]string{"method", "path"},
		)
	}

	if config.RecordSize {
		m.requestSize = promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: config.ServiceName,
				Name:      "http_request_size_bytes",
				Help:      "HTTP request size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 7), // 100B to 100MB
			},
			[]string{"method", "path"},
		)
		m.responseSize = promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: config.ServiceName,
				Name:      "http_response_size_bytes",
				Help:      "HTTP response size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 7),
			},
			[]string{"method", "path"},
		)
	}

	return m
}

// Handler wraps an http.Handler with metrics collection
func (m *MetricsMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip excluded paths
		if m.shouldExclude(r.URL.Path) {
			next.ServeHTTP(w, r)
			return
		}

		start := time.Now()

		// Wrap response writer to capture status and size
		rw := &metricsResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK, // Default if WriteHeader not called
		}

		// Execute handler
		next.ServeHTTP(rw, r)

		// Record metrics
		duration := time.Since(start).Seconds()
		method := r.Method
		path := r.URL.Path
		status := strconv.Itoa(rw.statusCode)

		m.requestsTotal.WithLabelValues(method, path, status).Inc()

		if m.config.RecordDuration {
			m.requestDuration.WithLabelValues(method, path).Observe(duration)
		}

		if m.config.RecordSize {
			if r.ContentLength > 0 {
				m.requestSize.WithLabelValues(method, path).Observe(float64(r.ContentLength))
			}
			if rw.bytesWritten > 0 {
				m.responseSize.WithLabelValues(method, path).Observe(float64(rw.bytesWritten))
			}
		}
	})
}

func (m *MetricsMiddleware) shouldExclude(path string) bool {
	for _, excluded := range m.config.ExcludePaths {
		if path == excluded || strings.HasPrefix(path, excluded) {
			return true
		}
	}
	return false
}

// metricsResponseWriter wraps http.ResponseWriter to capture metrics
type metricsResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	bytesWritten int
}

func (rw *metricsResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *metricsResponseWriter) Write(b []byte) (int, error) {
	n, err := rw.ResponseWriter.Write(b)
	rw.bytesWritten += n
	return n, err
}

// ============================================================================
// AUDIT MIDDLEWARE
// ============================================================================

// AuditLogger defines the interface for audit logging
//
// Implementations should persist audit logs to a durable store
// (database, message queue, external audit service)
type AuditLogger interface {
	// LogAuditEvent records a security-relevant event
	LogAuditEvent(ctx context.Context, event AuditEvent) error
}

// AuditEvent represents a security-relevant action
type AuditEvent struct {
	Timestamp   time.Time              `json:"timestamp"`
	UserID      string                 `json:"user_id"`
	Action      string                 `json:"action"`       // e.g., "user.delete", "user.update_email"
	Resource    string                 `json:"resource"`     // e.g., "user:uuid"
	Method      string                 `json:"method"`       // HTTP method
	Path        string                 `json:"path"`         // Request path
	StatusCode  int                    `json:"status_code"`  // Response status
	IPAddress   string                 `json:"ip_address"`   // Client IP
	UserAgent   string                 `json:"user_agent"`   // Client user agent
	RequestBody map[string]interface{} `json:"request_body"` // Sanitized request body
	Metadata    map[string]interface{} `json:"metadata"`     // Additional context
}

// AuditConfig configures audit logging behavior
type AuditConfig struct {
	// SensitiveActions are actions that require audit logging
	// Examples: "DELETE", "PATCH", "user.delete", "user.update_email"
	SensitiveActions []string

	// SanitizeFields are field names to redact from request bodies
	SanitizeFields []string

	// IncludeRequestBody enables request body logging (sanitized)
	IncludeRequestBody bool
}

// DefaultAuditConfig returns sensible defaults
func DefaultAuditConfig() AuditConfig {
	return AuditConfig{
		SensitiveActions: []string{
			"DELETE",
			"PATCH",
			"PUT",
		},
		SanitizeFields: []string{
			"password",
			"token",
			"secret",
			"api_key",
			"apikey",
			"authorization",
		},
		IncludeRequestBody: true,
	}
}

// AuditMiddleware logs security-relevant operations
//
// Use Cases:
// - Compliance requirements (GDPR, HIPAA, SOC2)
// - Security incident investigation
// - User activity tracking
// - Detecting suspicious patterns
//
// Thread-Safety: Safe if AuditLogger is thread-safe
type AuditMiddleware struct {
	logger AuditLogger
	config AuditConfig
}

// NewAuditMiddleware creates a new audit middleware
func NewAuditMiddleware(logger AuditLogger, config AuditConfig) *AuditMiddleware {
	return &AuditMiddleware{
		logger: logger,
		config: config,
	}
}

// Handler wraps an http.Handler with audit logging
func (m *AuditMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if this action requires auditing
		if !m.shouldAudit(r) {
			next.ServeHTTP(w, r)
			return
		}

		// Capture response status
		rw := &auditResponseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		// Execute handler
		next.ServeHTTP(rw, r)

		// Log audit event asynchronously (don't block response)
		go func() {
			event := m.buildAuditEvent(r, rw.statusCode)
			if err := m.logger.LogAuditEvent(context.Background(), event); err != nil {
				// Log as warn since audit logging failures are concerning but not critical
				// The request has already succeeded, we just failed to log it
				slog.Warn("audit event logging failed",
					slog.String("error", err.Error()),
					slog.String("action", event.Action),
					slog.String("user_id", event.UserID),
				)
			}
		}()
	})
}

func (m *AuditMiddleware) shouldAudit(r *http.Request) bool {
	for _, action := range m.config.SensitiveActions {
		if r.Method == action {
			return true
		}
		// Could also check path patterns here
		if strings.Contains(r.URL.Path, action) {
			return true
		}
	}
	return false
}

func (m *AuditMiddleware) buildAuditEvent(r *http.Request, statusCode int) AuditEvent {
	event := AuditEvent{
		Timestamp:  time.Now(),
		Action:     fmt.Sprintf("%s %s", r.Method, r.URL.Path),
		Resource:   extractResourceFromPath(r.URL.Path),
		Method:     r.Method,
		Path:       r.URL.Path,
		StatusCode: statusCode,
		IPAddress:  extractIPAddress(r),
		UserAgent:  r.UserAgent(),
		Metadata:   make(map[string]interface{}),
	}

	// Extract user ID from context (set by auth middleware)
	if userID, ok := r.Context().Value(ContextKeyUserID).(string); ok {
		event.UserID = userID
	}

	// Capture and sanitize request body if configured
	if m.config.IncludeRequestBody && r.Body != nil {
		event.RequestBody = m.sanitizeRequestBody(r)
	}

	return event
}

func (m *AuditMiddleware) sanitizeRequestBody(r *http.Request) map[string]interface{} {
	// Try to parse JSON body
	var body map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil
	}

	// Sanitize sensitive fields
	for _, field := range m.config.SanitizeFields {
		if _, exists := body[field]; exists {
			body[field] = "[REDACTED]"
		}
	}

	return body
}

type auditResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *auditResponseWriter) WriteHeader(statusCode int) {
	rw.statusCode = statusCode
	rw.ResponseWriter.WriteHeader(statusCode)
}

// extractResourceFromPath extracts resource identifier from URL path
// Example: "/users/123" -> "user:123"
func extractResourceFromPath(path string) string {
	parts := strings.Split(strings.Trim(path, "/"), "/")
	if len(parts) >= 2 {
		return fmt.Sprintf("%s:%s", strings.TrimSuffix(parts[0], "s"), parts[1])
	}
	return path
}

// extractIPAddress extracts client IP from request
func extractIPAddress(r *http.Request) string {
	// Check X-Forwarded-For header (set by proxies)
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	// Check X-Real-IP header
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		return xri
	}

	// Fall back to RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}

// ============================================================================
// ROLE-BASED ACCESS CONTROL (RBAC) MIDDLEWARE
// ============================================================================

// Role represents a user role
type Role string

const (
	RoleAdmin     Role = "admin"
	RoleUser      Role = "user"
	RoleModerator Role = "moderator"
	RoleGuest     Role = "guest"
)

// RBACMiddleware enforces role-based access control
//
// Responsibilities:
// - Check if authenticated user has required roles
// - Allow/deny access based on role configuration
// - Return 403 Forbidden for unauthorized access
//
// # Assumes auth middleware has already set ContextKeyUserRoles
//
// Thread-Safety: Safe for concurrent use (stateless)
type RBACMiddleware struct {
	requiredRoles []Role
	requireAll    bool // If true, user must have ALL roles; if false, ANY role
}

// NewRBACMiddleware creates RBAC middleware requiring ANY of the roles
func NewRBACMiddleware(roles ...Role) *RBACMiddleware {
	return &RBACMiddleware{
		requiredRoles: roles,
		requireAll:    false,
	}
}

// NewRBACMiddlewareAll creates RBAC middleware requiring ALL roles
func NewRBACMiddlewareAll(roles ...Role) *RBACMiddleware {
	return &RBACMiddleware{
		requiredRoles: roles,
		requireAll:    true,
	}
}

// Handler wraps an http.Handler with RBAC enforcement
func (m *RBACMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract user roles from context
		userRoles, ok := r.Context().Value(ContextKeyUserRoles).([]Role)
		if !ok {
			respondError(w, http.StatusUnauthorized, "Authentication required")
			return
		}

		// Check authorization
		if !m.isAuthorized(userRoles) {
			respondError(w, http.StatusForbidden, "Insufficient permissions")
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (m *RBACMiddleware) isAuthorized(userRoles []Role) bool {
	if len(m.requiredRoles) == 0 {
		return true // No roles required
	}

	userRoleSet := make(map[Role]bool)
	for _, role := range userRoles {
		userRoleSet[role] = true
	}

	if m.requireAll {
		// Must have ALL required roles
		for _, required := range m.requiredRoles {
			if !userRoleSet[required] {
				return false
			}
		}
		return true
	}

	// Must have ANY required role
	for _, required := range m.requiredRoles {
		if userRoleSet[required] {
			return true
		}
	}
	return false
}

// RequireAdmin is a convenience function for admin-only routes
func RequireAdmin() *RBACMiddleware {
	return NewRBACMiddleware(RoleAdmin)
}

// RequireAuthenticated is a convenience function requiring any authenticated user
func RequireAuthenticated() *RBACMiddleware {
	return NewRBACMiddleware(RoleAdmin, RoleUser, RoleModerator)
}

// ============================================================================
// PAGINATION MIDDLEWARE
// ============================================================================

// PaginationParams holds parsed pagination parameters
type PaginationParams struct {
	Limit  int
	Offset int
	Cursor string
	Type   domain.PaginationType
}

// PaginationConfig configures pagination parsing
type PaginationConfig struct {
	DefaultLimit int
	MaxLimit     int
	DefaultType  domain.PaginationType
}

// DefaultPaginationConfig returns sensible defaults
func DefaultPaginationConfig() PaginationConfig {
	return PaginationConfig{
		DefaultLimit: 20,
		MaxLimit:     100,
		DefaultType:  domain.PaginationTypeOffset,
	}
}

// PaginationMiddleware parses and validates pagination query parameters
//
// Supported Query Parameters:
// - limit: Maximum results (default: 20, max: 100)
// - offset: Number of results to skip (offset pagination)
// - cursor: Opaque cursor token (cursor pagination)
// - type: Pagination type ("offset" or "cursor")
//
// # Parsed parameters are stored in request context under ContextKeyPagination
//
// Thread-Safety: Safe for concurrent use (stateless)
type PaginationMiddleware struct {
	config PaginationConfig
}

// NewPaginationMiddleware creates a new pagination middleware
func NewPaginationMiddleware(config PaginationConfig) *PaginationMiddleware {
	return &PaginationMiddleware{config: config}
}

// Handler wraps an http.Handler with pagination parsing
func (m *PaginationMiddleware) Handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		params, err := m.parsePaginationParams(r)
		if err != nil {
			respondError(w, http.StatusBadRequest, err.Error())
			return
		}

		// Store in context
		ctx := context.WithValue(r.Context(), ContextKeyPagination, params)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (m *PaginationMiddleware) parsePaginationParams(r *http.Request) (*PaginationParams, error) {
	params := &PaginationParams{
		Limit: m.config.DefaultLimit,
		Type:  m.config.DefaultType,
	}

	// Parse limit
	if limitStr := r.URL.Query().Get("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err != nil {
			return nil, fmt.Errorf("invalid limit: %w", err)
		}
		if limit <= 0 {
			return nil, fmt.Errorf("limit must be positive")
		}
		if limit > m.config.MaxLimit {
			limit = m.config.MaxLimit
		}
		params.Limit = limit
	}

	// Parse offset (offset pagination)
	if offsetStr := r.URL.Query().Get("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err != nil {
			return nil, fmt.Errorf("invalid offset: %w", err)
		}
		if offset < 0 {
			return nil, fmt.Errorf("offset cannot be negative")
		}
		params.Offset = offset
		params.Type = domain.PaginationTypeOffset
	}

	// Parse cursor (cursor pagination)
	if cursor := r.URL.Query().Get("cursor"); cursor != "" {
		// Validate cursor format
		if _, err := domain.DecodeCursor(cursor); err != nil {
			return nil, fmt.Errorf("invalid cursor: %w", err)
		}
		params.Cursor = cursor
		params.Type = domain.PaginationTypeCursor
	}

	// Explicit type parameter
	if typeStr := r.URL.Query().Get("type"); typeStr != "" {
		paginationType := domain.PaginationType(typeStr)
		if paginationType != domain.PaginationTypeOffset && paginationType != domain.PaginationTypeCursor {
			return nil, fmt.Errorf("invalid pagination type: must be 'offset' or 'cursor'")
		}
		params.Type = paginationType
	}

	return params, nil
}

// GetPaginationFromContext extracts pagination params from context
//
// Returns nil if pagination middleware was not applied
func GetPaginationFromContext(ctx context.Context) *PaginationParams {
	params, _ := ctx.Value(ContextKeyPagination).(*PaginationParams)
	return params
}

// ============================================================================
// HELPER FUNCTIONS
// ============================================================================

// respondError writes JSON error response
func respondError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
