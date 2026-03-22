package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/pkg/response"
	"github.com/ulule/limiter/v3"
	"github.com/ulule/limiter/v3/drivers/store/memory"
)

// RateLimitMiddleware creates a rate limiting middleware
func RateLimitMiddleware(rate limiter.Rate) gin.HandlerFunc {
	store := memory.NewStore()
	instance := limiter.New(store, rate)

	return func(c *gin.Context) {
		// Get the limiter context
		limiterCtx, err := instance.Get(c, c.ClientIP())
		if err != nil {
			c.Next()
			return
		}

		// Check if rate limit is exceeded
		if limiterCtx.Reached {
			response.RateLimitExceeded(c)
			c.Abort()
			return
		}

		// Add rate limit headers
		c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limiterCtx.Limit))
		c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", limiterCtx.Remaining))
		c.Header("X-RateLimit-Reset", fmt.Sprintf("%d", limiterCtx.Reset))

		c.Next()
	}
}

// LoginRateLimit returns a rate limiter for login endpoint (5 requests per 15 minutes)
func LoginRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(limiter.Rate{
		Period: 15 * time.Minute,
		Limit:  5,
	})
}

// RegisterRateLimit returns a rate limiter for register endpoint (3 requests per hour)
func RegisterRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(limiter.Rate{
		Period: 1 * time.Hour,
		Limit:  3,
	})
}

// TokenValidationRateLimit returns a rate limiter for token validation (1000 requests per minute)
func TokenValidationRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  1000,
	})
}

// DefaultRateLimit returns a default rate limiter (100 requests per minute)
func DefaultRateLimit() gin.HandlerFunc {
	return RateLimitMiddleware(limiter.Rate{
		Period: 1 * time.Minute,
		Limit:  100,
	})
}
