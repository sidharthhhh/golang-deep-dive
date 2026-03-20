package main

import (
	"log"

	"github.com/sidharthhhh/go-auth-service/internal/config"
	"github.com/sidharthhhh/go-auth-service/internal/database"
	"github.com/sidharthhhh/go-auth-service/internal/handlers"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
	"github.com/sidharthhhh/go-auth-service/internal/routes"
	"github.com/sidharthhhh/go-auth-service/internal/service"
)

func main() {

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatal("Config error:", err)
	}

	db, err := database.NewMySQLConnection(cfg)
	if err != nil {
		log.Fatal("Database error:", err)
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := handlers.NewAuthHandler(authService)

	router := routes.SetupRouter(authHandler)

	log.Println("Server running on port", cfg.AppPort)

	router.Run(":" + cfg.AppPort)
}
