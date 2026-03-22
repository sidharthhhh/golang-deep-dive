package service

import (
	"context"
	"time"

	"github.com/sidharthhhh/go-auth-service/internal/repository"
	"github.com/sidharthhhh/go-auth-service/internal/utils"
	apperrors "github.com/sidharthhhh/go-auth-service/internal/pkg/errors"
)

// TokenValidationResult represents the result of token validation
type TokenValidationResult struct {
	Valid       bool      `json:"valid"`
	UserID      int       `json:"user_id,omitempty"`
	Email       string    `json:"email,omitempty"`
	Role        string    `json:"role,omitempty"`
	Permissions []string  `json:"permissions,omitempty"`
	ExpiresAt   time.Time `json:"expires_at,omitempty"`
	IssuedAt    time.Time `json:"issued_at,omitempty"`
}

// TokenInfo represents detailed token information
type TokenInfo struct {
	JTI       string    `json:"jti"`
	UserID    int       `json:"user_id"`
	Email     string    `json:"email"`
	Role      string    `json:"role"`
	ExpiresAt time.Time `json:"expires_at"`
	IssuedAt  time.Time `json:"issued_at"`
}

// TokenService handles token operations
type TokenService interface {
	ValidateToken(ctx context.Context, token string) (*TokenValidationResult, error)
	GetTokenInfo(ctx context.Context, token string) (*TokenInfo, error)
}

type tokenService struct {
	tokenRepo repository.TokenRepository
	jwtSecret string
}

// NewTokenService creates a new token service
func NewTokenService(tokenRepo repository.TokenRepository, jwtSecret string) TokenService {
	return &tokenService{
		tokenRepo: tokenRepo,
		jwtSecret: jwtSecret,
	}
}

// ValidateToken validates a JWT token and returns user information
func (s *tokenService) ValidateToken(ctx context.Context, tokenString string) (*TokenValidationResult, error) {
	// Parse and validate token
	claims, err := utils.ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		return &TokenValidationResult{Valid: false}, apperrors.ErrTokenInvalid
	}

	// Check if token is blacklisted
	isBlacklisted, err := s.tokenRepo.IsTokenBlacklisted(ctx, claims.ID)
	if err != nil {
		return &TokenValidationResult{Valid: false}, err
	}

	if isBlacklisted {
		return &TokenValidationResult{Valid: false}, apperrors.ErrTokenRevoked
	}

	// Get permissions based on role
	permissions := getPermissionsByRole(claims.Role)

	return &TokenValidationResult{
		Valid:       true,
		UserID:      claims.UserID,
		Email:       claims.Email,
		Role:        claims.Role,
		Permissions: permissions,
		ExpiresAt:   claims.ExpiresAt.Time,
		IssuedAt:    claims.IssuedAt.Time,
	}, nil
}

// GetTokenInfo returns detailed information about a token
func (s *tokenService) GetTokenInfo(ctx context.Context, tokenString string) (*TokenInfo, error) {
	// Parse token
	claims, err := utils.ValidateToken(tokenString, s.jwtSecret)
	if err != nil {
		return nil, apperrors.ErrTokenInvalid
	}

	// Check if token is blacklisted
	isBlacklisted, err := s.tokenRepo.IsTokenBlacklisted(ctx, claims.ID)
	if err != nil {
		return nil, err
	}

	if isBlacklisted {
		return nil, apperrors.ErrTokenRevoked
	}

	return &TokenInfo{
		JTI:       claims.ID,
		UserID:    claims.UserID,
		Email:     claims.Email,
		Role:      claims.Role,
		ExpiresAt: claims.ExpiresAt.Time,
		IssuedAt:  claims.IssuedAt.Time,
	}, nil
}

// getPermissionsByRole returns permissions based on user role
func getPermissionsByRole(role string) []string {
	switch role {
	case "super_admin":
		return []string{
			"users:read", "users:write", "users:delete",
			"admin:read", "admin:write",
			"system:read", "system:write",
		}
	case "admin":
		return []string{
			"users:read", "users:write",
			"admin:read",
		}
	case "user":
		return []string{
			"profile:read", "profile:write",
		}
	default:
		return []string{}
	}
}
