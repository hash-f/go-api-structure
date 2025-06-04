package logger

import (
	"log/slog"
	"os"
)

// New initializes and returns a new slog.Logger based on the application environment.
// appEnv: "local", "dev", "stage", "production"
func New(appEnv string) *slog.Logger {
	var handler slog.Handler

	switch appEnv {
	case "local":
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level:     slog.LevelDebug, // More verbose for local development
			AddSource: true,            // Include source file and line number
		})
	case "dev", "stage", "production":
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	default:
		// Fallback to a sensible default if appEnv is not recognized
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		})
	}

	logger := slog.New(handler)
	return logger
}
