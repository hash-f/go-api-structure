package api

import (
	"errors"
	"net/http"

	"go-api-structure/internal/api/dto"
	"go-api-structure/internal/auth"
	"go-api-structure/internal/store" // For store.ErrNotFound
)

// AuthHandler holds dependencies for authentication-related HTTP handlers.
type AuthHandler struct {
	authService *auth.AuthService
}

// NewAuthHandler creates a new AuthHandler with the given AuthService.
func NewAuthHandler(authService *auth.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// @Summary      Register a new user
// @Description  Creates a new user account with the provided details.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        user body dto.CreateUserRequest true "User registration details"
// @Success      201  {object}  dto.UserResponse "Successfully registered user"
// @Failure      400  {object}  map[string]string "Bad request (e.g., malformed JSON)"
// @Failure      409  {object}  map[string]string "Conflict (user already exists)"
// @Failure      422  {object}  map[string]string "Unprocessable entity (validation error)"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /auth/register [post]
// RegisterUser handles user registration requests.
// It expects a JSON body conforming to dto.CreateUserRequest.
// On success, it returns a 201 Created status with the new user's details (excluding password).
// On failure, it returns appropriate error responses (e.g., 400 for bad request, 422 for validation errors, 409 for conflict).
func (h *AuthHandler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var input dto.CreateUserRequest

	if !decodeAndValidate(w, r, &input) { // decodeAndValidate needs a pointer
		return // Errors handled by decodeAndValidate
	}

	// authService.Register expects *dto.CreateUserRequest (pointer type)
	createdUser, err := h.authService.Register(r.Context(), &input)
	if err != nil {
		switch err {
		case auth.ErrUserAlreadyExists:
			ErrorResponse(w, r, http.StatusConflict, "a user with this email or username already exists")
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	// Return a DTO that doesn't include sensitive info like password hash.
	userResponse := dto.NewUserResponse(createdUser)
	encode(w, r, http.StatusCreated, userResponse)
}

// @Summary      Log in a user
// @Description  Authenticates a user with email and password, returning a JWT and user details upon success.
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        credentials body dto.LoginUserRequest true "User login credentials"
// @Success      200  {object}  dto.LoginUserResponse "Successfully logged in"
// @Failure      400  {object}  map[string]string "Bad request (e.g., malformed JSON)"
// @Failure      401  {object}  map[string]string "Unauthorized (invalid credentials)"
// @Failure      422  {object}  map[string]string "Unprocessable entity (validation error)"
// @Failure      500  {object}  map[string]string "Internal server error"
// @Router       /auth/login [post]
// LoginUser handles user login requests.
// It expects an email and password in the request body.
// On successful authentication, it returns a JWT and user information.
func (h *AuthHandler) LoginUser(w http.ResponseWriter, r *http.Request) {
	var input dto.LoginUserRequest

	if !decodeAndValidate(w, r, &input) {
		return // Errors handled by decodeAndValidate
	}

	token, user, err := h.authService.Login(r.Context(), input.Email, input.Password)
	if err != nil {
		switch {
		case errors.Is(err, auth.ErrInvalidCredentials), errors.Is(err, store.ErrNotFound):
			ErrorResponse(w, r, http.StatusUnauthorized, "invalid email or password")
		default:
			ServerErrorResponse(w, r, err)
		}
		return
	}

	loginResponse := dto.LoginUserResponse{
		Token: token,
		User:  dto.NewUserResponse(user),
	}

	encode(w, r, http.StatusOK, loginResponse)
}
