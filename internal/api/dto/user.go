package dto

import (
	"net/mail"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"go-api-structure/internal/store/db" // For NewUserResponse
)

// CreateUserRequest defines the expected structure for a new user registration request.
// It includes fields for username, email, and password.
// It implements the Validator interface (defined in the api package) for self-validation.
type CreateUserRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Valid checks the validity of the CreateUserRequest fields.
// It returns a map of validation errors if any are found, otherwise nil.
func (r *CreateUserRequest) Valid() map[string]string {
	errors := make(map[string]string)

	// Username validation
	if strings.TrimSpace(r.Username) == "" {
		errors["username"] = "username must be provided"
	} else if utf8.RuneCountInString(r.Username) < 3 {
		errors["username"] = "username must be at least 3 characters long"
	} else if utf8.RuneCountInString(r.Username) > 50 {
		errors["username"] = "username must not be more than 50 characters long"
	}

	// Email validation
	if strings.TrimSpace(r.Email) == "" {
		errors["email"] = "email must be provided"
	} else {
		_, err := mail.ParseAddress(r.Email)
		if err != nil {
			errors["email"] = "email must be a valid email address"
		}
	}

	// Password validation
	if strings.TrimSpace(r.Password) == "" {
		errors["password"] = "password must be provided"
	} else if utf8.RuneCountInString(r.Password) < 8 {
		errors["password"] = "password must be at least 8 characters long"
	} else if utf8.RuneCountInString(r.Password) > 72 { // bcrypt has a 72-byte limit
		errors["password"] = "password must not be more than 72 characters long"
	}

	if len(errors) == 0 {
		return nil
	}
	return errors
}

// UserResponse defines the structure for user data returned by the API.
// It omits sensitive information like the password hash.
// Note: pgtype.Timestamptz from db.User needs to be converted to time.Time for JSON marshaling if not handled by `json:"created_at"` in db.User itself.
// sqlc's default JSON tags for pgtype.Timestamptz handle this correctly.
// If we needed custom formatting, we'd handle it here or in a custom MarshalJSON.
type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// NewUserResponse creates a new UserResponse DTO from a db.User model.
func NewUserResponse(user *db.User) *UserResponse {
	if user == nil {
		return nil
	}
	return &UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Time, // Convert pgtype.Timestamptz to time.Time
		UpdatedAt: user.UpdatedAt.Time, // Convert pgtype.Timestamptz to time.Time
	}
}

// LoginUserRequest defines the structure for a user login request.
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=72"`
}

// Valid checks if the LoginUserRequest fields are valid.
func (r *LoginUserRequest) Valid() map[string]string {
	errors := make(map[string]string)
	if r.Email == "" {
		errors["email"] = "Email is required"
	} else if !isValidEmail(r.Email) {
		errors["email"] = "Invalid email format"
	}

	if r.Password == "" {
		errors["password"] = "Password is required"
	} else if len(r.Password) < 8 {
		errors["password"] = "Password must be at least 8 characters long"
	} else if len(r.Password) > 72 { // bcrypt has a 72-byte limit
		errors["password"] = "Password must be at most 72 characters long"
	}

	if len(errors) > 0 {
		return errors
	}
	return nil
}

// LoginUserResponse defines the structure for a successful login response.
type LoginUserResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}

func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
