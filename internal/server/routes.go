package server

import (
	"go-api-structure/internal/api"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger" // Swagger UI handler
)

// addRoutes is responsible for setting up all the API routes.
func (s *Server) addRoutes() {
	// Health check endpoint
	s.router.Get("/health", s.handleHealthCheck())

	// Swagger UI endpoint
	s.router.Get("/swagger/*", httpSwagger.WrapHandler)

	// API v1 routes
	s.router.Route("/api/v1", s.apiRoutes)
}

func (s *Server) apiRoutes(r chi.Router) {
	// Authentication routes (e.g., /api/v1/auth/register, /api/v1/auth/login)
	r.Route("/auth", s.apiAuthRoutes)

	// User routes (e.g., /api/v1/users/me)
	r.Route("/users", s.apiUserRoutes)
}

func (s *Server) apiAuthRoutes(r chi.Router) {
	r.Post("/register", s.authHandler.RegisterUser)
	r.Post("/login", s.authHandler.LoginUser)
}

func (s *Server) apiUserRoutes(r chi.Router) {
	// Protected routes - require JWT authentication
	r.Group(func(r chi.Router) {
		r.Use(s.authService.JWTMiddleware(api.ErrorResponse))
		r.Get("/me", s.userHandler.GetMe)
	})
}
