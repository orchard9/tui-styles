package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/orchard9/go-core-http-toolkit/cache"
	"github.com/orchard9/go-core-http-toolkit/db"
	"github.com/orchard9/go-core-http-toolkit/middleware"
	"github.com/orchard9/go-core-http-toolkit/observability"
	"github.com/orchard9/peach/apps/email-worker/internal/config"
	"github.com/orchard9/peach/apps/email-worker/internal/handlers"
	"github.com/orchard9/peach/apps/email-worker/internal/worker"
	"github.com/redis/go-redis/v9"
)

func main() {
	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize observability FIRST (sets up slog.Default properly)
	shutdown, err := observability.Init(cfg.App.Name)
	if err != nil {
		log.Fatalf("Failed to init observability: %v", err)
	}
	defer func() {
		_ = shutdown(context.Background())
	}()

	// Connect to PostgreSQL (with graceful fallback) - now safe to use slog
	database, err := db.NewPostgresDB(db.Config{
		DSN:             cfg.Database.URL,
		MaxOpenConns:    cfg.Database.MaxOpenConns,
		MaxIdleConns:    cfg.Database.MaxIdleConns,
		ConnMaxLifetime: cfg.Database.ConnMaxLifetime,
	})
	if err != nil {
		slog.Warn("Failed to connect to database, running without database", slog.String("error", err.Error()))
		database = nil
	} else {
		slog.Info("Connected to database successfully")
	}
	if database != nil {
		defer func() {
			_ = database.Close()
		}()
	}
	// Connect to Redis (with graceful fallback)
	var redisClient *cache.RedisClient
	if cfg.Redis.URL != "" {
		redisOpts, err := redis.ParseURL(cfg.Redis.URL)
		if err != nil {
			slog.Warn("Failed to parse Redis URL, running without cache", slog.String("error", err.Error()))
		} else {
			redisClient, err = cache.NewRedisClient(cache.Config{
				Address:  redisOpts.Addr,
				Password: redisOpts.Password,
				DB:       redisOpts.DB,
				PoolSize: cfg.Redis.PoolSize,
			})
			if err != nil {
				slog.Warn("Failed to connect to Redis, running without cache", slog.String("error", err.Error()))
				redisClient = nil
			} else {
				slog.Info("Connected to Redis successfully")
			}
		}
	}
	if redisClient != nil {
		defer func() {
			_ = redisClient.Close()
		}()
	}

	// Initialize worker
	emailWorker := worker.New(cfg.Worker.TickRate)
	slog.Info("Worker initialized", "tick_rate", cfg.Worker.TickRate)
	// Setup HTTP router
	r := chi.NewRouter()

	// Middleware
	r.Use(chimiddleware.Recoverer)
	r.Use(middleware.RealIP())
	r.Use(middleware.RequestID())
	r.Use(middleware.Logger(middleware.LoggingConfig{
		Logger: slog.Default(),
	}))
	r.Use(middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// Routes
	// Create enhanced health handler (handles nil dependencies gracefully)
	healthHandler := handlers.NewHealthHandler(database, redisClient, emailWorker)
	r.Get("/health", healthHandler.HandleHealth)
	r.Get("/health/startup", healthHandler.HandleStartup)
	r.Get("/health/ready", healthHandler.HandleReady)
	// Public API routes with /api/v1 prefix
	r.Route("/api/v1", func(r chi.Router) {
		// Add public endpoints here
	})

	// Setup graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		slog.Info("Shutdown signal received")
		cancel()
	}()
	// Start HTTP server
	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		slog.Info("HTTP server starting", "port", cfg.Server.Port)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP server failed: %v", err)
		}
	}()

	// Start worker in background
	go emailWorker.Start(ctx)

	// Wait for shutdown
	<-ctx.Done()

	slog.Info("Shutting down services...")

	// Shutdown HTTP server
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown failed", "error", err)
	}

	// Worker will stop when ctx is cancelled (already done above)
	// Give it a moment to finish current tick
	time.Sleep(100 * time.Millisecond)

	if database != nil {
		_ = database.Close()
	}
	if redisClient != nil {
		_ = redisClient.Close()
	}

	slog.Info("Server shutdown complete")
}
