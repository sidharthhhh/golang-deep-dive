package models

import "time"

type UserRole string

const (
	RoleUser       UserRole = "user"
	RoleAdmin      UserRole = "admin"
	RoleSuperAdmin UserRole = "super_admin"
)

type User struct {
	ID                int64      `db:"id"`
	Email             string     `db:"email"`
	PasswordHash      string     `db:"password_hash"`
	Role              UserRole   `db:"role"`
	IsVerified        bool       `db:"is_verified"`
	CreatedAt         time.Time  `db:"created_at"`
	UpdatedAt         time.Time  `db:"updated_at"`
	PasswordChangedAt *time.Time `db:"password_changed_at"`
}
