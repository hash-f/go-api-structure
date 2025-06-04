package main

import (
	"context"
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"go-api-structure/internal/config"
	"go-api-structure/internal/database"
	"go-api-structure/internal/logger"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args, os.Getenv); err != nil {
		// If run() fails, this slog.Error will use the default Go logger
		// if setupLogger hasn't been called yet, or our configured logger if it has.
		slog.Error("application startup error", "error", err)
		os.Exit(1)
	}
}

// loadConfiguration loads application configuration.
func loadConfiguration(getenv func(string) string) (*config.Config, error) {
	cfg, err := config.Load(getenv)
	if err != nil {
		// It's okay to return a plainly formatted error here.
		// The caller (run) will return it, and main() will log it using slog.Error.
		// At that point, slog.Default() will be the Go standard logger if our logger isn't set up yet.
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	return cfg, nil
}

// setupLogger initializes the global logger based on the loaded configuration.
func setupLogger(cfg *config.Config) {
	appLogger := logger.New(cfg.AppEnv)
	slog.SetDefault(appLogger) // Set as default for global use

	// These messages will now use the configured logger.
	slog.Info("Application starting...", "pid", os.Getpid())
	slog.Info("Configuration loaded successfully", "http_port", cfg.HTTPPort, "app_env", cfg.AppEnv)
}

// setupDatabase initializes and returns a new database connection pool.
func setupDatabase(dsn string) (*sql.DB, error) {
	slog.Info("Initializing database connection...", "dsn", dsn)
	db, err := database.NewDB(dsn)
	if err != nil {
		slog.Error("Failed to initialize database connection", "error", err)
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}
	slog.Info("Database connection established successfully")
	return db, nil
}

func run(ctx context.Context, w io.Writer, args []string, getenv func(key string) string) error {
	cfg, err := loadConfiguration(getenv)
	if err != nil {
		// If config loading fails, this error will be returned to main(),
		// which will log it using the default slog handler (as our logger isn't set up).
		return err
	}

	// Config loaded successfully, now set up the logger.
	// Subsequent slog calls will use this configured logger.
	setupLogger(cfg)

	db, err := setupDatabase(cfg.DatabaseDSN)
	if err != nil {
		// This error will be logged by setupDatabase using our configured logger,
		// and then returned to main().
		return err
	}
	defer func() {
		slog.Info("Closing database connection pool...")
		if err := db.Close(); err != nil {
			slog.Error("Failed to close database connection pool", "error", err)
		}
	}()

	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	slog.Info("Application ready and listening")

	<-ctx.Done()

	slog.Info("Application shutting down...", "signal", ctx.Err())
	return nil
}
