package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/config"
	"github.com/sidharthhhh/go-auth-service/internal/handlers"
	"github.com/sidharthhhh/go-auth-service/internal/middleware"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
	"go.uber.org/zap"
)

// V1Handlers holds all v1 handlers
type V1Handlers struct {
	Auth   *handlers.AuthHandler
	Token  *handlers.TokenHandler
	User   *handlers.UserHandler
	Health *handlers.HealthHandler
}

// SetupV1Routes sets up all v1 API routes
func SetupV1Routes(
	router *gin.Engine,
	handlers *V1Handlers,
	jwtSecret string,
	tokenRepo repository.TokenRepository,
	corsConfig *config.CORSConfig,
	logger *zap.Logger,
) {
	// Global middleware
	router.Use(middleware.RequestIDMiddleware())
	router.Use(middleware.LoggingMiddleware(logger))
	router.Use(middleware.CORSMiddleware(corsConfig))

	// Health routes (no versioning)
	SetupHealthRoutes(router, handlers.Health)

	// API v1
	v1 := router.Group("/v1")
	{
		// Auth routes
		SetupAuthRoutes(v1, handlers.Auth, handlers.Token, jwtSecret, tokenRepo)

		// Admin routes
		SetupAdminRoutes(v1, handlers.Auth, handlers.User, jwtSecret, tokenRepo)
	}
}
