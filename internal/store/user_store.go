package store

import (
	"context"
	"errors"
	"go-api-structure/internal/store/db"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

// UserStore defines the interface for user-specific data operations.
// It will typically be implemented by a struct that has access to a *db.Queries object.
// For now, we'll list methods that correspond to our sqlc queries for users.
// We'll also need to consider how parameters are passed (e.g., DTOs vs. direct model types).
type UserStore interface {
	CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error)
	GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error)
	GetUserByEmail(ctx context.Context, email string) (db.User, error)
	GetUserByUsername(ctx context.Context, username string) (db.User, error)
	GetUserByAPIKey(ctx context.Context, apiKey string) (db.User, error)
	// TODO: Add UpdateUser, DeleteUser if needed later
}

// UserStore implementation
func (s *SQLStore) CreateUser(ctx context.Context, arg db.CreateUserParams) (db.User, error) {
	return s.Queries.CreateUser(ctx, arg)
}

func (s *SQLStore) GetUserByID(ctx context.Context, id uuid.UUID) (db.User, error) {
	user, err := s.Queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}

func (s *SQLStore) GetUserByEmail(ctx context.Context, email string) (db.User, error) {
	user, err := s.Queries.GetUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}

func (s *SQLStore) GetUserByUsername(ctx context.Context, username string) (db.User, error) {
	user, err := s.Queries.GetUserByUsername(ctx, username)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}

func (s *SQLStore) GetUserByAPIKey(ctx context.Context, apiKey string) (db.User, error) {
	user, err := s.Queries.GetUserByAPIKey(ctx, apiKey)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return db.User{}, ErrNotFound
		}
		return db.User{}, err
	}
	return user, nil
}
