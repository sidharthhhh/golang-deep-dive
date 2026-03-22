package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/utils"
	"go.uber.org/zap"
)

// LoggingMiddleware logs all HTTP requests
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method
		ip := c.ClientIP()

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)
		status := c.Writer.Status()
		requestID := c.GetString("request_id")

		// Log request
		utils.LogRequest(logger, method, path, requestID, ip, status, duration)

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, err := range c.Errors {
				utils.LogError(logger, err.Err, map[string]interface{}{
					"request_id": requestID,
					"path":       path,
					"method":     method,
				})
			}
		}
	}
}
