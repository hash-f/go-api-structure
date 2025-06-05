package user

import (
	"context"
	"go-api-structure/internal/store"
	"go-api-structure/internal/store/db" // For db.User type

	"github.com/google/uuid"
)

// ServiceInterface defines the operations for the user service.
type ServiceInterface interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*db.User, error)
	GetUserByAPIKey(ctx context.Context, apiKey string) (*db.User, error)
	// Add other user-specific business logic methods here if needed
}

// Service provides user-related operations.
type Service struct {
	userStore store.UserStore
}

// NewService creates a new UserService.
func NewService(userStore store.UserStore) *Service {
	return &Service{
		userStore: userStore,
	}
}

// GetUserByID retrieves a user by their ID.
func (s *Service) GetUserByID(ctx context.Context, id uuid.UUID) (*db.User, error) {
	user, err := s.userStore.GetUserByID(ctx, id)
	if err != nil {
		return nil, err // Error handling (e.g., store.ErrNotFound) is done in the store layer
	}
	return &user, nil
}

// GetUserByAPIKey retrieves a user by their API key.
func (s *Service) GetUserByAPIKey(ctx context.Context, apiKey string) (*db.User, error) {
	user, err := s.userStore.GetUserByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, err // Error handling (e.g., store.ErrNotFound) is done in the store layer
	}
	return &user, nil
}
