package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/pkg/response"
	"github.com/sidharthhhh/go-auth-service/internal/service"
	"github.com/sidharthhhh/go-auth-service/internal/utils"
	"go.uber.org/zap"
)

// TokenHandler handles token-related requests
type TokenHandler struct {
	tokenService service.TokenService
	logger       *zap.Logger
}

// NewTokenHandler creates a new token handler
func NewTokenHandler(tokenService service.TokenService, logger *zap.Logger) *TokenHandler {
	return &TokenHandler{
		tokenService: tokenService,
		logger:       logger,
	}
}

// ValidateTokenRequest represents the token validation request
type ValidateTokenRequest struct {
	Token string `json:"token" binding:"required"`
}

// TokenInfoResponse represents the token info response
type TokenInfoResponse struct {
	JTI       string   `json:"jti"`
	UserID    int      `json:"user_id"`
	Email     string   `json:"email"`
	Role      string   `json:"role"`
	ExpiresAt string   `json:"expires_at"`
	IssuedAt  string   `json:"issued_at"`
}

// ValidateToken validates a JWT token
// @Summary Validate JWT token
// @Description Validates a JWT token and returns user information
// @Tags Token
// @Accept json
// @Produce json
// @Param request body ValidateTokenRequest true "Token to validate"
// @Success 200 {object} response.Response{data=service.TokenValidationResult}
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /v1/auth/validate [post]
func (h *TokenHandler) ValidateToken(c *gin.Context) {
	var req ValidateTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.tokenService.ValidateToken(c.Request.Context(), req.Token)
	if err != nil {
		h.logger.Error("token validation failed",
			zap.Error(err),
			zap.String("request_id", c.GetString("request_id")),
		)
		response.Unauthorized(c, "Invalid or expired token")
		return
	}

	if !result.Valid {
		response.Unauthorized(c, "Token is not valid")
		return
	}

	response.Success(c, http.StatusOK, "Token is valid", result)
}

// GetTokenInfo returns detailed information about the current token
// @Summary Get token information
// @Description Returns detailed information about the authenticated token
// @Tags Token
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response{data=TokenInfoResponse}
// @Failure 401 {object} response.Response
// @Router /v1/auth/token-info [get]
func (h *TokenHandler) GetTokenInfo(c *gin.Context) {
	// Get token from Authorization header
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" || len(authHeader) < 8 {
		response.Unauthorized(c, "No token provided")
		return
	}

	tokenString := authHeader[7:] // Remove "Bearer " prefix

	tokenInfo, err := h.tokenService.GetTokenInfo(c.Request.Context(), tokenString)
	if err != nil {
		h.logger.Error("failed to get token info",
			zap.Error(err),
			zap.String("request_id", c.GetString("request_id")),
		)
		response.Unauthorized(c, "Invalid token")
		return
	}

	resp := TokenInfoResponse{
		JTI:       tokenInfo.JTI,
		UserID:    tokenInfo.UserID,
		Email:     tokenInfo.Email,
		Role:      tokenInfo.Role,
		ExpiresAt: tokenInfo.ExpiresAt.Format("2006-01-02T15:04:05Z07:00"),
		IssuedAt:  tokenInfo.IssuedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response.Success(c, http.StatusOK, "Token information retrieved", resp)
}

// RefreshToken generates a new token for the authenticated user
// @Summary Refresh JWT token
// @Description Generates a new JWT token for the authenticated user
// @Tags Token
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /v1/auth/refresh [post]
func (h *TokenHandler) RefreshToken(c *gin.Context) {
	// Get user info from context (set by auth middleware)
	userID := c.GetInt("user_id")
	email := c.GetString("email")
	role := c.GetString("role")
	jwtSecret := c.GetString("jwt_secret")

	// Generate new token
	token, err := utils.GenerateToken(userID, email, role, jwtSecret)
	if err != nil {
		h.logger.Error("failed to generate token",
			zap.Error(err),
			zap.Int("user_id", userID),
			zap.String("request_id", c.GetString("request_id")),
		)
		response.InternalError(c, "Failed to generate token")
		return
	}

	expiryInfo := utils.GetTokenExpiryInfo(role)

	data := gin.H{
		"token":      token,
		"email":      email,
		"role":       role,
		"expires_in": expiryInfo,
	}

	response.Success(c, http.StatusOK, "Token refreshed successfully", data)
}
