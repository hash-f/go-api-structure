package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"go-api-structure/internal/api/dto" // Assuming CreateUserRequest is here
	"go-api-structure/internal/store"
	"go-api-structure/internal/store/db" // sqlc generated models and params
)

var (
	ErrUserAlreadyExists = errors.New("user with this email or username already exists")
	ErrInvalidCredentials  = errors.New("invalid email or password")
)

// AuthService provides methods for user authentication and registration.
type AuthService struct {
	userStore   store.UserStore
	jwtSecret   string
	tokenExpiry time.Duration
}

// NewAuthService creates a new AuthService.
func NewAuthService(userStore store.UserStore, jwtSecret string, tokenExpiry time.Duration) *AuthService {
	return &AuthService{
		userStore:   userStore,
		jwtSecret:   jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

// Register creates a new user after validating input and hashing the password.
func (s *AuthService) Register(ctx context.Context, req *dto.CreateUserRequest) (*db.User, error) {
	hashedPassword, err := HashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("failed to hash password during registration: %w", err)
	}

	params := db.CreateUserParams{
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: hashedPassword,
	}

	user, err := s.userStore.CreateUser(ctx, params)
	if err != nil {
		// TODO: Check for specific database errors, e.g., unique constraint violation from pgx,
		// and wrap them into ErrUserAlreadyExists if applicable.
		return nil, fmt.Errorf("failed to create user: %w", err)
	}

	return &user, nil
}

// Login authenticates a user by email and password, returning a JWT if successful.
func (s *AuthService) Login(ctx context.Context, email, password string) (string, *db.User, error) {
	user, err := s.userStore.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrNotFound) {
			return "", nil, ErrInvalidCredentials
		}
		return "", nil, fmt.Errorf("failed to get user by email: %w", err)
	}

	if !CheckPasswordHash(password, user.PasswordHash) {
		return "", nil, ErrInvalidCredentials
	}

	// Create JWT claims
	claims := jwt.RegisteredClaims{
		Subject:   user.ID.String(), // Use user's UUID as subject
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.tokenExpiry)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		Issuer:    "go-api-structure", // Optional: identify the issuer
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", nil, fmt.Errorf("failed to sign token: %w", err)
	}

	return signedToken, &user, nil
}
