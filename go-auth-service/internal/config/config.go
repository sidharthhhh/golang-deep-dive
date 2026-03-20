package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort       string
	MysqlHost     string
	MysqlPort     string
	MysqlUser     string
	MysqlPassword string
	MysqlDB       string
	JWTSecret     string
}

func LoadConfig() *Config {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	return &Config{
		AppPort:       os.Getenv("APP_PORT"),
		MysqlHost:     os.Getenv("MYSQL_HOST"),
		MysqlPort:     os.Getenv("MYSQL_PORT"),
		MysqlUser:     os.Getenv("MYSQL_USER"),
		MysqlPassword: os.Getenv("MYSQL_PASSWORD"),
		MysqlDB:       os.Getenv("MYSQL_DB"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
	}
}
