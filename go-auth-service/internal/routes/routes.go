package routes

import "github.com/gin-gonic/gin"

func SetupRouter() *gin.Engine {

	router := gin.Default()

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"service": "go-auth-service",
		})
	})

	auth := router.Group("/auth")
	{
		auth.POST("/register", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "register endpoint"})
		})

		auth.POST("/login", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "login endpoint"})
		})
	}

	return router
}
