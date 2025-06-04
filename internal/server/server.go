package server

import (
	"log/slog"
	"net/http"
	"time"

	"go-api-structure/internal/config"
	"go-api-structure/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Server holds the dependencies for the HTTP server.
// This typically includes the application configuration, logger,
// data stores, and the router itself.
type Server struct {
	config *config.Config
	logger *slog.Logger
	store  store.Store
	router *chi.Mux
}

// NewServer creates and configures a new Server instance.
// It initializes the router, sets up dependencies, and prepares the server
// to handle requests. It returns an http.Handler (the configured router)
// which can be used with http.ListenAndServe.
func NewServer(cfg *config.Config, logger *slog.Logger, store store.Store) http.Handler {
	s := &Server{
		config: cfg,
		logger: logger,
		store:  store,
		router: chi.NewRouter(), // Initialize the chi router
	}

	// Global middleware
	s.router.Use(middleware.RequestID) // Injects a request ID into the context
	s.router.Use(middleware.RealIP)    // Sets X-Real-IP and X-Forwarded-For
	s.router.Use(slogMiddleware(logger)) // Custom slog logging middleware
	s.router.Use(middleware.Recoverer) // Recovers from panics

	// CORS configuration
	s.router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"}, // Allow all for now, tighten in production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any major browsers
	}))

	// TODO: Setup routes - Subtask 3.2
	// s.addRoutes()

	return s.router
}

// addRoutes will be responsible for setting up all the API routes.
// This will be implemented as part of Subtask 3.2.
// func (s *Server) addRoutes() {
// 	 s.router.Get("/health", s.handleHealthCheck())
// }

// handleHealthCheck is a simple handler for health checks.
// func (s *Server) handleHealthCheck() http.HandlerFunc {
// 	 return func(w http.ResponseWriter, r *http.Request) {
// 		 // TODO: Use JSON response helper from Subtask 3.3
// 		 w.WriteHeader(http.StatusOK)
// 		 w.Write([]byte("OK"))
// 	 }
// }

// slogMiddleware is a custom logging middleware using slog.
// It logs request details similar to chi's built-in logger but uses the structured logger.
func slogMiddleware(logger *slog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tstart := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

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
