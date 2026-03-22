package service

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"time"

	"github.com/sidharthhhh/go-auth-service/internal/models"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
	"github.com/sidharthhhh/go-auth-service/internal/utils"
)

var (
	ErrInvalidEmail      = errors.New("invalid email format")
	ErrUserExists        = errors.New("user already exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidAdminCode  = errors.New("invalid super admin code")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidRole       = errors.New("invalid role")
)

type LoginResult struct {
	Token      string
	User       *models.User
	ExpiresIn  string
}

type AuthService interface {
	RegisterUser(ctx context.Context, email, password string, superAdminCode *string) (*models.User, error)
	LoginUser(ctx context.Context, email, password string) (*LoginResult, error)
	LogoutUser(ctx context.Context, jti string, userID int, expiresAt time.Time) error
	PromoteToAdmin(ctx context.Context, superAdminID int, targetUserID int) error
}

type authService struct {
	userRepo       repository.UserRepository
	tokenRepo      repository.TokenRepository
	jwtSecret      string
	superAdminCode string
}

func NewAuthService(userRepo repository.UserRepository, tokenRepo repository.TokenRepository, jwtSecret, superAdminCode string) AuthService {
	return &authService{
		userRepo:       userRepo,
		tokenRepo:      tokenRepo,
		jwtSecret:      jwtSecret,
		superAdminCode: superAdminCode,
	}
}

func (s *authService) RegisterUser(ctx context.Context, email, password string, superAdminCode *string) (*models.User, error) {
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

	// Determine role based on super admin code
	role := models.RoleUser
	if superAdminCode != nil && *superAdminCode == s.superAdminCode {
		role = models.RoleSuperAdmin
	}

	// Hash password
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	// Create user
	user, err := s.userRepo.CreateUser(ctx, email, passwordHash, role)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) LoginUser(ctx context.Context, email, password string) (*LoginResult, error) {
	// Find user by email
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Verify password
	if err := utils.VerifyPassword(user.PasswordHash, password); err != nil {
		return nil, ErrInvalidPassword
	}

	// Generate JWT token with role
	token, err := utils.GenerateToken(int(user.ID), user.Email, string(user.Role), s.jwtSecret)
	if err != nil {
		return nil, err
	}

	// Get expiry info
	expiryInfo := utils.GetTokenExpiryInfo(string(user.Role))

	return &LoginResult{
		Token:     token,
		User:      user,
		ExpiresIn: expiryInfo,
	}, nil
}

func (s *authService) LogoutUser(ctx context.Context, jti string, userID int, expiresAt time.Time) error {
	return s.tokenRepo.BlacklistToken(ctx, jti, userID, expiresAt)
}

func (s *authService) PromoteToAdmin(ctx context.Context, superAdminID int, targetUserID int) error {
	// Update target user role to admin
	return s.userRepo.UpdateUserRole(ctx, targetUserID, models.RoleAdmin)
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
