package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/pkg/response"
	"github.com/sidharthhhh/go-auth-service/internal/service"
	"go.uber.org/zap"
)

// UserHandler handles user management requests
type UserHandler struct {
	userService service.UserService
	logger      *zap.Logger
}

// NewUserHandler creates a new user handler
func NewUserHandler(userService service.UserService, logger *zap.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// UserResponse represents a user in API responses (without sensitive data)
type UserResponse struct {
	ID         int64  `json:"id"`
	Email      string `json:"email"`
	Role       string `json:"role"`
	IsVerified bool   `json:"is_verified"`
	CreatedAt  string `json:"created_at"`
	UpdatedAt  string `json:"updated_at"`
}

// ListUsersResponse represents the paginated users list response
type ListUsersResponse struct {
	Users      []UserResponse `json:"users"`
	Total      int            `json:"total"`
	Page       int            `json:"page"`
	Limit      int            `json:"limit"`
	TotalPages int            `json:"total_pages"`
}

// UpdateUserRequest represents the user update request
type UpdateUserRequest struct {
	Email      *string `json:"email"`
	IsVerified *bool   `json:"is_verified"`
}

// BanUserRequest represents the ban user request
type BanUserRequest struct {
	Reason string `json:"reason" binding:"required"`
}

// ListUsers returns a paginated list of users
// @Summary List all users
// @Description Returns a paginated list of all users (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Success 200 {object} response.Response{data=ListUsersResponse}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Router /v1/admin/users [get]
func (h *UserHandler) ListUsers(c *gin.Context) {
	// Get pagination parameters
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	users, total, err := h.userService.GetAllUsers(c.Request.Context(), page, limit)
	if err != nil {
		h.logger.Error("failed to get users",
			zap.Error(err),
			zap.String("request_id", c.GetString("request_id")),
		)
		response.InternalError(c, "Failed to retrieve users")
		return
	}

	// Convert to response format (exclude sensitive data)
	userResponses := make([]UserResponse, len(users))
	for i, user := range users {
		userResponses[i] = UserResponse{
			ID:         user.ID,
			Email:      user.Email,
			Role:       string(user.Role),
			IsVerified: user.IsVerified,
			CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
			UpdatedAt:  user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
		}
	}

	totalPages := (total + limit - 1) / limit

	data := ListUsersResponse{
		Users:      userResponses,
		Total:      total,
		Page:       page,
		Limit:      limit,
		TotalPages: totalPages,
	}

	response.Success(c, http.StatusOK, "Users retrieved successfully", data)
}

// GetUser returns a specific user by ID
// @Summary Get user by ID
// @Description Returns detailed information about a specific user (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.Response{data=UserResponse}
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /v1/admin/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "Invalid user ID")
		return
	}

	user, err := h.userService.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		h.logger.Error("failed to get user",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.String("request_id", c.GetString("request_id")),
		)
		response.NotFound(c, "User not found")
		return
	}

	userResponse := UserResponse{
		ID:         user.ID,
		Email:      user.Email,
		Role:       string(user.Role),
		IsVerified: user.IsVerified,
		CreatedAt:  user.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:  user.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	response.Success(c, http.StatusOK, "User retrieved successfully", userResponse)
}

// UpdateUser updates a user's information
// @Summary Update user
// @Description Updates user information (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body UpdateUserRequest true "Update data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /v1/admin/users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "Invalid user ID")
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	updates := &service.UserUpdate{
		Email:      req.Email,
		IsVerified: req.IsVerified,
	}

	if err := h.userService.UpdateUser(c.Request.Context(), userID, updates); err != nil {
		h.logger.Error("failed to update user",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.String("request_id", c.GetString("request_id")),
		)
		response.InternalError(c, "Failed to update user")
		return
	}

	response.Success(c, http.StatusOK, "User updated successfully", nil)
}

// DeleteUser deletes a user
// @Summary Delete user
// @Description Deletes a user (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /v1/admin/users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "Invalid user ID")
		return
	}

	if err := h.userService.DeleteUser(c.Request.Context(), userID); err != nil {
		h.logger.Error("failed to delete user",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.String("request_id", c.GetString("request_id")),
		)
		response.InternalError(c, "Failed to delete user")
		return
	}

	response.Success(c, http.StatusOK, "User deleted successfully", nil)
}

// BanUser bans a user
// @Summary Ban user
// @Description Bans a user from the system (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body BanUserRequest true "Ban reason"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /v1/admin/users/{id}/ban [post]
func (h *UserHandler) BanUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "Invalid user ID")
		return
	}

	var req BanUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.userService.BanUser(c.Request.Context(), userID, req.Reason); err != nil {
		h.logger.Error("failed to ban user",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.String("request_id", c.GetString("request_id")),
		)
		response.InternalError(c, "Failed to ban user")
		return
	}

	response.Success(c, http.StatusOK, "User banned successfully", nil)
}

// UnbanUser unbans a user
// @Summary Unban user
// @Description Unbans a previously banned user (admin only)
// @Tags Admin
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} response.Response
// @Failure 401 {object} response.Response
// @Failure 403 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /v1/admin/users/{id}/unban [post]
func (h *UserHandler) UnbanUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		response.ValidationError(c, "Invalid user ID")
		return
	}

	if err := h.userService.UnbanUser(c.Request.Context(), userID); err != nil {
		h.logger.Error("failed to unban user",
			zap.Error(err),
			zap.Int64("user_id", userID),
			zap.String("request_id", c.GetString("request_id")),
		)
		response.InternalError(c, "Failed to unban user")
		return
	}

	response.Success(c, http.StatusOK, "User unbanned successfully", nil)
}
