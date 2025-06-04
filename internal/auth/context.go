package auth

import (
	"context"

	"go-api-structure/internal/store/db"
)

// contextKey is an unexported type for context keys defined in this package.
// This prevents collisions with context keys defined in other packages.
type contextKey string

// userContextKey is the key used to store the authenticated user in the request context.
const userContextKey = contextKey("user")

// ContextSetUser adds the user to the given context with the userContextKey.
func ContextSetUser(ctx context.Context, user *db.User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// GetUserFromContext retrieves the authenticated user from the request context.
// It returns nil if no user is found in the context or if the type is incorrect.
func GetUserFromContext(ctx context.Context) *db.User {
	user, ok := ctx.Value(userContextKey).(*db.User)
	if !ok {
		return nil
	}
	return user
}
