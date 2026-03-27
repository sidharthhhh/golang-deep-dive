	package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/handlers"
	"github.com/sidharthhhh/go-auth-service/internal/middleware"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
)

func SetupRouter(authHandler *handlers.AuthHandler, jwtSecret string, tokenRepo repository.TokenRepository, userRepo repository.UserRepository) *gin.Engine {

	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"service": "go-auth-service",
		})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/register", authHandler.Register)
		auth.POST("/login", authHandler.Login)
		
		// Refresh token endpoint (requires authentication)
		auth.POST("/refresh", middleware.AuthMiddleware(jwtSecret, tokenRepo, userRepo), func(c *gin.Context) {
			c.Set("jwt_secret", jwtSecret)
			authHandler.RefreshToken(c)
		})

		// Logout endpoint (requires authentication)
		auth.POST("/logout", middleware.AuthMiddleware(jwtSecret, tokenRepo, userRepo), func(c *gin.Context) {
			c.Set("jwt_secret", jwtSecret)
			authHandler.Logout(c)
		})
	}

	// Protected routes - requires authentication
	protected := router.Group("/api")
	protected.Use(middleware.AuthMiddleware(jwtSecret, tokenRepo, userRepo))
	{
		protected.GET("/profile", func(c *gin.Context) {
			userID := c.GetInt("user_id")
			email := c.GetString("email")
			role := c.GetString("role")
			c.JSON(200, gin.H{
				"user_id": userID,
				"email":   email,
				"role":    role,
				"message": "This is a protected route accessible by all authenticated users",
			})
		})
	}

	// Admin routes - requires admin or super_admin role
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtSecret, tokenRepo, userRepo))
	admin.Use(middleware.RequireRole("admin", "super_admin"))
	{
		admin.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "Welcome to admin dashboard",
				"role":    c.GetString("role"),
			})
		})

		admin.GET("/users", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "List of all users (admin only)",
			})
		})
	}

	// Super admin routes - requires super_admin role only
	superAdmin := router.Group("/super-admin")
	superAdmin.Use(middleware.AuthMiddleware(jwtSecret, tokenRepo, userRepo))
	superAdmin.Use(middleware.RequireRole("super_admin"))
	{
		superAdmin.POST("/promote", authHandler.PromoteToAdmin)

		superAdmin.GET("/system", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "System settings (super admin only)",
			})
		})
	}

	return router
}
