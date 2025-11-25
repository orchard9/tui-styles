// Package handlers contains HTTP request handlers for the service.
// It includes health checks, API endpoints, and middleware integration.
package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/orchard9/go-core-http-toolkit/cache"
	"github.com/orchard9/go-core-http-toolkit/db"
	"github.com/orchard9/go-core-http-toolkit/response"
	"github.com/orchard9/peach/apps/email-worker/internal/worker"
)

var (
	// Build metadata set at compile time
	Version    = "dev"
	CommitHash = "unknown"
	BuildTime  = "unknown"
)

// Health status constants
const (
	StatusHealthy   = "healthy"
	StatusUnhealthy = "unhealthy"
	StatusDisabled  = "disabled"
	StatusReady     = "ready"
)

var startTime = time.Now()

// HealthHandler handles comprehensive health check endpoints
type HealthHandler struct {
	db     *db.DB
	redis  *cache.RedisClient
	worker *worker.Worker
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(database *db.DB, redisClient *cache.RedisClient, w *worker.Worker) *HealthHandler {
	return &HealthHandler{
		db:     database,
		redis:  redisClient,
		worker: w,
	}
}

// SetBuildMetadata sets the build metadata for health responses
func SetBuildMetadata(version, commitHash, buildTime string) {
	Version = version
	CommitHash = commitHash
	BuildTime = buildTime
}

// SetStartTime sets the server start time
func SetStartTime(start time.Time) {
	startTime = start
}

// HandleHealth returns detailed health status for all components
func (h *HealthHandler) HandleHealth(w http.ResponseWriter, r *http.Request) {
	uptime := time.Since(startTime)
	now := time.Now()

	// Check database health
	var dbHealth map[string]interface{}
	if h.db != nil {
		dbStart := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := h.db.PingContext(ctx)
		dbDuration := time.Since(dbStart)

		if err != nil {
			dbHealth = map[string]interface{}{
				"status":       StatusUnhealthy,
				"message":      fmt.Sprintf("Database connection failed: %v", err),
				"lastCheck":    now.Format(time.RFC3339Nano),
				"responseTime": dbDuration.String(),
				"details": map[string]interface{}{
					"error": err.Error(),
				},
			}
		} else {
			stats := h.db.Stats()
			dbHealth = map[string]interface{}{
				"status":       StatusHealthy,
				"message":      "Database connection is healthy",
				"lastCheck":    now.Format(time.RFC3339Nano),
				"responseTime": dbDuration.String(),
				"details": map[string]interface{}{
					"idle":             stats.Idle,
					"in_use":           stats.InUse,
					"max_open_conns":   stats.MaxOpenConnections,
					"open_connections": stats.OpenConnections,
				},
			}
		}
	} else {
		dbHealth = map[string]interface{}{
			"status":       StatusDisabled,
			"message":      "Database not configured or unavailable",
			"lastCheck":    now.Format(time.RFC3339Nano),
			"responseTime": "0s",
		}
	}
	// Check cache health
	var cacheHealth map[string]interface{}
	if h.redis != nil {
		cacheStart := time.Now()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := h.redis.Set(ctx, "__health_check__", "ok", 1*time.Second)
		cacheDuration := time.Since(cacheStart)

		if err != nil {
			cacheHealth = map[string]interface{}{
				"status":       StatusUnhealthy,
				"message":      fmt.Sprintf("Cache connection failed: %v", err),
				"lastCheck":    now.Format(time.RFC3339Nano),
				"responseTime": cacheDuration.String(),
				"details": map[string]interface{}{
					"error": err.Error(),
				},
			}
		} else {
			cacheHealth = map[string]interface{}{
				"status":       StatusHealthy,
				"message":      "Cache service is healthy",
				"lastCheck":    now.Format(time.RFC3339Nano),
				"responseTime": cacheDuration.String(),
				"details": map[string]interface{}{
					"status": "available",
					"type":   "redis",
				},
			}
		}
	} else {
		cacheHealth = map[string]interface{}{
			"status":       StatusDisabled,
			"message":      "Cache not configured or unavailable",
			"lastCheck":    now.Format(time.RFC3339Nano),
			"responseTime": "0s",
		}
	}

	// Get worker health
	var workerHealth map[string]interface{}
	if h.worker != nil {
		stats := h.worker.GetStats()
		lastTickStr := ""
		if !stats.LastTick.IsZero() {
			lastTickStr = stats.LastTick.Format(time.RFC3339Nano)
		}
		workerHealth = map[string]interface{}{
			"tick_rate":  stats.TickRate,
			"last_tick":  lastTickStr,
			"tick_count": stats.TickCount,
			"running":    stats.Running,
		}
	} else {
		workerHealth = map[string]interface{}{
			"status": StatusDisabled,
		}
	}

	// Determine overall status
	overallStatus := StatusHealthy
	components := map[string]interface{}{
		"database": dbHealth,
		"cache":    cacheHealth,
		"worker":   workerHealth,
	}

	if dbHealth["status"] == StatusUnhealthy || (cacheHealth["status"] == StatusUnhealthy) {
		overallStatus = StatusUnhealthy
	}

	response := map[string]interface{}{
		"status":     overallStatus,
		"service":    "email-worker",
		"version":    Version,
		"commit":     CommitHash,
		"buildTime":  BuildTime,
		"uptime":     uptime.String(),
		"timestamp":  now.Format(time.RFC3339Nano),
		"components": components,
	}

	w.Header().Set("Content-Type", "application/json")
	if overallStatus == StatusUnhealthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	responseBytes, err := json.Marshal(response)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"status": "error", "message": "Failed to marshal response"}`))
		return
	}

	_, _ = w.Write(responseBytes)
}

// Health returns a simple health check handler for backward compatibility
func Health() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = response.JSON(map[string]interface{}{
			"status":  StatusHealthy,
			"service": "email-worker",
			"uptime":  time.Since(startTime).String(),
		}).WriteTo(w)
	}
}

// HandleStartup checks if the application has completed startup initialization
func (h *HealthHandler) HandleStartup(w http.ResponseWriter, r *http.Request) {
	var dbStatus string
	var dbMessage string

	if h.db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		err := h.db.PingContext(ctx)

		if err != nil {
			dbStatus = StatusUnhealthy
			dbMessage = fmt.Sprintf("Database not ready: %v", err)
		} else {
			dbStatus = StatusHealthy
			dbMessage = "Database connection established"
		}
	} else {
		dbStatus = StatusDisabled
		dbMessage = "Database not configured"
	}

	response := map[string]interface{}{
		"status":    dbStatus,
		"message":   dbMessage,
		"timestamp": time.Now().Format(time.RFC3339Nano),
	}

	w.Header().Set("Content-Type", "application/json")
	if dbStatus == StatusUnhealthy {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	responseBytes, _ := json.Marshal(response)
	_, _ = w.Write(responseBytes)
}

// HandleReady checks if the application is ready to serve traffic
func (h *HealthHandler) HandleReady(w http.ResponseWriter, r *http.Request) {
	now := time.Now()
	isReady := true
	readyComponents := make(map[string]string)

	if h.db != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := h.db.PingContext(ctx)

		if err != nil {
			readyComponents["database"] = fmt.Sprintf("not ready: %v", err)
			isReady = false
		} else {
			readyComponents["database"] = StatusReady
		}
	} else {
		readyComponents["database"] = "disabled"
	}

	if h.redis != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		err := h.redis.Set(ctx, "__readiness_check__", "ok", 1*time.Second)

		if err != nil {
			readyComponents["cache"] = fmt.Sprintf("not ready: %v", err)
			isReady = false
		} else {
			readyComponents["cache"] = StatusReady
		}
	} else {
		readyComponents["cache"] = "disabled"
	}

	status := StatusReady
	if !isReady {
		status = "not_ready"
	}

	response := map[string]interface{}{
		"status":     status,
		"timestamp":  now.Format(time.RFC3339Nano),
		"components": readyComponents,
	}

	w.Header().Set("Content-Type", "application/json")
	if !isReady {
		w.WriteHeader(http.StatusServiceUnavailable)
	} else {
		w.WriteHeader(http.StatusOK)
	}

	responseBytes, _ := json.Marshal(response)
	_, _ = w.Write(responseBytes)
}
