# Quick Start Guide - Go Auth Service

Get started with the Go Auth Service in 10 minutes!

## Prerequisites

- Go 1.21+
- MySQL 8.0+
- Postman (optional, for API testing)

---

## Step 1: Setup Database (2 minutes)

```bash
# Login to MySQL
mysql -u root -p

# Create database
CREATE DATABASE auth_service;
USE auth_service;

# Run migrations
SOURCE migrations/001_create_users_table.sql;
SOURCE migrations/002_add_role_column.sql;
SOURCE migrations/003_create_token_blacklist.sql;
SOURCE migrations/004_create_password_reset_tokens.sql;
```

---

## Step 2: Configure Environment (1 minute)

```bash
# Copy example env file
cp .env.example .env

# Edit .env file
nano .env
```

**Minimum required configuration:**
```env
APP_PORT=8080
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=your_password
MYSQL_DB=auth_service
JWT_SECRET=your-super-secret-jwt-key-min-32-chars-long
SUPER_ADMIN_CODE=YOUR_SUPER_ADMIN_CODE
```

---

## Step 3: Install Dependencies (1 minute)

```bash
cd go-auth-service
go mod download
go mod tidy
```

---

## Step 4: Run the Service (1 minute)

```bash
# Development mode
go run cmd/server/main.go

# Or build and run
go build -o auth-service cmd/server/main.go
./auth-service
```

**Expected output:**
```
{"level":"info","ts":"2024-01-08T12:00:00Z","msg":"Starting auth service","version":"1.0.0","port":"8080"}
{"level":"info","ts":"2024-01-08T12:00:01Z","msg":"Database connected successfully"}
{"level":"info","ts":"2024-01-08T12:00:01Z","msg":"Server starting","port":"8080"}
```

---

## Step 5: Test the API (5 minutes)

### Using cURL

#### 1. Health Check

```bash
curl http://localhost:8080/health
```

**Expected response:**
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "healthy",
    "version": "1.0.0"
  }
}
```

#### 2. Register User

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Expected response:**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "role": "user"
  }
}
```

#### 3. Login

```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type": application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

**Expected response:**
```json
{
  "success": true,
  "message": "Login successful",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "email": "user@example.com",
      "role": "user"
    },
    "expires_in": "7 days"
  }
}
```

#### 4. Get Profile (Protected Route)

```bash
# Save token from login response
TOKEN="your_jwt_token_here"

curl http://localhost:8080/v1/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

**Expected response:**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "role": "user",
    "is_verified": false
  }
}
```

#### 5. Validate Token (For Microservices)

```bash
curl -X POST http://localhost:8080/v1/auth/validate \
  -H "Content-Type: application/json" \
  -d "{\"token\":\"$TOKEN\"}"
```

**Expected response:**
```json
{
  "success": true,
  "message": "Token is valid",
  "data": {
    "valid": true,
    "user_id": 1,
    "email": "user@example.com",
    "role": "user",
    "permissions": ["users:read"]
  }
}
```

---

### Using Postman

1. **Import Collection**
   - Open Postman
   - Click "Import"
   - Select `docs/postman/Go_Auth_Service.postman_collection.json`

2. **Setup Environment**
   - Click "Environments"
   - Create new environment: "Local"
   - Add variable: `base_url` = `http://localhost:8080`

3. **Run Requests**
   - Select "Local" environment
   - Open "Authentication" folder
   - Run "Register User"
   - Run "Login User" (token auto-saved)
   - Run "Get Profile" (uses saved token)

---

## Step 6: Create Super Admin (Optional)

```bash
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "adminpass123",
    "super_admin_code": "YOUR_SUPER_ADMIN_CODE"
  }'
```

**Login as admin:**
```bash
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@example.com",
    "password": "adminpass123"
  }'
```

**Access admin dashboard:**
```bash
ADMIN_TOKEN="admin_jwt_token_here"

curl http://localhost:8080/v1/admin/dashboard \
  -H "Authorization: Bearer $ADMIN_TOKEN"
```

---

## Common Use Cases

### Use Case 1: User Registration and Login

```bash
# 1. Register
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"newuser@example.com","password":"pass123"}'

# 2. Login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"newuser@example.com","password":"pass123"}'

# 3. Save token and use in subsequent requests
```

### Use Case 2: Password Reset

```bash
# 1. Request password reset
curl -X POST http://localhost:8080/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com"}'

# 2. Reset password with token (from email)
curl -X POST http://localhost:8080/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"token":"RESET_TOKEN","new_password":"newpass123"}'
```

### Use Case 3: Token Validation (Microservices)

```bash
# Your microservice validates incoming tokens
curl -X POST http://localhost:8080/v1/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"token":"USER_JWT_TOKEN"}'

# Check response.data.valid and response.data.role
# Proceed with request if valid
```

---

## Docker Quick Start

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# Check logs
docker-compose logs -f auth-service

# Stop services
docker-compose down
```

### Using Docker

```bash
# Build image
docker build -t auth-service:latest .

# Run container
docker run -d \
  -p 8080:8080 \
  -e MYSQL_HOST=host.docker.internal \
  -e MYSQL_USER=root \
  -e MYSQL_PASSWORD=your_password \
  -e MYSQL_DB=auth_service \
  -e JWT_SECRET=your-secret-key \
  --name auth-service \
  auth-service:latest

# Check logs
docker logs -f auth-service
```

---

## Kubernetes Quick Start

```bash
# Apply configurations
kubectl apply -f k8s/configmap.yaml
kubectl apply -f k8s/secrets.yaml
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml

# Check status
kubectl get pods -l app=auth-service
kubectl get svc auth-service

# View logs
kubectl logs -f -l app=auth-service

# Port forward for local access
kubectl port-forward svc/auth-service 8080:8080
```

---

## Testing

### Run Tests

```bash
# All tests
go test ./tests/... -v

# With coverage
go test ./tests/... -v -cover

# Specific test
go test ./tests/... -v -run TestRegisterUser
```

### Using Test Scripts

```bash
# Linux/Mac
chmod +x tests/run_tests.sh
./tests/run_tests.sh

# Windows
tests\run_tests.bat
```

---

## Troubleshooting

### Issue: Cannot connect to database

**Solution:**
```bash
# Check MySQL is running
mysql -u root -p -e "SELECT 1"

# Verify database exists
mysql -u root -p -e "SHOW DATABASES LIKE 'auth_service'"

# Check .env configuration
cat .env | grep MYSQL
```

### Issue: Port already in use

**Solution:**
```bash
# Find process using port 8080
# Linux/Mac
lsof -i :8080

# Windows
netstat -ano | findstr :8080

# Change port in .env
APP_PORT=8081
```

### Issue: JWT token invalid

**Solution:**
```bash
# Check JWT_SECRET is set
echo $JWT_SECRET

# Verify token format
# Should be: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...

# Check token hasn't expired (7 days for users)
```

---

## Next Steps

1. **Read Documentation**
   - [API Documentation](API_DOCUMENTATION.md)
   - [Postman Guide](POSTMAN_GUIDE.md)
   - [Integration Guide](INTEGRATION_GUIDE.md)

2. **Integrate with Your Service**
   - Follow [Integration Guide](INTEGRATION_GUIDE.md)
   - Use token validation endpoint
   - Implement middleware

3. **Deploy to Production**
   - Review [Deployment Guide](../DEPLOYMENT_GUIDE.md)
   - Configure environment variables
   - Setup monitoring

4. **Customize**
   - Add custom roles
   - Implement email service
   - Add OAuth providers

---

## Useful Commands

```bash
# Development
go run cmd/server/main.go

# Build
go build -o auth-service cmd/server/main.go

# Test
go test ./tests/... -v

# Format code
go fmt ./...

# Lint
golangci-lint run ./...

# Database migrations
make migrate-up
make migrate-down

# Docker
docker-compose up -d
docker-compose logs -f
docker-compose down

# Kubernetes
kubectl apply -f k8s/
kubectl get pods
kubectl logs -f <pod-name>
```

---

## API Endpoints Summary

| Endpoint | Method | Auth | Description |
|----------|--------|------|-------------|
| /health | GET | No | Health check |
| /v1/auth/register | POST | No | Register user |
| /v1/auth/login | POST | No | Login user |
| /v1/auth/logout | POST | Yes | Logout user |
| /v1/auth/validate | POST | No | Validate token |
| /v1/auth/refresh | POST | Yes | Refresh token |
| /v1/api/profile | GET | Yes | Get profile |
| /v1/admin/users | GET | Admin | List users |
| /v1/admin/dashboard | GET | Admin | Admin dashboard |

---

## Support

Need help?
- Check [API Documentation](API_DOCUMENTATION.md)
- Review [Troubleshooting](#troubleshooting)
- Check logs: `docker-compose logs -f auth-service`
- Test with Postman collection

---

**Congratulations! 🎉**

Your Go Auth Service is now running and ready to use!
