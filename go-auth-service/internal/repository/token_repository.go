package repository

import (
	"context"
	"database/sql"
	"time"
)

type TokenRepository interface {
	BlacklistToken(ctx context.Context, jti string, userID int, expiresAt time.Time) error
	IsTokenBlacklisted(ctx context.Context, jti string) (bool, error)
	CleanupExpiredTokens(ctx context.Context) error
	BlacklistAllUserTokens(ctx context.Context, userID int) error
}

type tokenRepository struct {
	db *sql.DB
}

func NewTokenRepository(db *sql.DB) TokenRepository {
	return &tokenRepository{db: db}
}

func (r *tokenRepository) BlacklistToken(ctx context.Context, jti string, userID int, expiresAt time.Time) error {
	query := `INSERT INTO token_blacklist (token_jti, user_id, expires_at) VALUES (?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, jti, userID, expiresAt)
	return err
}

func (r *tokenRepository) IsTokenBlacklisted(ctx context.Context, jti string) (bool, error) {
	query := `SELECT COUNT(*) FROM token_blacklist WHERE token_jti = ? AND expires_at > NOW()`
	var count int
	err := r.db.QueryRowContext(ctx, query, jti).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *tokenRepository) CleanupExpiredTokens(ctx context.Context) error {
	query := `DELETE FROM token_blacklist WHERE expires_at <= NOW()`
	_, err := r.db.ExecContext(ctx, query)
	return err
}


func (r *tokenRepository) BlacklistAllUserTokens(ctx context.Context, userID int) error {
	// Instead of blacklisting individual tokens, we update the user's password_changed_at timestamp
	// The auth middleware will check if tokens were issued before this timestamp
	query := `UPDATE users SET password_changed_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}
