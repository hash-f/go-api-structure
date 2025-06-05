package server

import (
	"log/slog"
	"net/http"

	_ "go-api-structure/internal/docs" // Import for swagger docs generation

	"go-api-structure/internal/api"
	"go-api-structure/internal/auth"
	"go-api-structure/internal/config"
	"go-api-structure/internal/store"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Server holds the dependencies for the HTTP server.
// This typically includes the application configuration, logger,
// data stores, and the router itself.
type Server struct {
	config      *config.Config
	logger      *slog.Logger
	store       store.Store
	router      *chi.Mux
	authService *auth.AuthService
	authHandler *api.AuthHandler
	userHandler *api.UserHandler
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

	// Initialize services and handlers
	s.initDependencies()

	// Global middleware
	s.addMiddlewares()

	// Setup routes
	s.addRoutes()

	return s.router
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

func (s *Server) initDependencies() {
	s.authService = auth.NewAuthService(s.store, s.config.JWTSecret, s.config.JWTExpiryDuration)
	s.authHandler = api.NewAuthHandler(s.authService)
	s.userHandler = api.NewUserHandler(s.store)
}

func (s *Server) addMiddlewares() {
	s.router.Use(middleware.RequestID)           // Injects a request ID into the context
	s.router.Use(middleware.RealIP)              // Sets X-Real-IP and X-Forwarded-For
	s.router.Use(createSlogMiddleware(s.logger)) // Custom slog logging middleware
	s.router.Use(middleware.Recoverer)           // Recovers from panics
	s.router.Use(createCorsMiddleware())         // CORS configuration
}
