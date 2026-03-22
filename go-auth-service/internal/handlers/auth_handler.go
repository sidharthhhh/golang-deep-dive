package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/service"
	"github.com/sidharthhhh/go-auth-service/internal/utils"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
	}
}

type RegisterRequest struct {
	Email          string  `json:"email" binding:"required,email"`
	Password       string  `json:"password" binding:"required,min=8"`
	SuperAdminCode *string `json:"super_admin_code,omitempty"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type RegisterResponse struct {
	ID         int    `json:"id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
}

type LoginResponse struct {
	Token string `json:"token"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

type PromoteToAdminRequest struct {
	UserID int `json:"user_id" binding:"required"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authService.RegisterUser(c.Request.Context(), req.Email, req.Password, req.SuperAdminCode)
	if err != nil {
		switch err {
		case service.ErrInvalidEmail:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid email format"})
		case service.ErrUserExists:
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		ID:         int(user.ID),
		Email:      user.Email,
		Role:       string(user.Role),
		IsVerified: user.IsVerified,
	})
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.authService.LoginUser(c.Request.Context(), req.Email, req.Password)
	if err != nil {
		switch err {
		case service.ErrUserNotFound:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		case service.ErrInvalidPassword:
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token":      result.Token,
		"email":      result.User.Email,
		"role":       string(result.User.Role),
		"expires_in": result.ExpiresIn,
		"message":    "Login successful",
	})
}


func (h *AuthHandler) PromoteToAdmin(c *gin.Context) {
	// Get super admin ID from context (set by auth middleware)
	superAdminID := c.GetInt("user_id")
	role := c.GetString("role")

	// Verify caller is super admin
	if role != "super_admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "only super admins can promote users"})
		return
	}

	var req PromoteToAdminRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.authService.PromoteToAdmin(c.Request.Context(), superAdminID, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to promote user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user promoted to admin successfully"})
}


func (h *AuthHandler) RefreshToken(c *gin.Context) {
	// Get user info from context (set by auth middleware)
	userID := c.GetInt("user_id")
	email := c.GetString("email")
	role := c.GetString("role")

	// Generate new token
	token, err := utils.GenerateToken(userID, email, role, c.GetString("jwt_secret"))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	expiryInfo := utils.GetTokenExpiryInfo(role)

	c.JSON(http.StatusOK, gin.H{
		"token":      token,
		"email":      email,
		"role":       role,
		"expires_in": expiryInfo,
		"message":    "Token refreshed successfully",
	})
}


func (h *AuthHandler) Logout(c *gin.Context) {
	// Get token info from context
	userID := c.GetInt("user_id")
	
	// Get the token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "no token provided"})
		return
	}

	// Extract token (remove "Bearer " prefix)
	tokenString := authHeader[7:]
	
	// Parse token to get JTI and expiry
	jwtSecret := c.GetString("jwt_secret")
	claims, err := utils.ValidateToken(tokenString, jwtSecret)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	// Blacklist the token
	err = h.authService.LogoutUser(c.Request.Context(), claims.ID, userID, claims.ExpiresAt.Time)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to logout"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Logged out successfully",
	})
}
