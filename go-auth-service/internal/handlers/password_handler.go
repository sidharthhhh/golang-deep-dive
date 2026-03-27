package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/pkg/response"
	"github.com/sidharthhhh/go-auth-service/internal/service"
	"go.uber.org/zap"
)

// PasswordHandler handles password-related requests
type PasswordHandler struct {
	passwordService service.PasswordService
	logger          *zap.Logger
}

// NewPasswordHandler creates a new password handler
func NewPasswordHandler(passwordService service.PasswordService, logger *zap.Logger) *PasswordHandler {
	return &PasswordHandler{
		passwordService: passwordService,
		logger:          logger,
	}
}

// ForgotPasswordRequest represents the forgot password request
type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

// ResetPasswordRequest represents the reset password request
type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	Email       string `json:"email" binding:"required,email"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ChangePasswordRequest represents the change password request
type ChangePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

// ForgotPassword initiates the password reset process
// @Summary Forgot password
// @Description Sends a password reset token to the user's email
// @Tags Password
// @Accept json
// @Produce json
// @Param request body ForgotPasswordRequest true "Email address"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /v1/auth/forgot-password [post]
func (h *PasswordHandler) ForgotPassword(c *gin.Context) {
	var req ForgotPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	resetToken, err := h.passwordService.InitiatePasswordReset(c.Request.Context(), req.Email)
	if err != nil {
		h.logger.Error("failed to initiate password reset",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("request_id", c.GetString("request_id")),
		)
		// Don't reveal if user exists or not for security
		response.Success(c, http.StatusOK, "If the email exists, a password reset link has been sent", nil)
		return
	}

	// If resetToken is empty, user doesn't exist (but don't reveal this)
	if resetToken == "" {
		response.Success(c, http.StatusOK, "If the email exists, a password reset link has been sent", nil)
		return
	}

	// In production, send email with reset link
	// For development/testing, return the token
	data := gin.H{
		"message": "Password reset initiated",
	}

	// Include token in development mode (check environment)
	appEnv := c.GetString("APP_ENV")
	if appEnv == "" {
		// If APP_ENV not set in context, check from config
		appEnv = "development" // Default to development
	}
	
	if appEnv == "development" || appEnv == "" {
		data["reset_token"] = resetToken
		data["note"] = "In production, this token would be sent via email"
	}

	response.Success(c, http.StatusOK, "Password reset email sent", data)
}

// ResetPassword resets the user's password using a reset token
// @Summary Reset password
// @Description Resets user password using the reset token
// @Tags Password
// @Accept json
// @Produce json
// @Param request body ResetPasswordRequest true "Reset token and new password"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /v1/auth/reset-password [post]
func (h *PasswordHandler) ResetPassword(c *gin.Context) {
	var req ResetPasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	err := h.passwordService.ResetPassword(c.Request.Context(), req.Token, req.Email, req.NewPassword)
	if err != nil {
		h.logger.Error("failed to reset password",
			zap.Error(err),
			zap.String("email", req.Email),
			zap.String("request_id", c.GetString("request_id")),
		)
		response.Unauthorized(c, "Invalid or expired reset token")
		return
	}

	response.Success(c, http.StatusOK, "Password reset successfully", nil)
}

// ChangePassword changes the authenticated user's password
// @Summary Change password
// @Description Changes the password for the authenticated user
// @Tags Password
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body ChangePasswordRequest true "Old and new password"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Router /v1/auth/change-password [post]
func (h *PasswordHandler) ChangePassword(c *gin.Context) {
	userID := c.GetInt("user_id")
	if userID == 0 {
		response.Unauthorized(c, "User not authenticated")
		return
	}

	var req ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	err := h.passwordService.ChangePassword(c.Request.Context(), int64(userID), req.OldPassword, req.NewPassword)
	if err != nil {
		h.logger.Error("failed to change password",
			zap.Error(err),
			zap.Int("user_id", userID),
			zap.String("request_id", c.GetString("request_id")),
		)
		response.Unauthorized(c, "Invalid old password")
		return
	}

	response.Success(c, http.StatusOK, "Password changed successfully", nil)
}
