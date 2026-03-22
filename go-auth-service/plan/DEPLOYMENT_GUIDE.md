# Auth Service - Deployment Guide

## 🚀 Complete Implementation Summary

### ✅ All Features Implemented

1. **Foundation & Structure**
   - Standard response format
   - Custom error types
   - Structured logging (Zap)
   - Request ID tracking

2. **Token Management**
   - Token validation endpoint (for microservices)
   - Token refresh
   - Token blacklisting (logout)
   - Role-based expiry (7/30 days)

3. **Security**
   - CORS configuration
   - Rate limiting (login, register, token validation)
   - Password hashing (bcrypt)
   - JWT with permissions

4. **Health Checks**
   - Basic health (`/health`)
   - Readiness probe (`/health/ready`)
   - Liveness probe (`/health/live`)

5. **User Management**
   - List users (paginated)
   - Get user by ID
   - Update user
   - Delete user
   - Ban/unban user

6. **Password Management**
   - Forgot password
   - Reset password
   - Change password

7. **API Versioning**
   - v1 API structure
   - Organized routes

---

## 📦 Installation & Setup

### Prerequisites
```bash
- Go 1.21+
- MySQL 8.0+
- Docker & Docker Compose (optional)
- Kubernetes cluster (optional)
```

### 1. Install Dependencies

```bash
cd go-auth-service

# Install Go dependencies
make deps

# Or manually
go mod download
go mod tidy
```

### 2. Configure Environment

```bash
# Copy example env file
cp .env.example .env

# Edit .env with your values
nano .env
```

Required environment variables:
```env
APP_PORT=8080
APP_ENV=development
APP_VERSION=1.0.0
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your_password
MYSQL_DB=auth_service
JWT_SECRET=your-32-char-secret-here
SUPER_ADMIN_CODE=your-super-admin-code
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

### 3. Setup Database

```bash
# Create database
mysql -u root -p -e "CREATE DATABASE auth_service;"

# Run migrations
make migrate-up

# Or manually
mysql -u root -p auth_service < migrations/001_create_users_table.sql
mysql -u root -p auth_service < migrations/002_add_role_column.sql
mysql -u root -p auth_service < migrations/003_create_token_blacklist.sql
mysql -u root -p auth_service < migrations/004_create_password_reset_tokens.sql
```

### 4. Run the Service

```bash
# Development mode
make run

# Or with hot reload (requires air)
make dev

# Build and run
make build
./bin/auth-service
```

---

## 🐳 Docker Deployment

### Build Docker Image

```bash
make docker-build
```

### Run with Docker Compose

```bash
# Start all services (auth-service + MySQL + Redis)
make docker-run

# View logs
make docker-logs

# Stop services
make docker-stop

# Clean up
make docker-clean
```

### Docker Compose includes:
- Auth Service (port 8080)
- MySQL 8.0 (port 3306)
- Redis (port 6379)

---

## ☸️ Kubernetes Deployment

### 1. Build and Push Image

```bash
# Build image
docker build -t your-registry/auth-service:v1.0.0 .

# Push to registry
docker push your-registry/auth-service:v1.0.0
```

### 2. Update Kubernetes Manifests

Edit `k8s/secrets.yaml` with your secrets:
```yaml
stringData:
  mysql-user: "your-user"
  mysql-password: "your-password"
  jwt-secret: "your-jwt-secret-32-chars"
  super-admin-code: "your-super-admin-code"
```

Edit `k8s/configmap.yaml` with your config:
```yaml
data:
  mysql-host: "your-mysql-service"
  cors-allowed-origins: "https://your-domain.com"
```

### 3. Deploy to Kubernetes

```bash
# Deploy all resources
make k8s-deploy

# Check status
make k8s-status

# View logs
make k8s-logs

# Delete deployment
make k8s-delete
```

### Kubernetes Resources Created:
- **Deployment**: 3 replicas with health checks
- **Service**: ClusterIP service on port 80
- **ConfigMap**: Configuration values
- **Secret**: Sensitive data
- **Ingress**: External access with TLS
- **HPA**: Auto-scaling (3-10 pods)

---

## 🔧 Configuration

### CORS Configuration

Allow specific origins:
```env
CORS_ALLOWED_ORIGINS=http://localhost:3000,https://app.yourdomain.com,https://admin.yourdomain.com
```

### Rate Limiting

Default limits (configured in code):
- **Login**: 5 attempts / 15 minutes
- **Register**: 3 attempts / hour
- **Token Validation**: 1000 requests / minute
- **Default**: 100 requests / minute

### JWT Token Expiry

Role-based expiry:
- **User**: 7 days
- **Admin**: 7 days
- **Super Admin**: 30 days

---

## 📊 API Endpoints

### Public Endpoints
```
POST   /v1/auth/register           - Register new user
POST   /v1/auth/login              - Login user
POST   /v1/auth/validate           - Validate JWT token
POST   /v1/auth/forgot-password    - Request password reset
POST   /v1/auth/reset-password     - Reset password with token
```

### Protected Endpoints
```
POST   /v1/auth/refresh            - Refresh JWT token
POST   /v1/auth/logout             - Logout user
POST   /v1/auth/change-password    - Change password
GET    /v1/auth/token-info         - Get token information
GET    /v1/api/profile             - Get user profile
```

### Admin Endpoints
```
GET    /v1/admin/users             - List all users
GET    /v1/admin/users/:id         - Get user by ID
PUT    /v1/admin/users/:id         - Update user
DELETE /v1/admin/users/:id         - Delete user
POST   /v1/admin/users/:id/ban     - Ban user
POST   /v1/admin/users/:id/unban   - Unban user
```

### Super Admin Endpoints
```
POST   /v1/super-admin/promote     - Promote user to admin
```

### Health Endpoints
```
GET    /health                     - Basic health check
GET    /health/ready               - Readiness probe
GET    /health/live                - Liveness probe
```

---

## 🧪 Testing

### Run Tests

```bash
# Run all tests
make test

# Run with coverage
make test-coverage

# Run linter
make lint
```

### API Testing

```bash
# Test health endpoints
make api-test

# Manual testing
# 1. Register user
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# 2. Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# 3. Validate token (for microservices)
curl -X POST http://localhost:8080/v1/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"token":"YOUR_TOKEN_HERE"}'

# 4. Get profile (protected)
curl -X GET http://localhost:8080/v1/api/profile \
  -H "Authorization: Bearer YOUR_TOKEN_HERE"
```

---

## 🔍 Monitoring

### Health Checks

```bash
# Basic health
curl http://localhost:8080/health

# Readiness (checks DB connection)
curl http://localhost:8080/health/ready

# Liveness
curl http://localhost:8080/health/live
```

### Logs

```bash
# Docker logs
docker-compose logs -f auth-service

# Kubernetes logs
kubectl logs -f -l app=auth-service

# Or using make
make docker-logs
make k8s-logs
```

### Metrics

Logs are in JSON format for easy parsing:
```json
{
  "level": "info",
  "timestamp": "2024-01-08T12:00:00Z",
  "msg": "http_request",
  "request_id": "abc-123",
  "method": "POST",
  "path": "/v1/auth/login",
  "status": 200,
  "duration": "45ms",
  "ip": "192.168.1.1"
}
```

---

## 🔐 Security Best Practices

### 1. Environment Variables
- Never commit `.env` to git
- Use strong JWT secrets (32+ characters)
- Rotate secrets regularly

### 2. Database
- Use strong passwords
- Enable SSL/TLS connections
- Regular backups

### 3. Kubernetes Secrets
- Use Sealed Secrets or Vault in production
- Never commit `secrets.yaml` with real values
- Enable RBAC

### 4. Network
- Use HTTPS in production
- Configure firewall rules
- Enable rate limiting

### 5. Monitoring
- Set up alerts for failed logins
- Monitor rate limit hits
- Track token validation failures

---

## 🚨 Troubleshooting

### Service Won't Start

```bash
# Check logs
make docker-logs

# Check database connection
mysql -h localhost -u root -p -e "SELECT 1;"

# Verify environment variables
cat .env
```

### Database Connection Failed

```bash
# Check MySQL is running
docker-compose ps mysql

# Check credentials
mysql -h localhost -u root -p

# Run migrations
make migrate-up
```

### Token Validation Fails

```bash
# Check JWT secret matches
echo $JWT_SECRET

# Verify token format
curl -X POST http://localhost:8080/v1/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"token":"YOUR_TOKEN"}'
```

### Rate Limit Issues

```bash
# Check rate limit headers
curl -v http://localhost:8080/v1/auth/login

# Headers returned:
# X-RateLimit-Limit: 5
# X-RateLimit-Remaining: 3
# X-RateLimit-Reset: 2024-01-08T12:15:00Z
```

---

## 📈 Scaling

### Horizontal Scaling

Kubernetes HPA automatically scales based on:
- CPU usage (70% threshold)
- Memory usage (80% threshold)
- Min replicas: 3
- Max replicas: 10

### Database Scaling

For high traffic:
1. Use read replicas
2. Enable connection pooling
3. Add Redis caching for token validation

### Performance Optimization

1. **Token Validation**: Cache validation results in Redis
2. **Database**: Add indexes on frequently queried fields
3. **Rate Limiting**: Use Redis for distributed rate limiting

---

## 🎯 Production Checklist

- [ ] Update all secrets in `.env` and `k8s/secrets.yaml`
- [ ] Configure CORS for production domains
- [ ] Set up SSL/TLS certificates
- [ ] Enable database backups
- [ ] Configure log aggregation (ELK, Datadog)
- [ ] Set up monitoring and alerts
- [ ] Configure auto-scaling
- [ ] Test disaster recovery
- [ ] Document runbooks
- [ ] Set up CI/CD pipeline

---

## 📚 Additional Resources

### Makefile Commands

```bash
make help              # Show all available commands
make deps              # Install dependencies
make build             # Build application
make run               # Run locally
make test              # Run tests
make docker-build      # Build Docker image
make docker-run        # Run with Docker Compose
make k8s-deploy        # Deploy to Kubernetes
make migrate-up        # Run database migrations
```

### Project Structure

```
go-auth-service/
├── cmd/server/main.go           # Application entry point
├── internal/
│   ├── handlers/                # HTTP handlers
│   ├── services/                # Business logic
│   ├── repository/              # Database operations
│   ├── middleware/              # HTTP middleware
│   ├── routes/                  # Route definitions
│   ├── models/                  # Data models
│   ├── config/                  # Configuration
│   ├── utils/                   # Utilities
│   └── pkg/                     # Shared packages
├── migrations/                  # Database migrations
├── k8s/                         # Kubernetes manifests
├── Dockerfile                   # Docker build file
├── docker-compose.yml           # Docker Compose config
├── Makefile                     # Build commands
└── .env                         # Environment variables
```

---

## 🎉 Success!

Your auth microservice is now:
- ✅ Production-ready
- ✅ Scalable
- ✅ Secure
- ✅ Observable
- ✅ Well-documented

**Ready to integrate with other microservices!**

For questions or issues, check the troubleshooting section or review the implementation guides.
