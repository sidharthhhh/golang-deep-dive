package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new structured logger
func NewLogger(env string) (*zap.Logger, error) {
	var config zap.Config

	if env == "production" {
		config = zap.NewProductionConfig()
	} else {
		config = zap.NewDevelopmentConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return nil, err
	}

	return logger, nil
}

// LogRequest logs HTTP request information
func LogRequest(logger *zap.Logger, method, path, requestID, ip string, status int, duration time.Duration) {
	logger.Info("http_request",
		zap.String("request_id", requestID),
		zap.String("method", method),
		zap.String("path", path),
		zap.Int("status", status),
		zap.Duration("duration", duration),
		zap.String("ip", ip),
	)
}

// LogError logs error information
func LogError(logger *zap.Logger, err error, context map[string]interface{}) {
	fields := make([]zap.Field, 0, len(context)+1)
	fields = append(fields, zap.Error(err))

	for key, value := range context {
		fields = append(fields, zap.Any(key, value))
	}

	logger.Error("error", fields...)
}

// LogAuth logs authentication events
func LogAuth(logger *zap.Logger, action, email, ip string, success bool) {
	logger.Info("auth_event",
		zap.String("action", action),
		zap.String("email", email),
		zap.String("ip", ip),
		zap.Bool("success", success),
	)
}

// GetLoggerFromEnv creates a logger based on APP_ENV
func GetLoggerFromEnv() (*zap.Logger, error) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	return NewLogger(env)
}


// Helper functions for zap fields
func String(key, val string) zap.Field {
	return zap.String(key, val)
}

func Error(err error) zap.Field {
	return zap.Error(err)
}

func Int(key string, val int) zap.Field {
	return zap.Int(key, val)
}

func Bool(key string, val bool) zap.Field {
	return zap.Bool(key, val)
}
