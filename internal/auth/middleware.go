package auth

import (
	// "context" // No longer directly used here
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	// "go-api-structure/internal/api" // Removed to break import cycle
	"go-api-structure/internal/store"
	// "go-api-structure/internal/store/db" // No longer directly used here
)

// Middleware is a JWT authentication middleware.
// It checks for a valid JWT in the Authorization header.
// If the token is valid, it retrieves the user from the store and adds them to the request context using ContextSetUser.
// It calls the provided errorRenderer for sending HTTP error responses.
func (s *AuthService) Middleware(errorRenderer func(w http.ResponseWriter, r *http.Request, status int, message any)) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				errorRenderer(w, r, http.StatusUnauthorized, "authorization header missing")
				return
			}

			headerParts := strings.Split(authHeader, " ")
			if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
				errorRenderer(w, r, http.StatusUnauthorized, "invalid authorization header format")
				return
			}

			tokenString := headerParts[1]

			claims := &jwt.RegisteredClaims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
				}
				return []byte(s.jwtSecret), nil
			})

			if err != nil {
				if errors.Is(err, jwt.ErrTokenExpired) {
					errorRenderer(w, r, http.StatusUnauthorized, "token expired")
				} else {
					errorRenderer(w, r, http.StatusUnauthorized, "invalid token")
				}
				return
			}

			if !token.Valid {
				errorRenderer(w, r, http.StatusUnauthorized, "invalid token")
				return
			}

			userID, err := uuid.Parse(claims.Subject)
			if err != nil {
				// Log this error as it's likely a server-side issue with token generation or parsing claims
				// For the client, it's still an invalid token situation or internal error.
				errorRenderer(w, r, http.StatusUnauthorized, "invalid token claims") // Or InternalServerError
				return
			}

			user, err := s.userStore.GetUserByID(r.Context(), userID)
			if err != nil {
				if errors.Is(err, store.ErrNotFound) {
					errorRenderer(w, r, http.StatusUnauthorized, "user not found")
				} else {
					// Log internal error
					errorRenderer(w, r, http.StatusInternalServerError, "error retrieving user")
				}
				return
			}

			// Add user to context
			ctxWithUser := ContextSetUser(r.Context(), &user) // Use ContextSetUser from context.go
			next.ServeHTTP(w, r.WithContext(ctxWithUser))
		})
	}
}
