package config

import (
	"os"
	"testing"
)

func TestLoadDefaultConfig(t *testing.T) {
	// Clear any existing env vars that might interfere
	clearTestEnvVars()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	// Test default values
	if cfg.Server.Port != 34075 {
		t.Errorf("Server.Port = %d, want 34075", cfg.Server.Port)
	}
	if cfg.Server.Host != "localhost" {
		t.Errorf("Server.Host = %s, want localhost", cfg.Server.Host)
	}
	if cfg.Server.Env != "development" {
		t.Errorf("Server.Env = %s, want development", cfg.Server.Env)
	}
	if cfg.Logging.Level != "debug" {
		t.Errorf("Logging.Level = %s, want debug", cfg.Logging.Level)
	}
	if cfg.Logging.Format != "console" {
		t.Errorf("Logging.Format = %s, want console", cfg.Logging.Format)
	}
	if cfg.API.Version != "v1" {
		t.Errorf("API.Version = %s, want v1", cfg.API.Version)
	}
	if cfg.API.Prefix != "/api" {
		t.Errorf("API.Prefix = %s, want /api", cfg.API.Prefix)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name     string
		port     int
		env      string
		logLevel string
		wantErr  bool
	}{
		{
			name:     "valid config",
			port:     34075,
			env:      "development",
			logLevel: "debug",
			wantErr:  false,
		},
		{
			name:     "invalid port low",
			port:     500,
			env:      "development",
			logLevel: "debug",
			wantErr:  true,
		},
		{
			name:     "invalid port high",
			port:     70000,
			env:      "development",
			logLevel: "debug",
			wantErr:  true,
		},
		{
			name:     "invalid env",
			port:     34075,
			env:      "invalid",
			logLevel: "debug",
			wantErr:  true,
		},
		{
			name:     "invalid log level",
			port:     34075,
			env:      "development",
			logLevel: "invalid",
			wantErr:  true,
		},
		{
			name:     "production env valid",
			port:     34075,
			env:      "production",
			logLevel: "info",
			wantErr:  false,
		},
		{
			name:     "staging env valid",
			port:     34075,
			env:      "staging",
			logLevel: "warn",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{
					Port: tt.port,
					Host: "localhost",
					Env:  tt.env,
				},
				Logging: LoggingConfig{
					Level:  tt.logLevel,
					Format: "console",
				},
				API: APIConfig{
					Version: "v1",
					Prefix:  "/api",
				},
				CORS: CORSConfig{
					AllowedOrigins: []string{"http://localhost:3000"},
					AllowedMethods: []string{"GET", "POST"},
					AllowedHeaders: []string{"Content-Type"},
				},
			}

			err := cfg.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestEnvOverride(t *testing.T) {
	// Set env vars
	_ = os.Setenv("PORT", "9090")
	_ = os.Setenv("HOST", "0.0.0.0")
	_ = os.Setenv("ENV", "production")
	_ = os.Setenv("LOG_LEVEL", "error")
	_ = os.Setenv("LOG_FORMAT", "json")
	defer clearTestEnvVars()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	if cfg.Server.Port != 9090 {
		t.Errorf("Server.Port = %d, want 9090", cfg.Server.Port)
	}
	if cfg.Server.Host != "0.0.0.0" {
		t.Errorf("Server.Host = %s, want 0.0.0.0", cfg.Server.Host)
	}
	if cfg.Server.Env != "production" {
		t.Errorf("Server.Env = %s, want production", cfg.Server.Env)
	}
	if cfg.Logging.Level != "error" {
		t.Errorf("Logging.Level = %s, want error", cfg.Logging.Level)
	}
	if cfg.Logging.Format != "json" {
		t.Errorf("Logging.Format = %s, want json", cfg.Logging.Format)
	}
}

func TestCORSOverride(t *testing.T) {
	_ = os.Setenv("CORS_ALLOWED_ORIGINS", "https://masquerade.app,https://api.masquerade.app")
	_ = os.Setenv("CORS_ALLOWED_METHODS", "GET,POST,PUT")
	_ = os.Setenv("CORS_ALLOWED_HEADERS", "Content-Type,Authorization,X-Custom-Header")
	defer clearTestEnvVars()

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load() error = %v, want nil", err)
	}

	expectedOrigins := []string{"https://masquerade.app", "https://api.masquerade.app"}
	if len(cfg.CORS.AllowedOrigins) != len(expectedOrigins) {
		t.Fatalf("CORS.AllowedOrigins length = %d, want %d", len(cfg.CORS.AllowedOrigins), len(expectedOrigins))
	}
	for i, origin := range expectedOrigins {
		if cfg.CORS.AllowedOrigins[i] != origin {
			t.Errorf("CORS.AllowedOrigins[%d] = %s, want %s", i, cfg.CORS.AllowedOrigins[i], origin)
		}
	}

	expectedMethods := []string{"GET", "POST", "PUT"}
	if len(cfg.CORS.AllowedMethods) != len(expectedMethods) {
		t.Fatalf("CORS.AllowedMethods length = %d, want %d", len(cfg.CORS.AllowedMethods), len(expectedMethods))
	}

	expectedHeaders := []string{"Content-Type", "Authorization", "X-Custom-Header"}
	if len(cfg.CORS.AllowedHeaders) != len(expectedHeaders) {
		t.Fatalf("CORS.AllowedHeaders length = %d, want %d", len(cfg.CORS.AllowedHeaders), len(expectedHeaders))
	}
}

func TestMustLoad(t *testing.T) {
	clearTestEnvVars()

	// Should not panic with valid config
	cfg := MustLoad()
	if cfg == nil {
		t.Error("MustLoad() returned nil")
	}
}

func TestMustLoadPanic(t *testing.T) {
	_ = os.Setenv("PORT", "invalid")
	defer clearTestEnvVars()

	defer func() {
		if r := recover(); r == nil {
			t.Error("MustLoad() did not panic with invalid PORT")
		}
	}()

	MustLoad()
}

func TestIsDevelopment(t *testing.T) {
	tests := []struct {
		env  string
		want bool
	}{
		{"development", true},
		{"staging", false},
		{"production", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{Env: tt.env},
			}
			if got := cfg.IsDevelopment(); got != tt.want {
				t.Errorf("IsDevelopment() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsProduction(t *testing.T) {
	tests := []struct {
		env  string
		want bool
	}{
		{"production", true},
		{"staging", false},
		{"development", false},
	}

	for _, tt := range tests {
		t.Run(tt.env, func(t *testing.T) {
			cfg := &Config{
				Server: ServerConfig{Env: tt.env},
			}
			if got := cfg.IsProduction(); got != tt.want {
				t.Errorf("IsProduction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetEnv(t *testing.T) {
	_ = os.Setenv("TEST_VAR", "test_value")
	defer func() { _ = os.Unsetenv("TEST_VAR") }()

	if got := getEnv("TEST_VAR", "default"); got != "test_value" {
		t.Errorf("getEnv() = %s, want test_value", got)
	}

	if got := getEnv("NON_EXISTENT", "default"); got != "default" {
		t.Errorf("getEnv() = %s, want default", got)
	}
}

func TestGetEnvAsInt(t *testing.T) {
	_ = os.Setenv("TEST_INT", "42")
	defer func() { _ = os.Unsetenv("TEST_INT") }()

	if got := getEnvAsInt("TEST_INT", 0); got != 42 {
		t.Errorf("getEnvAsInt() = %d, want 42", got)
	}

	if got := getEnvAsInt("NON_EXISTENT", 99); got != 99 {
		t.Errorf("getEnvAsInt() = %d, want 99", got)
	}
}

func TestGetEnvAsIntPanic(t *testing.T) {
	_ = os.Setenv("TEST_INT", "not_a_number")
	defer func() { _ = os.Unsetenv("TEST_INT") }()

	defer func() {
		if r := recover(); r == nil {
			t.Error("getEnvAsInt() did not panic with invalid integer")
		}
	}()

	getEnvAsInt("TEST_INT", 0)
}

func TestGetEnvAsSlice(t *testing.T) {
	_ = os.Setenv("TEST_SLICE", "a,b,c")
	defer func() { _ = os.Unsetenv("TEST_SLICE") }()

	got := getEnvAsSlice("TEST_SLICE", []string{"default"})
	expected := []string{"a", "b", "c"}

	if len(got) != len(expected) {
		t.Fatalf("getEnvAsSlice() length = %d, want %d", len(got), len(expected))
	}

	for i, v := range expected {
		if got[i] != v {
			t.Errorf("getEnvAsSlice()[%d] = %s, want %s", i, got[i], v)
		}
	}

	got = getEnvAsSlice("NON_EXISTENT", []string{"default"})
	if len(got) != 1 || got[0] != "default" {
		t.Errorf("getEnvAsSlice() = %v, want [default]", got)
	}
}

// Helper function to clear test environment variables
func clearTestEnvVars() {
	_ = os.Unsetenv("PORT")
	_ = os.Unsetenv("HOST")
	_ = os.Unsetenv("ENV")
	_ = os.Unsetenv("LOG_LEVEL")
	_ = os.Unsetenv("LOG_FORMAT")
	_ = os.Unsetenv("API_VERSION")
	_ = os.Unsetenv("API_PREFIX")
	_ = os.Unsetenv("CORS_ALLOWED_ORIGINS")
	_ = os.Unsetenv("CORS_ALLOWED_METHODS")
	_ = os.Unsetenv("CORS_ALLOWED_HEADERS")
	_ = os.Unsetenv("TEST_VAR")
	_ = os.Unsetenv("TEST_INT")
	_ = os.Unsetenv("TEST_SLICE")
}
