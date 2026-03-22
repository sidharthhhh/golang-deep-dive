package config

import (
	"os"
	"strings"
	"time"
)

// CORSConfig holds CORS configuration
type CORSConfig struct {
	AllowedOrigins   []string
	AllowedMethods   []string
	AllowedHeaders   []string
	ExposedHeaders   []string
	AllowCredentials bool
	MaxAge           time.Duration
}

// LoadCORSConfig loads CORS configuration from environment
func LoadCORSConfig() *CORSConfig {
	// Get allowed origins from environment
	originsEnv := os.Getenv("CORS_ALLOWED_ORIGINS")
	var origins []string
	if originsEnv != "" {
		origins = strings.Split(originsEnv, ",")
	} else {
		// Default origins for development
		origins = []string{
			"http://localhost:3000",
			"http://localhost:3001",
			"http://localhost:8080",
		}
	}

	return &CORSConfig{
		AllowedOrigins: origins,
		AllowedMethods: []string{
			"GET",
			"POST",
			"PUT",
			"PATCH",
			"DELETE",
			"OPTIONS",
		},
		AllowedHeaders: []string{
			"Origin",
			"Content-Type",
			"Accept",
			"Authorization",
			"X-Request-ID",
		},
		ExposedHeaders: []string{
			"X-Request-ID",
		},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}
}
