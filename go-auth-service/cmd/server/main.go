package main

import (
	"log"

	"github.com/sidharthhhh/go-auth-service/internal/config"
	"github.com/sidharthhhh/go-auth-service/internal/database"
	"github.com/sidharthhhh/go-auth-service/internal/routes"
)

func main() {

	cfg := config.LoadConfig()

	db := database.NewMySQLConnection(cfg)
	defer db.Close()

	router := routes.SetupRouter()

	log.Println("Server running on port", cfg.AppPort)

	router.Run(":" + cfg.AppPort)
}
