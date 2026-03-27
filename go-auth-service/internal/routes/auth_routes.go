package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/handlers"
	"github.com/sidharthhhh/go-auth-service/internal/middleware"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
)

// SetupAuthRoutes sets up authentication routes
func SetupAuthRoutes(
	router *gin.RouterGroup,
	authHandler *handlers.AuthHandler,
	tokenHandler *handlers.TokenHandler,
	passwordHandler *handlers.PasswordHandler,
	jwtSecret string,
	tokenRepo repository.TokenRepository,
	userRepo repository.UserRepository,
) {
	auth := router.Group("/auth")
	{
		// Public endpoints
		auth.POST("/register",
			middleware.RegisterRateLimit(),
			authHandler.Register,
		)

		auth.POST("/login",
			middleware.LoginRateLimit(),
			authHandler.Login,
		)

		auth.POST("/validate",
			middleware.TokenValidationRateLimit(),
			tokenHandler.ValidateToken,
		)

		// Password reset endpoints (public)
		auth.POST("/forgot-password", passwordHandler.ForgotPassword)
		auth.POST("/reset-password", passwordHandler.ResetPassword)

		// Protected endpoints (require authentication)
		protected := auth.Group("")
		protected.Use(middleware.AuthMiddleware(jwtSecret, tokenRepo, userRepo))
		{
			protected.POST("/refresh", func(c *gin.Context) {
				c.Set("jwt_secret", jwtSecret)
				tokenHandler.RefreshToken(c)
			})

			protected.POST("/logout", func(c *gin.Context) {
				c.Set("jwt_secret", jwtSecret)
				authHandler.Logout(c)
			})

			protected.GET("/token-info", tokenHandler.GetTokenInfo)
			
			// Change password (requires authentication)
			protected.POST("/change-password", passwordHandler.ChangePassword)
		}
	}

	// User profile endpoint
	api := router.Group("/api")
	api.Use(middleware.AuthMiddleware(jwtSecret, tokenRepo, userRepo))
	{
		api.GET("/profile", func(c *gin.Context) {
			userID := c.GetInt("user_id")
			email := c.GetString("email")
			role := c.GetString("role")

			c.JSON(200, gin.H{
				"success": true,
				"data": gin.H{
					"user_id": userID,
					"email":   email,
					"role":    role,
				},
				"timestamp": c.GetString("timestamp"),
			})
		})
	}
}
