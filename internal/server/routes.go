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
	s.router.Route("/api/v1", func(r chi.Router) {
		// Authentication routes (e.g., /api/v1/auth/register, /api/v1/auth/login)
		r.Route("/auth", func(r chi.Router) {
			r.Post("/register", s.authHandler.RegisterUser)
			r.Post("/login", s.authHandler.LoginUser)
		})

		// User routes (e.g., /api/v1/users/me)
		r.Route("/users", func(r chi.Router) {
			// Protected routes - require JWT authentication
			r.Group(func(r chi.Router) {
				r.Use(s.authService.Middleware(api.ErrorResponse)) // Pass api.ErrorResponse as the error renderer
				r.Get("/me", s.userHandler.GetMe)
			})
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
