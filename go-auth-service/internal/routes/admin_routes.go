package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/handlers"
	"github.com/sidharthhhh/go-auth-service/internal/middleware"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
)

// SetupAdminRoutes sets up admin routes
func SetupAdminRoutes(
	router *gin.RouterGroup,
	authHandler *handlers.AuthHandler,
	userHandler *handlers.UserHandler,
	jwtSecret string,
	tokenRepo repository.TokenRepository,
	userRepo repository.UserRepository,
) {
	admin := router.Group("/admin")
	admin.Use(middleware.AuthMiddleware(jwtSecret, tokenRepo, userRepo))
	admin.Use(middleware.RequireRole("admin", "super_admin"))
	{
		// Dashboard
		admin.GET("/dashboard", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"message": "Welcome to admin dashboard",
				"data": gin.H{
					"role": c.GetString("role"),
				},
			})
		})

		// User management
		users := admin.Group("/users")
		{
			users.GET("", userHandler.ListUsers)
			users.GET("/:id", userHandler.GetUser)
			users.PUT("/:id", userHandler.UpdateUser)
			users.DELETE("/:id", userHandler.DeleteUser)
			users.POST("/:id/ban", userHandler.BanUser)
			users.POST("/:id/unban", userHandler.UnbanUser)
		}
	}

	// Super admin only routes
	superAdmin := router.Group("/super-admin")
	superAdmin.Use(middleware.AuthMiddleware(jwtSecret, tokenRepo, userRepo))
	superAdmin.Use(middleware.RequireRole("super_admin"))
	{
		superAdmin.POST("/promote", authHandler.PromoteToAdmin)

		superAdmin.GET("/system", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"success": true,
				"message": "System settings",
				"data": gin.H{
					"role": "super_admin",
				},
			})
		})
	}
}
