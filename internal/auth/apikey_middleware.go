package auth

import (
	"net/http"

	"go-api-structure/internal/user" // For user.ServiceInterface
)

const (
	APIKeyHeader = "X-API-Key" // Standard header for API keys
)

// APIKeyMiddleware creates a middleware that authenticates requests using an API key.
// It expects the AuthService to have a UserService instance provided to it.
func (s *AuthService) APIKeyMiddleware(userService user.ServiceInterface, errorFunc func(w http.ResponseWriter, r *http.Request, statusCode int, message any)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get(APIKeyHeader)
			if apiKey == "" {
				errorFunc(w, r, http.StatusUnauthorized, "API key required")
				return
			}

			dbUser, err := userService.GetUserByAPIKey(r.Context(), apiKey)
			if err != nil {
				// Consider logging the error for internal review: log.Printf("API key validation error: %v", err)
				// For the client, it's an invalid API key.
				errorFunc(w, r, http.StatusUnauthorized, "Invalid API key")
				return
			}

			// Add user to context using the existing mechanism (assuming ContextSetUser is available)
			ctxWithUser := ContextSetUser(r.Context(), dbUser)
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		})
	}
}
