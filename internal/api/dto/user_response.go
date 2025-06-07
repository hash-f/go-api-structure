package dto

import (
	"go-api-structure/internal/store/db"
	"time"

	"github.com/google/uuid"
)

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
