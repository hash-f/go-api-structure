package main

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx := context.Background()
	if err := run(ctx, os.Stdout, os.Args); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context, w io.Writer, args []string) error {
	ctx, cancel := signal.NotifyContext(ctx, os.Interrupt, syscall.SIGTERM)
	defer cancel()

	fmt.Fprintln(w, "Application starting...")

	// TODO: Initialize configuration, logger, database, server, etc.

	// Wait for context cancellation (e.g., SIGINT)
	<-ctx.Done()

	fmt.Fprintln(w, "Application shutting down...")
	// TODO: Graceful shutdown of resources

	return nil
}
