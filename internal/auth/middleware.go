package auth

import (
	"context"
	"errors" // Added missing import
	"fmt"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go-api-structure/internal/api"
	"go-api-structure/internal/store"
	"go-api-structure/internal/store/db" // For db.User
)

// Define a context key type for the authenticated user.
// This will be used by middleware to inject the user and by handlers to retrieve it.
type contextKey string

const userContextKey = contextKey("user")

// Authenticate is a middleware that checks for a valid JWT in the Authorization header.
// If the token is valid, it retrieves the user from the store and adds them to the request context.
func (s *AuthService) Authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			api.ErrorResponse(w, r, http.StatusUnauthorized, "authorization header missing")
			return
		}

		headerParts := strings.Split(authHeader, " ")
		if len(headerParts) != 2 || strings.ToLower(headerParts[0]) != "bearer" {
			api.ErrorResponse(w, r, http.StatusUnauthorized, "invalid authorization header format")
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
				api.ErrorResponse(w, r, http.StatusUnauthorized, "token expired")
			} else {
				api.ErrorResponse(w, r, http.StatusUnauthorized, "invalid token")
			}
			return
		}

		if !token.Valid {
			api.ErrorResponse(w, r, http.StatusUnauthorized, "invalid token")
			return
		}

		userID, err := uuid.Parse(claims.Subject)
		if err != nil {
			api.ErrorResponse(w, r, http.StatusInternalServerError, "error parsing user ID from token")
			return
		}

		user, err := s.userStore.GetUserByID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, store.ErrNotFound) {
				api.ErrorResponse(w, r, http.StatusUnauthorized, "user not found")
			} else {
				api.ErrorResponse(w, r, http.StatusInternalServerError, "error retrieving user")
			}
			return
		}

		// Add user to context
		ctx := context.WithValue(r.Context(), userContextKey, &user) // Store pointer to user
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserFromContext retrieves the authenticated user from the request context.
// It returns nil if no user is found in the context.
func GetUserFromContext(ctx context.Context) *db.User {
	user, ok := ctx.Value(userContextKey).(*db.User)
	if !ok {
		return nil
	}
	return user
}
