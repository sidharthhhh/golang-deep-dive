package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/sidharthhhh/go-auth-service/internal/models"
)

type UserRepository interface {
	CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error)
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) CreateUser(ctx context.Context, email, passwordHash string) (*models.User, error) {
	query := `INSERT INTO users (email, password_hash, is_verified, created_at, updated_at) 
	          VALUES (?, ?, ?, ?, ?)`

	now := time.Now()
	result, err := r.db.ExecContext(ctx, query, email, passwordHash, false, now, now)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &models.User{
		ID:           int(id),
		Email:        email,
		PasswordHash: passwordHash,
		IsVerified:   false,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

func (r *userRepository) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `SELECT id, email, password_hash, is_verified, created_at, updated_at 
	          FROM users WHERE email = ?`

	user := &models.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.PasswordHash,
		&user.IsVerified,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}
