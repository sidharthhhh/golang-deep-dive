package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
	"github.com/sidharthhhh/go-auth-service/internal/utils"
)

func AuthMiddleware(jwtSecret string, tokenRepo repository.TokenRepository, userRepo repository.UserRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := parts[1]
		claims, err := utils.ValidateToken(token, jwtSecret)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid or expired token"})
			c.Abort()
			return
		}

		// Check if token is blacklisted (logged out)
		isBlacklisted, err := tokenRepo.IsTokenBlacklisted(c.Request.Context(), claims.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to validate token"})
			c.Abort()
			return
		}
		if isBlacklisted {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token has been revoked"})
			c.Abort()
			return
		}

		// Check if password was changed after token was issued
		user, err := userRepo.GetUserByID(c.Request.Context(), int64(claims.UserID))
		if err == nil && user.PasswordChangedAt != nil {
			// If token was issued before password change, invalidate it
			if claims.IssuedAt.Time.Before(*user.PasswordChangedAt) {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "token invalidated due to password change"})
				c.Abort()
				return
			}
		}

		// Set user info in context
		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Set("role", claims.Role)

		c.Next()
	}
}


// RequireRole middleware checks if user has required role
func RequireRole(roles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userRole := c.GetString("role")
		if userRole == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			c.Abort()
			return
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range roles {
			if userRole == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{"error": "insufficient permissions"})
			c.Abort()
			return
		}

		c.Next()
	}
}
