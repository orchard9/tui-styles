package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Config holds all application configuration
type Config struct {
	Server  ServerConfig
	Logging LoggingConfig
	API     APIConfig
	CORS    CORSConfig
}

// ServerConfig holds server-specific configuration
type ServerConfig struct {
	Port int
	Host string
	Env  string // development, staging, production
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string // debug, info, warn, error
	Format string // json, console
}

// APIConfig holds API configuration
type APIConfig struct {
	Version string
	Prefix  string
}

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string
}

// Load configuration from environment variables (set by envault)
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port: getEnvAsInt("PORT", 34075),
			Host: getEnv("HOST", "localhost"),
			Env:  getEnv("ENV", "development"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "debug"),
			Format: getEnv("LOG_FORMAT", "console"),
		},
		API: APIConfig{
			Version: getEnv("API_VERSION", "v1"),
			Prefix:  getEnv("API_PREFIX", "/api"),
		},
		CORS: CORSConfig{
			AllowedOrigins: getEnvAsSlice("CORS_ALLOWED_ORIGINS", []string{
				"http://localhost:34072",
				"http://localhost:34073",
			}),
			AllowedMethods: getEnvAsSlice("CORS_ALLOWED_METHODS", []string{
				"GET", "POST", "PUT", "DELETE", "OPTIONS",
			}),
			AllowedHeaders: getEnvAsSlice("CORS_ALLOWED_HEADERS", []string{
				"Content-Type", "Authorization",
			}),
		},
	}

	if err := cfg.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return cfg, nil
}

// MustLoad loads config or panics
func MustLoad() *Config {
	cfg, err := Load()
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return cfg
}

// Validate configuration values
func (c *Config) Validate() error {
	// Port validation
	if c.Server.Port < 1024 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid port: %d (must be 1024-65535)", c.Server.Port)
	}

	// Environment validation
	validEnvs := map[string]bool{"development": true, "staging": true, "production": true}
	if !validEnvs[c.Server.Env] {
		return fmt.Errorf("invalid env: %s (must be development, staging, or production)", c.Server.Env)
	}

	// Log level validation
	validLevels := map[string]bool{"debug": true, "info": true, "warn": true, "error": true}
	if !validLevels[c.Logging.Level] {
		return fmt.Errorf("invalid log level: %s", c.Logging.Level)
	}

	return nil
}

// IsDevelopment returns true if running in development environment
func (c *Config) IsDevelopment() bool {
	return c.Server.Env == "development"
}

// IsProduction returns true if running in production environment
func (c *Config) IsProduction() bool {
	return c.Server.Env == "production"
}

// Helper functions for environment variable parsing

// getEnv reads an environment variable or returns the default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt reads an environment variable as an integer or returns the default value
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		panic(fmt.Sprintf("invalid integer value for %s: %s", key, valueStr))
	}
	return value
}

// getEnvAsSlice reads an environment variable as a comma-separated slice or returns the default value
func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	return strings.Split(valueStr, ",")
}
