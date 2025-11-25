// Package config provides configuration management for the service.
// It handles loading configuration from environment variables and provides
// validation for all required settings.
package config

import (
	"context"
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

// Config holds all application configuration
type Config struct {
	App           AppConfig           `json:"app"`
	Server        ServerConfig        `json:"server" validate:"required"`
	Database      DatabaseConfig      `json:"database,omitempty"`
	Redis         RedisConfig         `json:"redis,omitempty"`
	JWT           JWTConfig           `json:"jwt,omitempty"`
	Observability ObservabilityConfig `json:"observability,omitempty"`
	Worker        WorkerConfig        `json:"worker"`
}

// AppConfig holds application-specific configuration
type AppConfig struct {
	Name        string `json:"name" validate:"required"`
	Environment string `json:"environment" validate:"required,oneof=development staging production"`
	Debug       bool   `json:"debug"`
}

// ServerConfig holds server configuration
type ServerConfig struct {
	Host         string        `json:"host" validate:"required"`
	Port         int           `json:"port" validate:"required,min=1,max=65535"`
	ReadTimeout  time.Duration `json:"read_timeout"`
	WriteTimeout time.Duration `json:"write_timeout"`
	IdleTimeout  time.Duration `json:"idle_timeout"`
	GRPCPort     int           `json:"grpc_port,omitempty"`
}

// DatabaseConfig holds database configuration
type DatabaseConfig struct {
	URL             string        `json:"url" validate:"required"`
	MaxOpenConns    int           `json:"max_open_conns" validate:"min=1"`
	MaxIdleConns    int           `json:"max_idle_conns" validate:"min=0"`
	ConnMaxLifetime time.Duration `json:"conn_max_lifetime"`
}

// RedisConfig holds Redis configuration
type RedisConfig struct {
	URL      string `json:"url" validate:"required"`
	PoolSize int    `json:"pool_size" validate:"min=1"`
}

// JWTConfig holds JWT configuration
type JWTConfig struct {
	Secret     string        `json:"secret" validate:"required,min=32"`
	Issuer     string        `json:"issuer" validate:"required"`
	Expiration time.Duration `json:"expiration" validate:"required,min=1m"`
}

// ObservabilityConfig holds observability configuration
type ObservabilityConfig struct {
	ServiceName    string  `json:"service_name"`
	ServiceVersion string  `json:"service_version"`
	LogLevel       string  `json:"log_level" validate:"oneof=debug info warn error"`
	OTLPEndpoint   string  `json:"otlp_endpoint"`
	SamplingRatio  float64 `json:"sampling_ratio" validate:"min=0,max=1"`
}

// WorkerConfig holds worker configuration
type WorkerConfig struct {
	TickRate time.Duration `json:"tick_rate" validate:"required,min=1s"`
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	// First, load secrets from Google Secret Manager if CONFIG_FROM_GSM_SECRET is set
	// This must happen before viper reads environment variables
	ctx := context.Background()
	if err := loadFromGSM(ctx); err != nil {
		return nil, fmt.Errorf("load from GSM: %w", err)
	}

	viper.SetConfigType("env")
	viper.SetEnvPrefix("")
	viper.AutomaticEnv()

	// Set defaults
	setDefaults()

	// In development, optionally load from .env file
	if viper.GetString("APP_ENVIRONMENT") == "development" || viper.GetString("APP_ENVIRONMENT") == "" {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig() // Ignore error - fallback to env vars
	}

	// Map env vars to config struct
	config := &Config{
		App: AppConfig{
			Name:        viper.GetString("APP_NAME"),
			Environment: viper.GetString("APP_ENVIRONMENT"),
			Debug:       viper.GetBool("APP_DEBUG"),
		},
		Server: ServerConfig{
			Host:         viper.GetString("SERVER_HOST"),
			Port:         viper.GetInt("SERVER_PORT"),
			ReadTimeout:  viper.GetDuration("SERVER_READ_TIMEOUT"),
			WriteTimeout: viper.GetDuration("SERVER_WRITE_TIMEOUT"),
			IdleTimeout:  viper.GetDuration("SERVER_IDLE_TIMEOUT"),
		},
		Database: DatabaseConfig{
			URL:             viper.GetString("DATABASE_URL"),
			MaxOpenConns:    viper.GetInt("DATABASE_MAX_OPEN_CONNS"),
			MaxIdleConns:    viper.GetInt("DATABASE_MAX_IDLE_CONNS"),
			ConnMaxLifetime: viper.GetDuration("DATABASE_CONN_MAX_LIFETIME"),
		},
		Redis: RedisConfig{
			URL:      viper.GetString("REDIS_URL"),
			PoolSize: viper.GetInt("REDIS_POOL_SIZE"),
		},
		JWT: JWTConfig{
			Secret:     viper.GetString("JWT_SECRET"),
			Issuer:     viper.GetString("JWT_ISSUER"),
			Expiration: viper.GetDuration("JWT_EXPIRATION"),
		},
		Observability: ObservabilityConfig{
			ServiceName:    viper.GetString("OBSERVABILITY_SERVICE_NAME"),
			ServiceVersion: viper.GetString("OBSERVABILITY_SERVICE_VERSION"),
			LogLevel:       viper.GetString("OBSERVABILITY_LOG_LEVEL"),
			OTLPEndpoint:   viper.GetString("OBSERVABILITY_OTLP_ENDPOINT"),
			SamplingRatio:  viper.GetFloat64("OBSERVABILITY_SAMPLING_RATIO"),
		},
		Worker: WorkerConfig{
			TickRate: viper.GetDuration("TICK_RATE"),
		},
	}

	// Validate configuration
	validate := validator.New()
	if err := validate.Struct(config); err != nil {
		return nil, fmt.Errorf("validate config: %w", err)
	}

	return config, nil
}

func setDefaults() {
	// App defaults
	viper.SetDefault("APP_NAME", "email-worker")
	viper.SetDefault("APP_ENVIRONMENT", "development")
	viper.SetDefault("APP_DEBUG", false)

	// Server defaults
	viper.SetDefault("SERVER_HOST", "0.0.0.0")
	viper.SetDefault("SERVER_PORT", 20323)
	viper.SetDefault("SERVER_READ_TIMEOUT", "30s")
	viper.SetDefault("SERVER_WRITE_TIMEOUT", "30s")
	viper.SetDefault("SERVER_IDLE_TIMEOUT", "60s")

	// Database defaults
	viper.SetDefault("DATABASE_MAX_OPEN_CONNS", 25)
	viper.SetDefault("DATABASE_MAX_IDLE_CONNS", 5)
	viper.SetDefault("DATABASE_CONN_MAX_LIFETIME", "5m")

	// Redis defaults
	viper.SetDefault("REDIS_POOL_SIZE", 10)

	// JWT defaults
	viper.SetDefault("JWT_EXPIRATION", "24h")
	viper.SetDefault("JWT_ISSUER", "email-worker")

	// Observability defaults
	viper.SetDefault("OBSERVABILITY_SERVICE_NAME", "email-worker")
	viper.SetDefault("OBSERVABILITY_SERVICE_VERSION", "0.1.0")
	viper.SetDefault("OBSERVABILITY_LOG_LEVEL", "info")
	viper.SetDefault("OBSERVABILITY_OTLP_ENDPOINT", "http://localhost:4318")
	viper.SetDefault("OBSERVABILITY_SAMPLING_RATIO", 0.1)

	// Worker defaults
	viper.SetDefault("TICK_RATE", "30s")
}
