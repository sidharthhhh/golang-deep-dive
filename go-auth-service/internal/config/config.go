package config

import (
	"errors"
	"fmt"
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

func LoadConfig() (*Config, error) {

	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found")
	}

	cfg := &Config{
		AppPort:       os.Getenv("APP_PORT"),
		MysqlHost:     os.Getenv("MYSQL_HOST"),
		MysqlPort:     os.Getenv("MYSQL_PORT"),
		MysqlUser:     os.Getenv("MYSQL_USER"),
		MysqlPassword: os.Getenv("MYSQL_PASSWORD"),
		MysqlDB:       os.Getenv("MYSQL_DB"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
	}

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	var missingFields []string

	if c.AppPort == "" {
		missingFields = append(missingFields, "APP_PORT")
	}
	if c.MysqlHost == "" {
		missingFields = append(missingFields, "MYSQL_HOST")
	}
	if c.MysqlPort == "" {
		missingFields = append(missingFields, "MYSQL_PORT")
	}
	if c.MysqlUser == "" {
		missingFields = append(missingFields, "MYSQL_USER")
	}
	if c.MysqlPassword == "" {
		missingFields = append(missingFields, "MYSQL_PASSWORD")
	}
	if c.MysqlDB == "" {
		missingFields = append(missingFields, "MYSQL_DB")
	}
	if c.JWTSecret == "" {
		missingFields = append(missingFields, "JWT_SECRET")
	}

	if len(missingFields) > 0 {
		return fmt.Errorf("missing required environment variables: %v", missingFields)
	}

	if len(c.JWTSecret) < 32 {
		return errors.New("JWT_SECRET must be at least 32 characters long")
	}

	return nil
}
