package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/handlers"
)

// SetupHealthRoutes sets up health check routes
func SetupHealthRoutes(router *gin.Engine, healthHandler *handlers.HealthHandler) {
	health := router.Group("/health")
	{
		health.GET("", healthHandler.Health)
		health.GET("/ready", healthHandler.Ready)
		health.GET("/live", healthHandler.Live)
	}
}
