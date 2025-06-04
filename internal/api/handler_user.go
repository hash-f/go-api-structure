package api

import (
	"net/http"

	"go-api-structure/internal/api/dto"
	"go-api-structure/internal/auth"
	"go-api-structure/internal/store"
)

// UserHandler holds dependencies for user-related HTTP handlers.
type UserHandler struct {
	userStore store.UserStore
	// We might not need authService directly if middleware handles auth
	// and GetUserFromContext is a standalone function.
}

// NewUserHandler creates a new UserHandler with necessary dependencies.
func NewUserHandler(userStore store.UserStore) *UserHandler {
	return &UserHandler{
		userStore: userStore,
	}
}

// GetMe handles requests for the authenticated user's details.
// It expects the user to be authenticated by the JWT middleware.
func (h *UserHandler) GetMe(w http.ResponseWriter, r *http.Request) {
	user := auth.GetUserFromContext(r.Context())
	if user == nil {
		// This case should ideally be prevented by the auth middleware.
		// If it occurs, it implies a misconfiguration or an issue with the middleware chain.
		ErrorResponse(w, r, http.StatusUnauthorized, "no authenticated user found in context")
		return
	}

	userResponse := dto.NewUserResponse(user)
	encode(w, r, http.StatusOK, userResponse)
}
