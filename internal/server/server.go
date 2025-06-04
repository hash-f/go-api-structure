package server

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"go-api-structure/internal/api"
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

	// Setup routes
	s.addRoutes()

	return s.router
}

// addRoutes is responsible for setting up all the API routes.
func (s *Server) addRoutes() {
	// Health check endpoint
	s.router.Get("/health", s.handleHealthCheck())

	// API v1 routes
	s.router.Route("/api/v1", func(r chi.Router) {
		// Authentication routes (e.g., /api/v1/auth/register, /api/v1/auth/login)
		r.Route("/auth", func(r chi.Router) {
			// r.Post("/register", s.handleRegisterUser()) // Placeholder
			// r.Post("/login", s.handleLoginUser())       // Placeholder
		})

		// User routes (e.g., /api/v1/users/me)
		r.Route("/users", func(r chi.Router) {
			// r.Get("/me", s.handleGetUserMe()) // Placeholder, requires auth middleware
		})

		// Vendor routes (e.g., /api/v1/vendors)
		r.Route("/vendors", func(r chi.Router) {
			// r.Post("/", s.handleCreateVendor())    // Placeholder, requires auth
			// r.Get("/", s.handleListUserVendors()) // Placeholder, requires auth
			// r.Get("/{vendorID}", s.handleGetVendor()) // Placeholder, requires auth & ownership
			// r.Put("/{vendorID}", s.handleUpdateVendor()) // Placeholder, requires auth & ownership
			// r.Delete("/{vendorID}", s.handleDeleteVendor()) // Placeholder, requires auth & ownership
		})

		// Merchant routes (e.g., /api/v1/merchants)
		r.Route("/merchants", func(r chi.Router) {
			// r.Post("/", s.handleCreateMerchant())    // Placeholder, requires auth
			// r.Get("/", s.handleListUserMerchants()) // Placeholder, requires auth
			// r.Get("/{merchantID}", s.handleGetMerchant()) // Placeholder, requires auth & ownership
			// r.Put("/{merchantID}", s.handleUpdateMerchant()) // Placeholder, requires auth & ownership
			// r.Delete("/{merchantID}", s.handleDeleteMerchant()) // Placeholder, requires auth & ownership
		})
	})
}

// handleHealthCheck is a simple handler for health checks.
func (s *Server) handleHealthCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// For now, a simple response. Later, we'll use JSON helpers (Subtask 3.3).
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status": "ok"}`))
	}
}

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
