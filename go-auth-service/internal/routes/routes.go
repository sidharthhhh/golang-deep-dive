package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/handlers"
)

func SetupRouter(authHandler *handlers.AuthHandler) *gin.Engine {

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
	}

	return router
}
