package api

import (
	"errors" // For store.ErrNotFound
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"go-api-structure/internal/api/dto"
	"go-api-structure/internal/auth"
	"go-api-structure/internal/store" // For store.ErrNotFound
	"go-api-structure/internal/user"  // New import
)

// UserHandler holds dependencies for user-related HTTP handlers.
type UserHandler struct {
	userService user.ServiceInterface // Changed from userStore
}

// NewUserHandler creates a new UserHandler with necessary dependencies.
func NewUserHandler(userService user.ServiceInterface) *UserHandler { // Changed parameter
	return &UserHandler{
		userService: userService, // Changed assignment
	}
}

// @Summary      Get current user's details
// @Description  Retrieves the details of the currently authenticated user.
// @Tags         Users
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  dto.UserResponse "Successfully retrieved user details"
// @Failure      401  {object}  map[string]string "Unauthorized (e.g., no user in context, invalid token)"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /users/me [get]
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

// @Summary      Get user details by ID
// @Description  Retrieves the details of a user by their ID.
// @Tags         Users
// @Produce      json
// @Param        id   path      string  true  "User ID (UUID format)"
// @Security     APIKey
// @Success      200  {object}  dto.UserResponse "Successfully retrieved user details"
// @Failure      400  {object}  map[string]string "Invalid user ID format"
// @Failure      401  {object}  map[string]string "Unauthorized (e.g., invalid API key)"
// @Failure      404  {object}  map[string]string "User not found"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /users/{id} [get]
func (h *UserHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	// The APIKeyMiddleware should have already authenticated the request.
	// We don't need to re-check the API key here.
	// We can optionally retrieve the authenticated user from context if needed for authorization,
	// e.g., if only admins can fetch any user, or users can only fetch themselves.
	// For now, let's assume if the API key is valid, the request is allowed to proceed
	// to fetch the user specified by the path ID.

	userIDStr := chi.URLParam(r, "id")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ErrorResponse(w, r, http.StatusBadRequest, "Invalid user ID format")
		return
	}

	targetUser, err := h.userService.GetUserByID(r.Context(), userID)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			ErrorResponse(w, r, http.StatusNotFound, "User not found")
		} else {
			// Log the error for internal review
			ErrorResponse(w, r, http.StatusInternalServerError, "Failed to retrieve user")
		}
		return
	}

	userResponse := dto.NewUserResponse(targetUser)
	encode(w, r, http.StatusOK, userResponse)
}
