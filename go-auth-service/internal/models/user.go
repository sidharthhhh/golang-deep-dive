package models

import "time"

type User struct {
	ID           int64     `db:"id"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	IsVerified   bool      `db:"is_verified"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}
