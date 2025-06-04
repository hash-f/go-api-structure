package main

import (
	"context"
	"fmt"
	"io"       // Keep for now, though direct use might be replaced by logger
	"log/slog" // Import slog
	"os"
	"os/signal"
	"syscall"

	"go-api-structure/internal/config"
	"go-api-structure/internal/database" // Import database package
	"go-api-structure/internal/logger" // Import your logger package
)

func main() {
	ctx := context.Background()
	// Pass os.Getenv as the getenv function
	if err := run(ctx, os.Stdout, os.Args, os.Getenv); err != nil {
		// Use a basic slog here if run() fails before logger is initialized,
		// or just fmt.Fprintf for bootstrap errors.
		slog.Error("application startup error", "error", err) // Or fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// Add getenv to the function signature
func run(ctx context.Context, w io.Writer, args []string, getenv func(key string) string) error {
	// Initialize logger first, as it's used for subsequent setup messages.
	// To do this, we need to load config first, or at least the APP_ENV part.
	// For simplicity, load full config, then init logger.

	cfg, err := config.Load(getenv)
	if err != nil {
		// If config fails, we can't initialize the logger as intended.
		// Fallback to a basic logger or direct stderr.
		return fmt.Errorf("failed to load configuration: %w", err)
	}

	// Initialize logger using AppEnv from config
	appLogger := logger.New(cfg.AppEnv)
	slog.SetDefault(appLogger) // Set as default for global use

	slog.Info("Application starting...", "pid", os.Getpid())

	slog.Info("Configuration loaded successfully", "http_port", cfg.HTTPPort, "app_env", cfg.AppEnv)

	// Initialize database connection
	slog.Info("Initializing database connection...", "dsn", cfg.DatabaseDSN) // Consider masking DSN in production logs if sensitive
	db, err := database.NewDB(cfg.DatabaseDSN)
	if err != nil {
		slog.Error("Failed to initialize database connection", "error", err)
		return fmt.Errorf("failed to initialize database: %w", err)
	}
	defer func() {
		slog.Info("Closing database connection pool...")
		if err := db.Close(); err != nil {
			slog.Error("Failed to close database connection pool", "error", err)
		}
	}()
	slog.Info("Database connection established successfully")

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// TODO: Initialize server (using cfg.HTTPPort, db), etc.
	// These components should receive the logger instance and db instance if they need them.

	slog.Info("Application ready and listening")

	// Wait for context cancellation (e.g., SIGINT)
	<-ctx.Done()

	slog.Info("Application shutting down...", "signal", ctx.Err())
	// Graceful shutdown of resources is handled by deferred calls (e.g., db.Close())

	return nil
}
