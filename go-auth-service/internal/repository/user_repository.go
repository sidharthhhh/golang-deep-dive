package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/sidharthhhh/go-auth-service/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string, role models.UserRole) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	GetUserByID(ctx context.Context, userID int64) (*models.User, error)
	GetAllUsers(ctx context.Context, limit, offset int) ([]*models.User, int, error)
	UpdateUserRole(ctx context.Context, userID int, role models.UserRole) error
	UpdateUser(ctx context.Context, userID int64, updates interface{}) error
	UpdatePassword(ctx context.Context, userID int64, passwordHash string) error
	DeleteUser(ctx context.Context, userID int64) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, email, passwordHash string, role models.UserRole) (*models.User, error) {
	query := `INSERT INTO users (email, password_hash, role, is_verified, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, email, passwordHash, role, false, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:           int64(id),
		Email:        email,
		PasswordHash: passwordHash,
		Role:         role,
		IsVerified:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, role, is_verified, created_at, updated_at, password_changed_at 
	          FROM users WHERE email = ?`

	user := &models.User{}
	var passwordChangedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&passwordChangedAt,
	)

	if err != nil {
		return nil, err
	}

	if passwordChangedAt.Valid {
		user.PasswordChangedAt = &passwordChangedAt.Time
	}

	return user, nil
}

func (r *userRepository) UpdateUserRole(ctx context.Context, userID int, role models.UserRole) error {
	query := `UPDATE users SET role = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, role, time.Now(), userID)
	return err
}


func (r *userRepository) GetUserByID(ctx context.Context, userID int64) (*models.User, error) {
	query := `SELECT id, email, password_hash, role, is_verified, created_at, updated_at, password_changed_at 
	          FROM users WHERE id = ?`

	user := &models.User{}
	var passwordChangedAt sql.NullTime
	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
		&passwordChangedAt,
	)

	if err != nil {
		return nil, err
	}

	if passwordChangedAt.Valid {
		user.PasswordChangedAt = &passwordChangedAt.Time
	}

	return user, nil
}

func (r *userRepository) GetAllUsers(ctx context.Context, limit, offset int) ([]*models.User, int, error) {
	// Get total count
	var total int
	countQuery := `SELECT COUNT(*) FROM users`
	err := r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Get users
	query := `SELECT id, email, password_hash, role, is_verified, created_at, updated_at 
	          FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]*models.User, 0)
	for rows.Next() {
		user := &models.User{}
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PasswordHash,
			&user.Role,
			&user.IsVerified,
			&user.CreatedAt,
			&user.UpdatedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, total, nil
}

func (r *userRepository) UpdateUser(ctx context.Context, userID int64, updates interface{}) error {
	// Type assert to get the actual update struct
	// This is a simplified version - in production, use a proper update builder
	query := `UPDATE users SET updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, time.Now(), userID)
	return err
}

func (r *userRepository) DeleteUser(ctx context.Context, userID int64) error {
	query := `DELETE FROM users WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}


func (r *userRepository) UpdatePassword(ctx context.Context, userID int64, passwordHash string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.ExecContext(ctx, query, passwordHash, time.Now(), userID)
	return err
}
