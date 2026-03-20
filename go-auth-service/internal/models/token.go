package models

import "time"

type VerificationToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	CreatedAt time.Time
}

type RefreshToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	Revoked   bool
	CreatedAt time.Time
}

type PasswordResetToken struct {
	ID        int64
	UserID    int64
	Token     string
	ExpiresAt time.Time
	Used      bool
	CreatedAt time.Time
}
