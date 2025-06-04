package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	// Import your config package
	"go-api-structure/internal/config" // Assuming your module is go-api-structure
)

func main() {
	ctx := context.Background()
	// Pass os.Getenv as the getenv function
	if err := run(ctx, os.Stdout, os.Args, os.Getenv); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

// Add getenv to the function signature
func run(ctx context.Context, w io.Writer, args []string, getenv func(key string) string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Fprintln(w, "Application starting...")

	// Load configuration
	cfg, err := config.Load(getenv)
	if err != nil {
		return fmt.Errorf("failed to load configuration: %w", err)
	}
	fmt.Fprintf(w, "Configuration loaded successfully. HTTP Port: %s\n", cfg.HTTPPort)

	// TODO: Initialize logger (using cfg if needed), database (using cfg.DatabaseDSN), server (using cfg.HTTPPort), etc.

	// Wait for context cancellation (e.g., SIGINT)
	<-ctx.Done()

	fmt.Fprintln(w, "Application shutting down...")
	// TODO: Graceful shutdown of resources

	return nil
}
