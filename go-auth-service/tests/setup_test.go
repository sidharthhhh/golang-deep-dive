package tests

import (
	"database/sql"
	"log"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/config"
	"github.com/sidharthhhh/go-auth-service/internal/database"
	"github.com/sidharthhhh/go-auth-service/internal/handlers"
	"github.com/sidharthhhh/go-auth-service/internal/repository"
	"github.com/sidharthhhh/go-auth-service/internal/routes"
	"github.com/sidharthhhh/go-auth-service/internal/service"
	"github.com/sidharthhhh/go-auth-service/internal/utils"
	"go.uber.org/zap"
)

var (
	TestRouter    *gin.Engine
	TestDB        *sql.DB
	TestLogger    *zap.Logger
	TestConfig    *config.Config
	UserRepo      repository.UserRepository
	TokenRepo     repository.TokenRepository
	PasswordRepo  repository.PasswordResetRepository
	AuthService   service.AuthService
	TokenService  service.TokenService
	UserService   service.UserService
	PasswordService service.PasswordService
)

// SetupTestEnvironment initializes the test environment
func SetupTestEnvironment() {
	// Set test environment
	os.Setenv("APP_ENV", "test")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("MYSQL_HOST", "localhost")
	os.Setenv("MYSQL_PORT", "3306")
	os.Setenv("MYSQL_USER", "root")
	os.Setenv("MYSQL_PASSWORD", "hjkl")
	os.Setenv("MYSQL_DB", "auth_service_test")
	os.Setenv("JWT_SECRET", "test-secret-key-for-testing-purposes-only")
	os.Setenv("SUPER_ADMIN_CODE", "TEST_SUPER_ADMIN_CODE")
	os.Setenv("CORS_ALLOWED_ORIGINS", "http://localhost:3000")

	// Load config
	var err error
	TestConfig, err = config.LoadConfig()
	if err != nil {
		log.Fatal("Failed to load test config:", err)
	}

	// Setup logger
	TestLogger, err = utils.GetLoggerFromEnv()
	if err != nil {
		log.Fatal("Failed to setup test logger:", err)
	}

	// Connect to database
	TestDB, err = database.NewMySQLConnection(TestConfig)
	if err != nil {
		log.Fatal("Failed to connect to test database:", err)
	}

	// Initialize repositories
	UserRepo = repository.NewUserRepository(TestDB)
	TokenRepo = repository.NewTokenRepository(TestDB)
	PasswordRepo = repository.NewPasswordResetRepository(TestDB)

	// Initialize services
	AuthService = service.NewAuthService(UserRepo, TokenRepo, TestConfig.JWTSecret, TestConfig.SuperAdminCode)
	TokenService = service.NewTokenService(TokenRepo, TestConfig.JWTSecret)
	UserService = service.NewUserService(UserRepo)
	PasswordService = service.NewPasswordService(UserRepo, PasswordRepo, TokenRepo)

	// Initialize handlers
	authHandler := handlers.NewAuthHandler(AuthService)
	tokenHandler := handlers.NewTokenHandler(TokenService, TestLogger)
	userHandler := handlers.NewUserHandler(UserService, TestLogger)
	healthHandler := handlers.NewHealthHandler(TestDB, TestLogger, "test-1.0.0")

	// Setup router
	gin.SetMode(gin.TestMode)
	TestRouter = gin.New()

	v1Handlers := &routes.V1Handlers{
		Auth:   authHandler,
		Token:  tokenHandler,
		User:   userHandler,
		Health: healthHandler,
	}

	corsConfig := config.LoadCORSConfig()
	routes.SetupV1Routes(TestRouter, v1Handlers, TestConfig.JWTSecret, TokenRepo, corsConfig, TestLogger)
}

// CleanupTestEnvironment cleans up after tests
func CleanupTestEnvironment() {
	if TestDB != nil {
		// Clean up test data
		TestDB.Exec("DELETE FROM password_reset_tokens")
		TestDB.Exec("DELETE FROM token_blacklist")
		TestDB.Exec("DELETE FROM users")
		TestDB.Close()
	}
	if TestLogger != nil {
		TestLogger.Sync()
	}
}

// TestMain runs before all tests
func TestMain(m *testing.M) {
	SetupTestEnvironment()
	code := m.Run()
	CleanupTestEnvironment()
	os.Exit(code)
}
