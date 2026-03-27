package service

import (
	"context"
	"database/sql"
	"time"

	apperrors "github.com/sidharthhhh/go-auth-service/internal/pkg/errors"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
	"github.com/sidharthhhh/go-auth-service/internal/utils"
)

// PasswordService handles password-related operations
type PasswordService interface {
	InitiatePasswordReset(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, email, newPassword string) error
	ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error
}

type passwordService struct {
	userRepo     repository.UserRepository
	passwordRepo repository.PasswordResetRepository
	tokenRepo    repository.TokenRepository
}

// NewPasswordService creates a new password service
func NewPasswordService(
	userRepo repository.UserRepository,
	passwordRepo repository.PasswordResetRepository,
	tokenRepo repository.TokenRepository,
) PasswordService {
	return &passwordService{
		userRepo:     userRepo,
		passwordRepo: passwordRepo,
		tokenRepo:    tokenRepo,
	}
}

// InitiatePasswordReset creates a password reset token
func (s *passwordService) InitiatePasswordReset(ctx context.Context, email string) (string, error) {
	// Find user by email
	user, err := s.userRepo.FindUserByEmail(ctx, email)
	if err != nil {
		if err == sql.ErrNoRows {
			// Don't reveal if user exists - return empty token
			return "", nil
		}
		// Log actual error but return empty token for security
		return "", err
	}

	// Generate reset token
	resetToken, err := utils.GeneratePasswordResetToken()
	if err != nil {
		return "", err
	}

	// Store reset token (expires in 1 hour)
	expiresAt := time.Now().Add(1 * time.Hour)
	err = s.passwordRepo.CreatePasswordResetToken(ctx, user.ID, resetToken, expiresAt)
	if err != nil {
		// Return error so handler can log it
		return "", err
	}

	// TODO: Send email with reset link
	// emailService.SendPasswordResetEmail(user.Email, resetToken)

	return resetToken, nil
}

// ResetPassword resets the password using a reset token
func (s *passwordService) ResetPassword(ctx context.Context, token, email, newPassword string) error {
	// Validate reset token
	userID, err := s.passwordRepo.ValidatePasswordResetToken(ctx, token)
	if err != nil {
		return apperrors.ErrTokenInvalid
	}

	// Verify the email matches the user who requested the reset
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return apperrors.ErrTokenInvalid
	}
	
	if user.Email != email {
		return apperrors.ErrTokenInvalid
	}

	// Hash new password
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update user password
	err = s.userRepo.UpdatePassword(ctx, userID, passwordHash)
	if err != nil {
		return err
	}

	// Invalidate reset token
	err = s.passwordRepo.DeletePasswordResetToken(ctx, token)
	if err != nil {
		return err
	}

	// Invalidate all user sessions (logout from all devices)
	err = s.tokenRepo.BlacklistAllUserTokens(ctx, int(userID))
	if err != nil {
		// Log error but don't fail the password reset
		return nil
	}

	return nil
}

// ChangePassword changes the user's password
func (s *passwordService) ChangePassword(ctx context.Context, userID int64, oldPassword, newPassword string) error {
	// Get user
	user, err := s.userRepo.GetUserByID(ctx, userID)
	if err != nil {
		return apperrors.ErrNotFound
	}

	// Verify old password
	err = utils.VerifyPassword(user.PasswordHash, oldPassword)
	if err != nil {
		return apperrors.ErrInvalidPassword
	}

	// Hash new password
	passwordHash, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// Update password
	err = s.userRepo.UpdatePassword(ctx, userID, passwordHash)
	if err != nil {
		return err
	}

	// Invalidate all user sessions (logout from all devices)
	err = s.tokenRepo.BlacklistAllUserTokens(ctx, int(userID))
	if err != nil {
		// Log error but don't fail the password change
		return nil
	}

	return nil
}
