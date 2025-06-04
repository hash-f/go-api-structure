package server

import (
	"context"
	"go-api-structure/internal/api"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5/middleware"
)

// slogMiddleware is a custom logging middleware using slog.
// It logs request details similar to chi's built-in logger but uses the structured logger.
func slogMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tstart := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			ctx := context.WithValue(r.Context(), api.GetLoggerKey(), logger)
			r = r.WithContext(ctx)

			defer func() {
				logger.Info("Served request",
					"status", ww.Status(),
					"method", r.Method,
					"path", r.URL.Path,
					"query", r.URL.RawQuery,
					"request_id", middleware.GetReqID(r.Context()),
					"duration", time.Since(tstart),
					"bytes_written", ww.BytesWritten(),
				)
			}()

			next.ServeHTTP(ww, r)
		}
		return http.HandlerFunc(fn)
	}
}
