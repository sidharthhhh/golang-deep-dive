package repository

import (
	"context"
	"database/sql"
	"time"
)

// PasswordResetRepository handles password reset token operations
type PasswordResetRepository interface {
	CreatePasswordResetToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error
	ValidatePasswordResetToken(ctx context.Context, token string) (int64, error)
	DeletePasswordResetToken(ctx context.Context, token string) error
	CleanupExpiredTokens(ctx context.Context) error
}

type passwordResetRepository struct {
	db *sql.DB
}

// NewPasswordResetRepository creates a new password reset repository
func NewPasswordResetRepository(db *sql.DB) PasswordResetRepository {
	return &passwordResetRepository{db: db}
}

// CreatePasswordResetToken creates a new password reset token
func (r *passwordResetRepository) CreatePasswordResetToken(ctx context.Context, userID int64, token string, expiresAt time.Time) error {
	query := `INSERT INTO password_reset_tokens (user_id, token, expires_at, created_at) 
	          VALUES (?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, userID, token, expiresAt, time.Now())
	return err
}

// ValidatePasswordResetToken validates a reset token and returns the user ID
func (r *passwordResetRepository) ValidatePasswordResetToken(ctx context.Context, token string) (int64, error) {
	query := `SELECT user_id FROM password_reset_tokens 
	          WHERE token = ? AND expires_at > NOW() AND used = FALSE`

	var userID int64
	err := r.db.QueryRowContext(ctx, query, token).Scan(&userID)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

// DeletePasswordResetToken deletes a reset token
func (r *passwordResetRepository) DeletePasswordResetToken(ctx context.Context, token string) error {
	query := `UPDATE password_reset_tokens SET used = TRUE WHERE token = ?`
	_, err := r.db.ExecContext(ctx, query, token)
	return err
}

// CleanupExpiredTokens removes expired reset tokens
func (r *passwordResetRepository) CleanupExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM password_reset_tokens WHERE expires_at <= NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}
