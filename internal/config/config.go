package config

import (
	"fmt"
	"os" // Used for default in Load if getenv is nil, or for direct use if preferred

	"github.com/joho/godotenv" // For loading .env files
)

// Config holds the application configuration.
type Config struct {
	AppEnv      string // e.g., "local", "dev", "stage", "production"
	HTTPPort    string
	DatabaseDSN string
	JWTSecret   string
	// Add other configuration fields as needed
}

// Load reads configuration from environment variables.
// It uses the provided getenv function, which makes it testable.
// If a .env file exists, it will be loaded first.
func Load(getenv func(key string) string) (*Config, error) {
	// Attempt to load .env file. This is useful for local development.
	// In production, environment variables are usually set directly.
	_ = godotenv.Load() // Ignore error if .env file doesn't exist

	if getenv == nil {
		getenv = os.Getenv // Default to os.Getenv if no custom function is provided
	}

	cfg := &Config{}

	cfg.AppEnv = getenv("APP_ENV")
	if cfg.AppEnv == "" {
		cfg.AppEnv = "local" // Default to local environment
	}

	cfg.HTTPPort = getenv("HTTP_PORT")
	if cfg.HTTPPort == "" {
		cfg.HTTPPort = "8080" // Default port
	}

	cfg.DatabaseDSN = getenv("DATABASE_DSN")
	if cfg.DatabaseDSN == "" {
		return nil, fmt.Errorf("DATABASE_DSN environment variable is required")
	}

	cfg.JWTSecret = getenv("JWT_SECRET")
	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	// TODO: Load other configuration fields

	return cfg, nil
}
