# Integration Guide - Go Auth Service

Complete guide for integrating the Go Auth Service with your microservices.

## Table of Contents

1. [Overview](#overview)
2. [Architecture](#architecture)
3. [Integration Methods](#integration-methods)
4. [Code Examples](#code-examples)
5. [Best Practices](#best-practices)
6. [Troubleshooting](#troubleshooting)

---

## Overview

The Go Auth Service provides centralized authentication and authorization for your microservices architecture.

### Key Features

- ✅ JWT-based authentication
- ✅ Role-based access control (RBAC)
- ✅ Token validation endpoint
- ✅ Token refresh mechanism
- ✅ User management
- ✅ Rate limiting
- ✅ Request tracking

---

## Architecture

```
┌─────────────┐         ┌──────────────────┐         ┌─────────────┐
│   Client    │────────▶│   API Gateway    │────────▶│   Service   │
│  (Browser)  │         │  (Load Balancer) │         │     A/B/C   │
└─────────────┘         └──────────────────┘         └─────────────┘
                                 │                            │
                                 │                            │
                                 ▼                            ▼
                        ┌──────────────────┐         ┌─────────────┐
                        │   Auth Service   │◀────────│  Validate   │
                        │  (Port 8080)     │         │   Token     │
                        └──────────────────┘         └─────────────┘
                                 │
                                 ▼
                        ┌──────────────────┐
                        │   MySQL DB       │
                        └──────────────────┘
```

### Flow

1. Client authenticates with Auth Service
2. Auth Service returns JWT token
3. Client includes token in requests to other services
4. Services validate token with Auth Service
5. Services process request if token is valid

---

## Integration Methods

### Method 1: Token Validation (Recommended)

Validate tokens by calling the Auth Service validation endpoint.

**Pros:**
- Always up-to-date validation
- Centralized token management
- Supports token revocation

**Cons:**
- Network call overhead
- Dependency on Auth Service availability

---

### Method 2: JWT Verification (Advanced)

Verify JWT tokens locally using the shared secret.

**Pros:**
- No network call
- Faster validation
- Works offline

**Cons:**
- Cannot check token revocation
- Requires shared secret management
- Clock synchronization needed

---

## Code Examples

### Go Service Integration

#### Method 1: Token Validation API

```go
package auth

import (
    "bytes"
    "encoding/json"
    "fmt"
    "net/http"
    "time"
)

type AuthClient struct {
    BaseURL    string
    HTTPClient *http.Client
}

type TokenValidationRequest struct {
    Token string `json:"token"`
}

type TokenValidationResponse struct {
    Success bool `json:"success"`
    Data    struct {
        Valid       bool     `json:"valid"`
        UserID      int      `json:"user_id"`
        Email       string   `json:"email"`
        Role        string   `json:"role"`
        Permissions []string `json:"permissions"`
        ExpiresAt   string   `json:"expires_at"`
    } `json:"data"`
}

func NewAuthClient(baseURL string) *AuthClient {
    return &AuthClient{
        BaseURL: baseURL,
        HTTPClient: &http.Client{
            Timeout: 5 * time.Second,
        },
    }
}

func (c *AuthClient) ValidateToken(token string) (*TokenValidationResponse, error) {
    reqBody := TokenValidationRequest{Token: token}
    jsonData, err := json.Marshal(reqBody)
    if err != nil {
        return nil, err
    }

    req, err := http.NewRequest("POST", c.BaseURL+"/v1/auth/validate", bytes.NewBuffer(jsonData))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/json")

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("token validation failed: %d", resp.StatusCode)
    }

    var result TokenValidationResponse
    if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
        return nil, err
    }

    return &result, nil
}

// Middleware for Gin
func (c *AuthClient) AuthMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        authHeader := ctx.GetHeader("Authorization")
        if authHeader == "" {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Missing authorization header"})
            ctx.Abort()
            return
        }

        // Extract token
        token := strings.TrimPrefix(authHeader, "Bearer ")
        if token == authHeader {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
            ctx.Abort()
            return
        }

        // Validate token
        validation, err := c.ValidateToken(token)
        if err != nil || !validation.Success || !validation.Data.Valid {
            ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            ctx.Abort()
            return
        }

        // Set user info in context
        ctx.Set("user_id", validation.Data.UserID)
        ctx.Set("user_email", validation.Data.Email)
        ctx.Set("user_role", validation.Data.Role)
        ctx.Set("permissions", validation.Data.Permissions)

        ctx.Next()
    }
}

// Role-based middleware
func (c *AuthClient) RequireRole(roles ...string) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        userRole, exists := ctx.Get("user_role")
        if !exists {
            ctx.JSON(http.StatusForbidden, gin.H{"error": "No role found"})
            ctx.Abort()
            return
        }

        roleStr := userRole.(string)
        for _, role := range roles {
            if roleStr == role {
                ctx.Next()
                return
            }
        }

        ctx.JSON(http.StatusForbidden, gin.H{"error": "Insufficient permissions"})
        ctx.Abort()
    }
}
```

**Usage:**

```go
package main

import (
    "github.com/gin-gonic/gin"
    "your-service/auth"
)

func main() {
    router := gin.Default()
    
    // Initialize auth client
    authClient := auth.NewAuthClient("http://auth-service:8080")
    
    // Public routes
    router.GET("/health", healthHandler)
    
    // Protected routes
    protected := router.Group("/api")
    protected.Use(authClient.AuthMiddleware())
    {
        protected.GET("/profile", profileHandler)
        protected.GET("/data", dataHandler)
    }
    
    // Admin routes
    admin := router.Group("/admin")
    admin.Use(authClient.AuthMiddleware())
    admin.Use(authClient.RequireRole("admin", "super_admin"))
    {
        admin.GET("/users", listUsersHandler)
        admin.POST("/settings", updateSettingsHandler)
    }
    
    router.Run(":8081")
}

func profileHandler(c *gin.Context) {
    userID := c.GetInt("user_id")
    email := c.GetString("user_email")
    
    c.JSON(200, gin.H{
        "user_id": userID,
        "email":   email,
    })
}
```

---

#### Method 2: Local JWT Verification

```go
package auth

import (
    "errors"
    "github.com/golang-jwt/jwt/v5"
    "time"
)

type Claims struct {
    UserID int    `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func VerifyToken(tokenString string, jwtSecret string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
        return []byte(jwtSecret), nil
    })

    if err != nil {
        return nil, err
    }

    if claims, ok := token.Claims.(*Claims); ok && token.Valid {
        // Check expiration
        if claims.ExpiresAt.Time.Before(time.Now()) {
            return nil, errors.New("token expired")
        }
        return claims, nil
    }

    return nil, errors.New("invalid token")
}

// Middleware
func JWTMiddleware(jwtSecret string) gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        if authHeader == "" {
            c.JSON(401, gin.H{"error": "Missing authorization header"})
            c.Abort()
            return
        }

        token := strings.TrimPrefix(authHeader, "Bearer ")
        claims, err := VerifyToken(token, jwtSecret)
        if err != nil {
            c.JSON(401, gin.H{"error": "Invalid token"})
            c.Abort()
            return
        }

        c.Set("user_id", claims.UserID)
        c.Set("user_email", claims.Email)
        c.Set("user_role", claims.Role)
        c.Next()
    }
}
```

---

### Node.js/Express Integration

```javascript
const axios = require('axios');

class AuthClient {
    constructor(baseURL) {
        this.baseURL = baseURL;
        this.client = axios.create({
            baseURL: baseURL,
            timeout: 5000,
        });
    }

    async validateToken(token) {
        try {
            const response = await this.client.post('/v1/auth/validate', {
                token: token
            });
            return response.data;
        } catch (error) {
            throw new Error('Token validation failed');
        }
    }

    // Middleware
    authMiddleware() {
        return async (req, res, next) => {
            const authHeader = req.headers.authorization;
            
            if (!authHeader) {
                return res.status(401).json({ error: 'Missing authorization header' });
            }

            const token = authHeader.replace('Bearer ', '');
            
            try {
                const validation = await this.validateToken(token);
                
                if (!validation.success || !validation.data.valid) {
                    return res.status(401).json({ error: 'Invalid token' });
                }

                // Attach user info to request
                req.user = {
                    id: validation.data.user_id,
                    email: validation.data.email,
                    role: validation.data.role,
                    permissions: validation.data.permissions
                };

                next();
            } catch (error) {
                return res.status(401).json({ error: 'Token validation failed' });
            }
        };
    }

    // Role-based middleware
    requireRole(...roles) {
        return (req, res, next) => {
            if (!req.user) {
                return res.status(403).json({ error: 'No user found' });
            }

            if (roles.includes(req.user.role)) {
                next();
            } else {
                return res.status(403).json({ error: 'Insufficient permissions' });
            }
        };
    }
}

// Usage
const express = require('express');
const app = express();

const authClient = new AuthClient('http://auth-service:8080');

// Public routes
app.get('/health', (req, res) => {
    res.json({ status: 'healthy' });
});

// Protected routes
app.get('/api/profile', 
    authClient.authMiddleware(), 
    (req, res) => {
        res.json({
            user_id: req.user.id,
            email: req.user.email
        });
    }
);

// Admin routes
app.get('/admin/users', 
    authClient.authMiddleware(),
    authClient.requireRole('admin', 'super_admin'),
    (req, res) => {
        res.json({ users: [] });
    }
);

app.listen(8081);
```

---

### Python/Flask Integration

```python
import requests
from functools import wraps
from flask import request, jsonify

class AuthClient:
    def __init__(self, base_url):
        self.base_url = base_url
        self.session = requests.Session()
        self.session.headers.update({'Content-Type': 'application/json'})
    
    def validate_token(self, token):
        try:
            response = self.session.post(
                f'{self.base_url}/v1/auth/validate',
                json={'token': token},
                timeout=5
            )
            response.raise_for_status()
            return response.json()
        except requests.exceptions.RequestException:
            return None
    
    def auth_required(self, f):
        @wraps(f)
        def decorated_function(*args, **kwargs):
            auth_header = request.headers.get('Authorization')
            
            if not auth_header:
                return jsonify({'error': 'Missing authorization header'}), 401
            
            token = auth_header.replace('Bearer ', '')
            validation = self.validate_token(token)
            
            if not validation or not validation.get('success') or not validation.get('data', {}).get('valid'):
                return jsonify({'error': 'Invalid token'}), 401
            
            # Attach user info to request
            request.user = {
                'id': validation['data']['user_id'],
                'email': validation['data']['email'],
                'role': validation['data']['role'],
                'permissions': validation['data']['permissions']
            }
            
            return f(*args, **kwargs)
        return decorated_function
    
    def role_required(self, *roles):
        def decorator(f):
            @wraps(f)
            def decorated_function(*args, **kwargs):
                if not hasattr(request, 'user'):
                    return jsonify({'error': 'No user found'}), 403
                
                if request.user['role'] in roles:
                    return f(*args, **kwargs)
                else:
                    return jsonify({'error': 'Insufficient permissions'}), 403
            return decorated_function
        return decorator

# Usage
from flask import Flask
app = Flask(__name__)

auth_client = AuthClient('http://auth-service:8080')

@app.route('/health')
def health():
    return jsonify({'status': 'healthy'})

@app.route('/api/profile')
@auth_client.auth_required
def profile():
    return jsonify({
        'user_id': request.user['id'],
        'email': request.user['email']
    })

@app.route('/admin/users')
@auth_client.auth_required
@auth_client.role_required('admin', 'super_admin')
def list_users():
    return jsonify({'users': []})

if __name__ == '__main__':
    app.run(port=8081)
```

---

## Best Practices

### 1. Token Caching

Cache validation results to reduce Auth Service load:

```go
type TokenCache struct {
    cache map[string]*CachedToken
    mu    sync.RWMutex
    ttl   time.Duration
}

type CachedToken struct {
    Validation *TokenValidationResponse
    ExpiresAt  time.Time
}

func (c *TokenCache) Get(token string) (*TokenValidationResponse, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    
    cached, exists := c.cache[token]
    if !exists || time.Now().After(cached.ExpiresAt) {
        return nil, false
    }
    
    return cached.Validation, true
}

func (c *TokenCache) Set(token string, validation *TokenValidationResponse) {
    c.mu.Lock()
    defer c.mu.Unlock()
    
    c.cache[token] = &CachedToken{
        Validation: validation,
        ExpiresAt:  time.Now().Add(c.ttl),
    }
}
```

### 2. Circuit Breaker

Implement circuit breaker for Auth Service calls:

```go
import "github.com/sony/gobreaker"

func NewAuthClientWithCircuitBreaker(baseURL string) *AuthClient {
    cb := gobreaker.NewCircuitBreaker(gobreaker.Settings{
        Name:        "AuthService",
        MaxRequests: 3,
        Interval:    time.Minute,
        Timeout:     10 * time.Second,
    })
    
    return &AuthClient{
        BaseURL:        baseURL,
        HTTPClient:     &http.Client{Timeout: 5 * time.Second},
        CircuitBreaker: cb,
    }
}
```

### 3. Retry Logic

Add retry logic for transient failures:

```go
func (c *AuthClient) ValidateTokenWithRetry(token string, maxRetries int) (*TokenValidationResponse, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        result, err := c.ValidateToken(token)
        if err == nil {
            return result, nil
        }
        
        lastErr = err
        time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)
    }
    
    return nil, lastErr
}
```

### 4. Request Timeout

Always set timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

req = req.WithContext(ctx)
```

### 5. Error Handling

Handle different error scenarios:

```go
func handleAuthError(err error, statusCode int) error {
    switch statusCode {
    case 401:
        return errors.New("invalid or expired token")
    case 429:
        return errors.New("rate limit exceeded")
    case 503:
        return errors.New("auth service unavailable")
    default:
        return fmt.Errorf("auth error: %v", err)
    }
}
```

---

## Configuration

### Environment Variables

```bash
# Auth Service Configuration
AUTH_SERVICE_URL=http://auth-service:8080
AUTH_SERVICE_TIMEOUT=5s
AUTH_CACHE_TTL=5m
AUTH_RETRY_COUNT=3

# JWT Configuration (for local verification)
JWT_SECRET=your-jwt-secret-key
```

### Docker Compose

```yaml
version: '3.8'

services:
  auth-service:
    image: auth-service:latest
    ports:
      - "8080:8080"
    environment:
      - MYSQL_HOST=mysql
      - JWT_SECRET=${JWT_SECRET}
    depends_on:
      - mysql

  your-service:
    image: your-service:latest
    ports:
      - "8081:8081"
    environment:
      - AUTH_SERVICE_URL=http://auth-service:8080
    depends_on:
      - auth-service
```

### Kubernetes

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: your-service-config
data:
  AUTH_SERVICE_URL: "http://auth-service:8080"
  AUTH_CACHE_TTL: "5m"
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: your-service
spec:
  replicas: 3
  template:
    spec:
      containers:
      - name: your-service
        image: your-service:latest
        envFrom:
        - configMapRef:
            name: your-service-config
```

---

## Troubleshooting

### Issue: Token Validation Fails

**Symptoms:**
- 401 Unauthorized responses
- "Invalid token" errors

**Solutions:**
1. Check token format: `Bearer <token>`
2. Verify token hasn't expired
3. Ensure Auth Service is reachable
4. Check JWT secret matches (if using local verification)

### Issue: High Latency

**Symptoms:**
- Slow API responses
- Timeout errors

**Solutions:**
1. Implement token caching
2. Use connection pooling
3. Add circuit breaker
4. Consider local JWT verification

### Issue: Rate Limiting

**Symptoms:**
- 429 Too Many Requests
- X-RateLimit-Remaining: 0

**Solutions:**
1. Implement exponential backoff
2. Cache validation results
3. Use local JWT verification
4. Request rate limit increase

---

## Monitoring

### Metrics to Track

```go
// Prometheus metrics
var (
    authValidationDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name: "auth_validation_duration_seconds",
            Help: "Duration of auth validation requests",
        },
        []string{"status"},
    )
    
    authValidationTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "auth_validation_total",
            Help: "Total number of auth validation requests",
        },
        []string{"status"},
    )
)
```

### Health Check

```go
func (c *AuthClient) HealthCheck() error {
    resp, err := c.HTTPClient.Get(c.BaseURL + "/health")
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("auth service unhealthy: %d", resp.StatusCode)
    }
    
    return nil
}
```

---

## Support

For integration issues:
- Check Auth Service logs
- Verify network connectivity
- Review token format
- Check rate limits
- Monitor Auth Service health
