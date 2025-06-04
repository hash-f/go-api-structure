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

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	// TODO: Initialize database (using cfg.DatabaseDSN), server (using cfg.HTTPPort), etc.
	// These components should receive the logger instance if they need to log.

	slog.Info("Application ready and listening")

	// Wait for context cancellation (e.g., SIGINT)
	<-ctx.Done()

	slog.Info("Application shutting down...", "signal", ctx.Err())
	// TODO: Graceful shutdown of resources

	return nil
}
