package service

import (
	"context"
	"database/sql"
	"errors"
	"regexp"

	"github.com/sidharthhhh/go-auth-service/internal/models"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
	"github.com/sidharthhhh/go-auth-service/internal/utils"
)

var (
	ErrInvalidEmail    = errors.New("invalid email format")
	ErrUserExists      = errors.New("user already exists")
	ErrInvalidPassword = errors.New("invalid password")
	ErrUserNotFound    = errors.New("user not found")
)

type AuthService interface {
	RegisterUser(ctx context.Context, email, password string) (*models.User, error)
	LoginUser(ctx context.Context, email, password string) (string, error)
}

type authService struct {
	userRepo  repository.UserRepository
	jwtSecret string
}

func NewAuthService(userRepo repository.UserRepository, jwtSecret string) AuthService {
	return &authService{
		userRepo:  userRepo,
		jwtSecret: jwtSecret,
	}
}

func (s *authService) RegisterUser(ctx context.Context, email, password string) (*models.User, error) {
	// Validate email format
	if !isValidEmail(email) {
		return nil, ErrInvalidEmail
	}

	// Check if user already exists
	existingUser, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrUserExists
	}

	// Hash password
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user, err := s.userRepo.CreateUser(ctx, email, passwordHash)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) LoginUser(ctx context.Context, email, password string) (string, error) {
	// Placeholder for login implementation
	return "", errors.New("not implemented")
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
