package response

import (
	"time"

	"github.com/gin-gonic/gin"
)

// Response represents the standard API response structure
type Response struct {
	Success   bool         `json:"success"`
	Message   string       `json:"message,omitempty"`
	Data      interface{}  `json:"data,omitempty"`
	Error     *ErrorDetail `json:"error,omitempty"`
	RequestID string       `json:"request_id,omitempty"`
	Timestamp string       `json:"timestamp"`
}

// ErrorDetail represents detailed error information
type ErrorDetail struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}

// Success sends a successful response
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	requestID, _ := c.Get("request_id")
	c.JSON(statusCode, Response{
		Success:   true,
		Message:   message,
		Data:      data,
		RequestID: requestID.(string),
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// Error sends an error response
func Error(c *gin.Context, statusCode int, code, message, details string) {
	requestID, _ := c.Get("request_id")
	c.JSON(statusCode, Response{
		Success: false,
		Error: &ErrorDetail{
			Code:    code,
			Message: message,
			Details: details,
		},
		RequestID: requestID.(string),
		Timestamp: time.Now().Format(time.RFC3339),
	})
}

// ValidationError sends a validation error response
func ValidationError(c *gin.Context, details string) {
	Error(c, 400, "VALIDATION_ERROR", "Invalid input", details)
}

// Unauthorized sends an unauthorized error response
func Unauthorized(c *gin.Context, message string) {
	Error(c, 401, "UNAUTHORIZED", message, "")
}

// Forbidden sends a forbidden error response
func Forbidden(c *gin.Context, message string) {
	Error(c, 403, "FORBIDDEN", message, "")
}

// NotFound sends a not found error response
func NotFound(c *gin.Context, message string) {
	Error(c, 404, "NOT_FOUND", message, "")
}

// InternalError sends an internal server error response
func InternalError(c *gin.Context, message string) {
	Error(c, 500, "INTERNAL_ERROR", message, "")
}

// RateLimitExceeded sends a rate limit exceeded error response
func RateLimitExceeded(c *gin.Context) {
	Error(c, 429, "RATE_LIMIT_EXCEEDED", "Too many requests", "Please try again later")
}
