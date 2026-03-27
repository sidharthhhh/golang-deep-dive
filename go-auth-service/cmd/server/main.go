package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/config"
	"github.com/sidharthhhh/go-auth-service/internal/database"
	"github.com/sidharthhhh/go-auth-service/internal/handlers"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
	"github.com/sidharthhhh/go-auth-service/internal/routes"
	"github.com/sidharthhhh/go-auth-service/internal/service"
	"github.com/sidharthhhh/go-auth-service/internal/utils"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	// Setup logger
	logger, err := utils.GetLoggerFromEnv()
	if err != nil {
		log.Fatal("Logger error:", err)
	}
	defer logger.Sync()

	logger.Info("Starting auth service",
		utils.String("version", getVersion()),
		utils.String("port", cfg.AppPort),
	)

	// Connect to database
	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		logger.Fatal("Database error", utils.Error(err))
	}
	defer db.Close()

	logger.Info("Database connected successfully")

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	tokenRepo := repository.NewTokenRepository(db)
	passwordResetRepo := repository.NewPasswordResetRepository(db)

	// Initialize services
	authService := service.NewAuthService(userRepo, tokenRepo, cfg.JWTSecret, cfg.SuperAdminCode)
	tokenService := service.NewTokenService(tokenRepo, cfg.JWTSecret)
	userService := service.NewUserService(userRepo)
	passwordService := service.NewPasswordService(userRepo, passwordResetRepo, tokenRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(authService)
	tokenHandler := handlers.NewTokenHandler(tokenService, logger)
	userHandler := handlers.NewUserHandler(userService, logger)
	passwordHandler := handlers.NewPasswordHandler(passwordService, logger)
	healthHandler := handlers.NewHealthHandler(db, logger, getVersion())

	// Setup router
	router := gin.New()

	// Create handlers struct
	v1Handlers := &routes.V1Handlers{
		Auth:     authHandler,
		Token:    tokenHandler,
		User:     userHandler,
		Password: passwordHandler,
		Health:   healthHandler,
	}

	// Load CORS config
	corsConfig := config.LoadCORSConfig()

	// Setup all routes
	routes.SetupV1Routes(router, v1Handlers, cfg.JWTSecret, tokenRepo, userRepo, corsConfig, logger)

	// Set APP_ENV globally for handlers
	router.Use(func(c *gin.Context) {
		c.Set("APP_ENV", cfg.AppEnv)
		c.Next()
	})

	// Start server
	logger.Info("Server starting", utils.String("port", cfg.AppPort))
	if err := router.Run(":" + cfg.AppPort); err != nil {
		logger.Fatal("Server error", utils.Error(err))
	}
}

func getVersion() string {
	version := os.Getenv("APP_VERSION")
	if version == "" {
		return "1.0.0"
	}
	return version
}
