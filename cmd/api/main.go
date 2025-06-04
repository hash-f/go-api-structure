package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"go-api-structure/internal/config"
	"go-api-structure/internal/database"
	"go-api-structure/internal/logger"
	"go-api-structure/internal/server"
	"go-api-structure/internal/store"
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

// setupDatabase initializes and returns a new pgx database connection pool.
func setupDatabase(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	slog.Info("Initializing database connection pool (pgx)...", "dsn", dsn)
	pool, err := database.NewPgxPool(ctx, dsn)
	if err != nil {
		slog.Error("Failed to initialize database connection pool (pgx)", "error", err)
		return nil, fmt.Errorf("failed to initialize database (pgx): %w", err)
	}
	slog.Info("Database connection pool (pgx) established successfully")
	return pool, nil
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

	// Pass the main context to setupDatabase for pgxpool.New
	db, err := setupDatabase(ctx, cfg.DatabaseDSN)
	if err != nil {
		// This error will be logged by setupDatabase using our configured logger,
		// and then returned to main().
		return err
	}
	defer func() {
		slog.Info("Closing database connection pool (pgx)...")
		db.Close() // pgxpool.Pool.Close() doesn't return an error
		slog.Info("Database connection pool (pgx) closed.")
	}()

	appStore := store.NewStore(db)

	// Get the configured logger. slog.Default() returns the logger set by setupLogger.
	appLogger := slog.Default()

	httpHandler := server.NewServer(cfg, appLogger, appStore)

	srv := &http.Server{
		Addr:         ":" + cfg.HTTPPort,
		Handler:      httpHandler,
		ReadTimeout:  5 * time.Second,  // Example timeout
		WriteTimeout: 10 * time.Second, // Example timeout
		IdleTimeout:  15 * time.Second, // Example timeout
	}

	// Channel to listen for errors from the HTTP server goroutine
	serverErrors := make(chan error, 1)

	// Start HTTP server in a goroutine
	go func() {
		slog.Info("HTTP server starting", "address", srv.Addr)
		serverErrors <- srv.ListenAndServe()
	}()

	// Graceful shutdown
	ctxShutdown, cancelShutdown := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancelShutdown()

	select {
	case err := <-serverErrors:
		if !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "error", err)
			return fmt.Errorf("http server error: %w", err)
		}
		slog.Info("HTTP server stopped.")
	case <-ctxShutdown.Done():
		slog.Info("Shutdown signal received, initiating graceful shutdown...", "signal", ctxShutdown.Err())
		shutdownCtx, cancelTimeout := context.WithTimeout(context.Background(), 10*time.Second) // Timeout for shutdown
		defer cancelTimeout()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			slog.Error("HTTP server graceful shutdown failed", "error", err)
			return fmt.Errorf("http server graceful shutdown failed: %w", err)
		}
		slog.Info("HTTP server shutdown gracefully.")
	}

	return nil
}
