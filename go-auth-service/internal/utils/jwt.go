package utils

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("token has expired")
)

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT token for a user with role-based expiry
// - user: 7 days
// - admin: 7 days
// - super_admin: 30 days (longer for convenience)
func GenerateToken(userID int, email, role, secret string) (string, error) {
	// Determine token expiry based on role
	var expiryDuration time.Duration
	switch role {
	case "user":
		expiryDuration = 7 * 24 * time.Hour // 7 days for users
	case "admin":
		expiryDuration = 7 * 24 * time.Hour // 7 days for admins
	case "super_admin":
		expiryDuration = 30 * 24 * time.Hour // 30 days for super admins
	default:
		expiryDuration = 7 * 24 * time.Hour // default to 7 days
	}

	// Generate unique JTI (JWT ID) for token tracking
	jti, err := generateJTI()
	if err != nil {
		return "", err
	}

	claims := Claims{
		UserID: userID,
		Email:  email,
		Role:   role,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        jti,
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expiryDuration)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// generateJTI generates a unique JWT ID
func generateJTI() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// ValidateToken verifies and parses a JWT token
func ValidateToken(tokenString, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrInvalidToken
}


// GetTokenExpiry returns the expiry duration for a given role
func GetTokenExpiry(role string) time.Duration {
	switch role {
	case "user":
		return 7 * 24 * time.Hour // 7 days for users
	case "admin":
		return 7 * 24 * time.Hour // 7 days for admins
	case "super_admin":
		return 30 * 24 * time.Hour // 30 days for super admins
	default:
		return 7 * 24 * time.Hour
	}
}

// GetTokenExpiryInfo returns human-readable expiry information
func GetTokenExpiryInfo(role string) string {
	switch role {
	case "user":
		return "7 days"
	case "admin":
		return "7 days"
	case "super_admin":
		return "30 days"
	default:
		return "7 days"
	}
}
