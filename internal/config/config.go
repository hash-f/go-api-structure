package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv" // For loading .env files
)

// Config holds the application configuration.
type Config struct {
	AppEnv            string // e.g., "local", "dev", "stage", "production"
	HTTPPort          string
	DatabaseDSN       string
	JWTSecret         string
	JWTExpiryDuration time.Duration
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

	jwtExpiryMinutesStr := getenv("JWT_EXPIRY_MINUTES")
	if jwtExpiryMinutesStr == "" {
		jwtExpiryMinutesStr = "60" // Default to 60 minutes
	}
	jwtExpiryMinutes, err := strconv.Atoi(jwtExpiryMinutesStr)
	if err != nil {
		return nil, fmt.Errorf("invalid JWT_EXPIRY_MINUTES: %w", err)
	}
	cfg.JWTExpiryDuration = time.Duration(jwtExpiryMinutes) * time.Minute

	// Add loading for other config fields here

	return cfg, nil
}
