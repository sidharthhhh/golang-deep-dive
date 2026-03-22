# Auth Microservice - Implementation Plan

## 🎯 Goal
Transform the auth service into a production-ready microservice with proper code structure, separation of concerns, and all critical features.

---

## 📁 New Project Structure

```
go-auth-service/
├── cmd/
│   └── server/
│       └── main.go                          # Entry point
├── internal/
│   ├── config/
│   │   ├── config.go                        # Configuration management
│   │   └── cors.go                          # NEW: CORS configuration
│   ├── handlers/
│   │   ├── auth_handler.go                  # Auth endpoints (login, register, logout)
│   │   ├── token_handler.go                 # NEW: Token validation, refresh
│   │   ├── user_handler.go                  # NEW: User management (admin)
│   │   ├── health_handler.go                # NEW: Health check endpoints
│   │   └── password_handler.go              # NEW: Password reset flow
│   ├── middleware/
│   │   ├── auth_middleware.go               # JWT validation
│   │   ├── rate_limit_middleware.go         # NEW: Rate limiting
│   │   ├── logging_middleware.go            # NEW: Request logging
│   │   ├── cors_middleware.go               # NEW: CORS handling
│   │   └── request_id_middleware.go         # NEW: Request ID tracking
│   ├── models/
│   │   ├── user.go                          # User model
│   │   ├── token.go                         # Token models
│   │   ├── audit_log.go                     # NEW: Audit log model
│   │   └── response.go                      # NEW: Standard API responses
│   ├── repository/
│   │   ├── user_repository.go               # User database operations
│   │   ├── token_repository.go              # Token blacklist operations
│   │   └── audit_repository.go              # NEW: Audit log operations
│   ├── service/
│   │   ├── auth_service.go                  # Auth business logic
│   │   ├── token_service.go                 # NEW: Token validation service
│   │   ├── user_service.go                  # NEW: User management service
│   │   └── audit_service.go                 # NEW: Audit logging service
│   ├── routes/
│   │   ├── routes.go                        # Main router setup
│   │   ├── auth_routes.go                   # NEW: Auth route group
│   │   ├── admin_routes.go                  # NEW: Admin route group
│   │   ├── health_routes.go                 # NEW: Health route group
│   │   └── v1.go                            # NEW: API v1 versioning
│   ├── utils/
│   │   ├── jwt.go                           # JWT utilities
│   │   ├── password.go                      # Password hashing
│   │   ├── token.go                         # Token generation
│   │   ├── logger.go                        # NEW: Structured logger
│   │   ├── validator.go                     # NEW: Input validation
│   │   └── response.go                      # NEW: Response helpers
│   └── pkg/
│       ├── errors/
│       │   └── errors.go                    # NEW: Custom error types
│       └── response/
│           └── response.go                  # Standard response format
├── migrations/
│   ├── 001_create_users_table.sql
│   ├── 002_add_role_column.sql
│   ├── 003_create_token_blacklist.sql
│   └── 004_create_audit_logs.sql            # NEW: Audit logs table
├── docs/
│   ├── swagger.yaml                         # NEW: OpenAPI/Swagger spec
│   └── api.md                               # NEW: API documentation
├── .env
├── .env.example                             # NEW: Example environment file
├── Dockerfile                               # NEW: Docker configuration
├── docker-compose.yml                       # NEW: Local development setup
├── Makefile                                 # NEW: Build commands
└── README.md                                # NEW: Project documentation
```

---

## 🚀 Implementation Phases

### **Phase 1: Foundation & Structure (Day 1)**

#### 1.1 Standard Response Format
**File**: `internal/pkg/response/response.go`
```go
type Response struct {
    Success   bool        `json:"success"`
    Message   string      `json:"message,omitempty"`
    Data      interface{} `json:"data,omitempty"`
    Error     *ErrorDetail `json:"error,omitempty"`
    RequestID string      `json:"request_id,omitempty"`
    Timestamp string      `json:"timestamp"`
}

type ErrorDetail struct {
    Code    string `json:"code"`
    Message string `json:"message"`
    Details string `json:"details,omitempty"`
}
```

#### 1.2 Custom Error Types
**File**: `internal/pkg/errors/errors.go`
```go
type AppError struct {
    Code       string
    Message    string
    StatusCode int
    Err        error
}

var (
    ErrUnauthorized     = &AppError{Code: "UNAUTHORIZED", StatusCode: 401}
    ErrForbidden        = &AppError{Code: "FORBIDDEN", StatusCode: 403}
    ErrNotFound         = &AppError{Code: "NOT_FOUND", StatusCode: 404}
    ErrRateLimitExceeded = &AppError{Code: "RATE_LIMIT_EXCEEDED", StatusCode: 429}
)
```

#### 1.3 Structured Logger
**File**: `internal/utils/logger.go`
```go
// Using zap for structured logging
func NewLogger() *zap.Logger
func LogRequest(logger *zap.Logger, method, path string, duration time.Duration)
func LogError(logger *zap.Logger, err error, context map[string]interface{})
```

#### 1.4 Request ID Middleware
**File**: `internal/middleware/request_id_middleware.go`
```go
func RequestIDMiddleware() gin.HandlerFunc {
    // Generate UUID for each request
    // Add to context and response header
}
```

---

### **Phase 2: Token Validation Endpoint (Day 1-2)**

#### 2.1 Token Service
**File**: `internal/service/token_service.go`
```go
type TokenService interface {
    ValidateToken(ctx context.Context, token string) (*TokenValidationResult, error)
    GetTokenInfo(ctx context.Context, token string) (*TokenInfo, error)
}

type TokenValidationResult struct {
    Valid       bool
    UserID      int
    Email       string
    Role        string
    Permissions []string
    ExpiresAt   time.Time
}
```

#### 2.2 Token Handler
**File**: `internal/handlers/token_handler.go`
```go
type TokenHandler struct {
    tokenService service.TokenService
    logger       *zap.Logger
}

// POST /auth/validate
func (h *TokenHandler) ValidateToken(c *gin.Context)

// POST /auth/refresh
func (h *TokenHandler) RefreshToken(c *gin.Context)

// GET /auth/token-info
func (h *TokenHandler) GetTokenInfo(c *gin.Context)
```

#### 2.3 Token Routes
**File**: `internal/routes/auth_routes.go`
```go
func SetupAuthRoutes(router *gin.RouterGroup, handlers *Handlers, middleware *Middleware) {
    auth := router.Group("/auth")
    {
        // Public endpoints
        auth.POST("/register", handlers.Auth.Register)
        auth.POST("/login", handlers.Auth.Login)
        auth.POST("/validate", handlers.Token.ValidateToken)
        
        // Protected endpoints
        protected := auth.Group("")
        protected.Use(middleware.Auth)
        {
            protected.POST("/refresh", handlers.Token.RefreshToken)
            protected.POST("/logout", handlers.Auth.Logout)
            protected.GET("/token-info", handlers.Token.GetTokenInfo)
        }
    }
}
```

---

### **Phase 3: CORS & Rate Limiting (Day 2)**

#### 3.1 CORS Configuration
**File**: `internal/config/cors.go`
```go
type CORSConfig struct {
    AllowedOrigins   []string
    AllowedMethods   []string
    AllowedHeaders   []string
    ExposedHeaders   []string
    AllowCredentials bool
    MaxAge           int
}

func LoadCORSConfig() *CORSConfig
```

#### 3.2 CORS Middleware
**File**: `internal/middleware/cors_middleware.go`
```go
func CORSMiddleware(config *config.CORSConfig) gin.HandlerFunc {
    return cors.New(cors.Config{
        AllowOrigins:     config.AllowedOrigins,
        AllowMethods:     config.AllowedMethods,
        AllowHeaders:     config.AllowedHeaders,
        ExposeHeaders:    config.ExposedHeaders,
        AllowCredentials: config.AllowCredentials,
        MaxAge:           time.Duration(config.MaxAge) * time.Hour,
    })
}
```

#### 3.3 Rate Limiting Middleware
**File**: `internal/middleware/rate_limit_middleware.go`
```go
import "github.com/ulule/limiter/v3"

func RateLimitMiddleware(rate limiter.Rate) gin.HandlerFunc {
    // IP-based rate limiting
    // Different limits for different endpoints
}

// Usage:
// auth.POST("/login", RateLimitMiddleware(limiter.Rate{
//     Period: 15 * time.Minute,
//     Limit:  5,
// }), handler.Login)
```

---

### **Phase 4: Logging & Monitoring (Day 2-3)**

#### 4.1 Logging Middleware
**File**: `internal/middleware/logging_middleware.go`
```go
func LoggingMiddleware(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        
        c.Next()
        
        duration := time.Since(start)
        logger.Info("request",
            zap.String("request_id", c.GetString("request_id")),
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("duration", duration),
            zap.String("ip", c.ClientIP()),
        )
    }
}
```

#### 4.2 Audit Logging
**File**: `internal/models/audit_log.go`
```go
type AuditLog struct {
    ID        int64
    UserID    int64
    Action    string  // login, logout, register, role_change, etc.
    Resource  string  // user, token, admin
    IPAddress string
    UserAgent string
    Success   bool
    Details   string  // JSON
    CreatedAt time.Time
}
```

**File**: `internal/service/audit_service.go`
```go
type AuditService interface {
    LogAction(ctx context.Context, log *models.AuditLog) error
    GetUserAuditLogs(ctx context.Context, userID int64, limit int) ([]*models.AuditLog, error)
}
```

---

### **Phase 5: Health Checks (Day 3)**

#### 5.1 Health Handler
**File**: `internal/handlers/health_handler.go`
```go
type HealthHandler struct {
    db     *sql.DB
    logger *zap.Logger
}

// GET /health
func (h *HealthHandler) Health(c *gin.Context) {
    c.JSON(200, gin.H{
        "status":    "healthy",
        "version":   "1.0.0",
        "timestamp": time.Now().Format(time.RFC3339),
    })
}

// GET /health/ready (Kubernetes readiness probe)
func (h *HealthHandler) Ready(c *gin.Context) {
    // Check database connection
    if err := h.db.Ping(); err != nil {
        c.JSON(503, gin.H{"status": "not ready", "error": "database unavailable"})
        return
    }
    c.JSON(200, gin.H{"status": "ready"})
}

// GET /health/live (Kubernetes liveness probe)
func (h *HealthHandler) Live(c *gin.Context) {
    c.JSON(200, gin.H{"status": "alive"})
}
```

#### 5.2 Health Routes
**File**: `internal/routes/health_routes.go`
```go
func SetupHealthRoutes(router *gin.Engine, handler *handlers.HealthHandler) {
    health := router.Group("/health")
    {
        health.GET("", handler.Health)
        health.GET("/ready", handler.Ready)
        health.GET("/live", handler.Live)
    }
}
```

---

### **Phase 6: API Versioning (Day 3)**

#### 6.1 Version 1 Routes
**File**: `internal/routes/v1.go`
```go
func SetupV1Routes(router *gin.Engine, handlers *Handlers, middleware *Middleware) {
    v1 := router.Group("/v1")
    v1.Use(middleware.RequestID)
    v1.Use(middleware.Logging)
    v1.Use(middleware.CORS)
    
    {
        SetupAuthRoutes(v1, handlers, middleware)
        SetupAdminRoutes(v1, handlers, middleware)
    }
}
```

---

### **Phase 7: User Management APIs (Day 4)**

#### 7.1 User Service
**File**: `internal/service/user_service.go`
```go
type UserService interface {
    GetAllUsers(ctx context.Context, page, limit int) ([]*models.User, int, error)
    GetUserByID(ctx context.Context, userID int64) (*models.User, error)
    UpdateUser(ctx context.Context, userID int64, updates *UserUpdate) error
    DeleteUser(ctx context.Context, userID int64) error
    BanUser(ctx context.Context, userID int64, reason string) error
    UnbanUser(ctx context.Context, userID int64) error
}
```

#### 7.2 User Handler
**File**: `internal/handlers/user_handler.go`
```go
type UserHandler struct {
    userService  service.UserService
    auditService service.AuditService
    logger       *zap.Logger
}

// GET /admin/users
func (h *UserHandler) ListUsers(c *gin.Context)

// GET /admin/users/:id
func (h *UserHandler) GetUser(c *gin.Context)

// PUT /admin/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context)

// DELETE /admin/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context)

// POST /admin/users/:id/ban
func (h *UserHandler) BanUser(c *gin.Context)

// POST /admin/users/:id/unban
func (h *UserHandler) UnbanUser(c *gin.Context)
```

#### 7.3 Admin Routes
**File**: `internal/routes/admin_routes.go`
```go
func SetupAdminRoutes(router *gin.RouterGroup, handlers *Handlers, middleware *Middleware) {
    admin := router.Group("/admin")
    admin.Use(middleware.Auth)
    admin.Use(middleware.RequireRole("admin", "super_admin"))
    {
        // User management
        users := admin.Group("/users")
        {
            users.GET("", handlers.User.ListUsers)
            users.GET("/:id", handlers.User.GetUser)
            users.PUT("/:id", handlers.User.UpdateUser)
            users.DELETE("/:id", handlers.User.DeleteUser)
            users.POST("/:id/ban", handlers.User.BanUser)
            users.POST("/:id/unban", handlers.User.UnbanUser)
        }
        
        // Super admin only
        superAdmin := admin.Group("/super")
        superAdmin.Use(middleware.RequireRole("super_admin"))
        {
            superAdmin.POST("/promote", handlers.Auth.PromoteToAdmin)
        }
    }
}
```

---

### **Phase 8: Password Reset Flow (Day 4-5)**

#### 8.1 Password Handler
**File**: `internal/handlers/password_handler.go`
```go
type PasswordHandler struct {
    authService service.AuthService
    logger      *zap.Logger
}

// POST /auth/forgot-password
func (h *PasswordHandler) ForgotPassword(c *gin.Context) {
    // Generate reset token
    // Send email (or return token for testing)
}

// POST /auth/reset-password
func (h *PasswordHandler) ResetPassword(c *gin.Context) {
    // Validate reset token
    // Update password
    // Invalidate all user sessions
}

// POST /auth/change-password (authenticated)
func (h *PasswordHandler) ChangePassword(c *gin.Context) {
    // Verify old password
    // Update to new password
}
```

---

### **Phase 9: Swagger Documentation (Day 5)**

#### 9.1 Swagger Setup
**File**: `docs/swagger.yaml`
```yaml
openapi: 3.0.0
info:
  title: Auth Microservice API
  version: 1.0.0
  description: Authentication and authorization microservice

paths:
  /v1/auth/login:
    post:
      summary: Login user
      tags: [Authentication]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: '#/components/schemas/LoginRequest'
      responses:
        '200':
          description: Login successful
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/LoginResponse'
```

#### 9.2 Swagger Integration
```go
import swaggerFiles "github.com/swaggo/files"
import ginSwagger "github.com/swaggo/gin-swagger"

// In routes setup
router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
```

---

### **Phase 10: Docker & Deployment (Day 5-6)**

#### 10.1 Dockerfile
**File**: `Dockerfile`
```dockerfile
# Multi-stage build
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o auth-service cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/auth-service .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8080
CMD ["./auth-service"]
```

#### 10.2 Docker Compose
**File**: `docker-compose.yml`
```yaml
version: '3.8'
services:
  auth-service:
    build: .
    ports:
      - "8080:8080"
    environment:
      - APP_ENV=development
      - MYSQL_HOST=mysql
    depends_on:
      - mysql
  
  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: password
      MYSQL_DATABASE: auth_service
    ports:
      - "3306:3306"
```

#### 10.3 Makefile
**File**: `Makefile`
```makefile
.PHONY: build run test docker-build docker-run migrate

build:
	go build -o bin/auth-service cmd/server/main.go

run:
	go run cmd/server/main.go

test:
	go test -v ./...

docker-build:
	docker build -t auth-service:latest .

docker-run:
	docker-compose up -d

migrate-up:
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/auth_service" up

migrate-down:
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/auth_service" down
```

---

## 📋 Implementation Checklist

### Day 1: Foundation
- [ ] Create standard response format
- [ ] Create custom error types
- [ ] Setup structured logger (zap)
- [ ] Implement request ID middleware
- [ ] Create token validation endpoint
- [ ] Create token service

### Day 2: Security & Middleware
- [ ] Implement CORS configuration
- [ ] Implement CORS middleware
- [ ] Implement rate limiting middleware
- [ ] Implement logging middleware
- [ ] Create audit log model
- [ ] Create audit service

### Day 3: Health & Versioning
- [ ] Create health check endpoints
- [ ] Implement readiness probe
- [ ] Implement liveness probe
- [ ] Setup API versioning (v1)
- [ ] Reorganize routes by version

### Day 4: User Management
- [ ] Create user service
- [ ] Create user handler
- [ ] Implement list users endpoint
- [ ] Implement get user endpoint
- [ ] Implement update user endpoint
- [ ] Implement delete user endpoint
- [ ] Implement ban/unban endpoints

### Day 5: Password & Documentation
- [ ] Create password handler
- [ ] Implement forgot password
- [ ] Implement reset password
- [ ] Implement change password
- [ ] Create Swagger documentation
- [ ] Setup Swagger UI

### Day 6: Deployment
- [ ] Create Dockerfile
- [ ] Create docker-compose.yml
- [ ] Create Makefile
- [ ] Write README.md
- [ ] Create .env.example
- [ ] Test full deployment

---

## 🎯 Success Criteria

### Functional Requirements
✅ Token validation endpoint works
✅ CORS allows frontend access
✅ Rate limiting prevents abuse
✅ All requests are logged
✅ Health checks respond correctly
✅ API documentation is accessible
✅ User management APIs work
✅ Password reset flow works

### Non-Functional Requirements
✅ Response time < 100ms for token validation
✅ Can handle 1000 requests/second
✅ Logs are structured (JSON)
✅ Docker image < 50MB
✅ Zero downtime deployment
✅ API documentation is up-to-date

---

## 📊 Testing Strategy

### Unit Tests
```go
// internal/service/token_service_test.go
func TestValidateToken(t *testing.T)
func TestExpiredToken(t *testing.T)
func TestBlacklistedToken(t *testing.T)
```

### Integration Tests
```go
// tests/integration/auth_test.go
func TestLoginFlow(t *testing.T)
func TestTokenValidationFlow(t *testing.T)
func TestRateLimiting(t *testing.T)
```

### Load Tests
```bash
# Using k6 or Apache Bench
k6 run load-test.js
```

---

## 🚀 Deployment Strategy

### Local Development
```bash
make run
```

### Docker
```bash
make docker-build
make docker-run
```

### Kubernetes
```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
kubectl apply -f k8s/ingress.yaml
```

---

## 📝 Next Steps After Implementation

1. **Monitoring**: Add Prometheus metrics
2. **Tracing**: Add OpenTelemetry
3. **Caching**: Add Redis for token validation
4. **Email**: Integrate email service for password reset
5. **MFA**: Add two-factor authentication
6. **OAuth**: Add social login (Google, GitHub)

---

Ready to start implementation? Let me know which phase you want to begin with!
