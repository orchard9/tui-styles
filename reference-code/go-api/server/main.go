// Package main is the entry point for the Creator API HTTP server.
// It initializes the server with configuration, sets up routes, and handles graceful shutdown.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/masquerade/creator-api/internal/api"
	"github.com/masquerade/creator-api/internal/config"
	"github.com/rs/zerolog/log"
)

// Version is the current version of the Creator API.
// This will be replaced with build-time injection in the future.
const Version = "0.1.0"

func main() {
	// Load configuration from environment variables (set by envault)
	cfg, err := config.Load()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed to load config")
	}

	// Log startup information
	log.Info().
		Str("version", Version).
		Str("env", cfg.Server.Env).
		Int("port", cfg.Server.Port).
		Str("host", cfg.Server.Host).
		Msg("Starting Creator API server")

	// Create router with all API routes
	router := api.NewRouter()

	// Create HTTP server with timeouts
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:         addr,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in goroutine so it doesn't block shutdown handling
	go func() {
		log.Info().Msgf("Server listening on %s", addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal().Err(err).Msg("Server failed to start")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	// Capture SIGINT (Ctrl+C) and SIGTERM (container orchestrator)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Info().Msg("Server shutting down...")

	// Graceful shutdown with 30 second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal().Err(err).Msg("Server forced to shutdown")
	}

	log.Info().Msg("Server exited")
}
