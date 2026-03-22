package service

import (
	"context"
	"database/sql"

	"github.com/sidharthhhh/go-auth-service/internal/models"
	apperrors "github.com/sidharthhhh/go-auth-service/internal/pkg/errors"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
)

// UserUpdate represents user update data
type UserUpdate struct {
	Email      *string
	IsVerified *bool
}

// UserService handles user management operations
type UserService interface {
	GetAllUsers(ctx context.Context, page, limit int) ([]*models.User, int, error)
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	UpdateUser(ctx context.Context, userID int64, updates *UserUpdate) error
	DeleteUser(ctx context.Context, userID int64) error
	BanUser(ctx context.Context, userID int64, reason string) error
	UnbanUser(ctx context.Context, userID int64) error
}

type userService struct {
	userRepo repository.UserRepository
}

// NewUserService creates a new user service
func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

// GetAllUsers returns paginated list of users
func (s *userService) GetAllUsers(ctx context.Context, page, limit int) ([]*models.User, int, error) {
	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit
	users, total, err := s.userRepo.GetAllUsers(ctx, limit, offset)
	if err != nil {
		return nil, 0, err
	}

	return users, total, nil
}

// GetUserByID returns a user by ID
func (s *userService) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, apperrors.ErrNotFound
		}
		return nil, err
	}

	return user, nil
}

// UpdateUser updates user information
func (s *userService) UpdateUser(ctx context.Context, userID int64, updates *UserUpdate) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperrors.ErrNotFound
		}
		return err
	}

	// Update user
	return s.userRepo.UpdateUser(ctx, userID, updates)
}

// DeleteUser deletes a user
func (s *userService) DeleteUser(ctx context.Context, userID int64) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperrors.ErrNotFound
		}
		return err
	}

	return s.userRepo.DeleteUser(ctx, userID)
}

// BanUser bans a user
func (s *userService) BanUser(ctx context.Context, userID int64, reason string) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperrors.ErrNotFound
		}
		return err
	}

	// TODO: Implement ban logic (add banned field to user model)
	// For now, we'll just mark user as not verified
	return s.userRepo.UpdateUser(ctx, userID, &UserUpdate{
		IsVerified: boolPtr(false),
	})
}

// UnbanUser unbans a user
func (s *userService) UnbanUser(ctx context.Context, userID int64) error {
	// Check if user exists
	_, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return apperrors.ErrNotFound
		}
		return err
	}

	// TODO: Implement unban logic
	return s.userRepo.UpdateUser(ctx, userID, &UserUpdate{
		IsVerified: boolPtr(true),
	})
}

func boolPtr(b bool) *bool {
	return &b
}
