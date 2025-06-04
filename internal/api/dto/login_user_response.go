package dto

// LoginUserResponse defines the structure for a successful login response.
type LoginUserResponse struct {
	Token string        `json:"token"`
	User  *UserResponse `json:"user"`
}
