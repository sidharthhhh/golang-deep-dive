# Testing Guide - Go Auth Service

Complete guide for testing the Go Auth Service API.

## Table of Contents

1. [Overview](#overview)
2. [Test Setup](#test-setup)
3. [Running Tests](#running-tests)
4. [Test Coverage](#test-coverage)
5. [API Endpoint Tests](#api-endpoint-tests)
6. [Manual Testing](#manual-testing)
7. [Troubleshooting](#troubleshooting)

---

## Overview

The test suite provides comprehensive coverage of all API endpoints, including:

- ✅ Authentication (register, login, logout)
- ✅ Token validation and management
- ✅ Password management (forgot, reset, change)
- ✅ User profile and management
- ✅ Admin operations
- ✅ Health checks
- ✅ Integration workflows
- ✅ Middleware (CORS, rate limiting, request ID)

### Test Statistics

- **Total Test Files**: 7
- **Total Test Cases**: 20+
- **Test Scenarios**: 50+
- **Expected Coverage**: >80%

---

## Test Setup

### 1. Prerequisites

```bash
# Install Go (1.21 or higher)
go version

# Install MySQL
mysql --version

# Install test dependencies
go get github.com/stretchr/testify/assert
go mod tidy
```

### 2. Database Setup

```bash
# Create test database
mysql -u root -p

CREATE DATABASE auth_service_test;
USE auth_service_test;

# Run migrations
SOURCE migrations/001_create_users_table.sql;
SOURCE migrations/002_add_role_column.sql;
SOURCE migrations/003_create_token_blacklist.sql;
SOURCE migrations/004_create_password_reset_tokens.sql;
```

### 3. Environment Configuration

Test environment is configured in `tests/setup_test.go`:

```go
MYSQL_DB=auth_service_test
JWT_SECRET=test-secret-key-for-testing-purposes-only
SUPER_ADMIN_CODE=TEST_SUPER_ADMIN_CODE
```

---

## Running Tests

### Quick Start

```bash
# Linux/Mac
chmod +x tests/run_tests.sh
./tests/run_tests.sh

# Windows
tests\run_tests.bat
```

### All Tests

```bash
cd go-auth-service
go test ./tests/... -v
```

### Specific Test Files

```bash
# Authentication tests
go test ./tests/auth_test.go ./tests/setup_test.go -v

# Health check tests
go test ./tests/health_test.go ./tests/setup_test.go -v

# Password tests
go test ./tests/password_test.go ./tests/setup_test.go -v

# User management tests
go test ./tests/user_test.go ./tests/setup_test.go -v

# Admin tests
go test ./tests/admin_test.go ./tests/setup_test.go -v

# Integration tests
go test ./tests/integration_test.go ./tests/setup_test.go -v
```

### Specific Test Functions

```bash
# Run single test
go test ./tests/... -v -run TestRegisterUser

# Run tests matching pattern
go test ./tests/... -v -run "TestLogin.*"

# Run tests in specific file
go test ./tests/... -v -run "TestComplete.*"
```

### With Coverage

```bash
# Generate coverage report
go test ./tests/... -v -cover -coverprofile=coverage.out

# View coverage in browser
go tool cover -html=coverage.out -o coverage.html
open coverage.html  # Mac
start coverage.html # Windows
xdg-open coverage.html # Linux
```

### Parallel Execution

```bash
# Run tests in parallel
go test ./tests/... -v -parallel 4
```

---

## Test Coverage

### Coverage by Module

| Module | Coverage | Status |
|--------|----------|--------|
| Handlers | >85% | ✅ |
| Services | >80% | ✅ |
| Repositories | >75% | ✅ |
| Middleware | >90% | ✅ |
| Utils | >85% | ✅ |

### Viewing Coverage

```bash
# Terminal output
go test ./tests/... -cover

# Detailed HTML report
go test ./tests/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```

---

## API Endpoint Tests

### Authentication Endpoints

#### POST /v1/auth/register

```bash
# Test: Valid registration
curl -X POST http://localhost:8080/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Expected: 201 Created
# Response: {"success":true,"data":{...}}
```

#### POST /v1/auth/login

```bash
# Test: Valid login
curl -X POST http://localhost:8080/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Expected: 200 OK
# Response: {"success":true,"data":{"token":"..."}}
```

#### POST /v1/auth/logout

```bash
# Test: Valid logout
curl -X POST http://localhost:8080/v1/auth/logout \
  -H "Authorization: Bearer YOUR_TOKEN"

# Expected: 200 OK
# Response: {"success":true,"message":"Logged out successfully"}
```

#### POST /v1/auth/validate

```bash
# Test: Token validation
curl -X POST http://localhost:8080/v1/auth/validate \
  -H "Content-Type: application/json" \
  -d '{"token":"YOUR_TOKEN"}'

# Expected: 200 OK
# Response: {"success":true,"data":{"valid":true,...}}
```

### Health Check Endpoints

#### GET /health

```bash
curl http://localhost:8080/health

# Expected: 200 OK
# Response: {"success":true,"message":"Service is healthy"}
```

#### GET /health/ready

```bash
curl http://localhost:8080/health/ready

# Expected: 200 OK
# Response: {"success":true,"data":{"status":"ready","database":"healthy"}}
```

#### GET /health/live

```bash
curl http://localhost:8080/health/live

# Expected: 200 OK
# Response: {"success":true,"data":{"status":"alive"}}
```

### Password Management Endpoints

#### POST /v1/auth/forgot-password

```bash
curl -X POST http://localhost:8080/v1/auth/forgot-password \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com"}'

# Expected: 200 OK
```

#### POST /v1/auth/reset-password

```bash
curl -X POST http://localhost:8080/v1/auth/reset-password \
  -H "Content-Type: application/json" \
  -d '{"token":"RESET_TOKEN","new_password":"newpass123"}'

# Expected: 200 OK
```

#### POST /v1/auth/change-password

```bash
curl -X POST http://localhost:8080/v1/auth/change-password \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{"old_password":"oldpass","new_password":"newpass"}'

# Expected: 200 OK
```

### User Management Endpoints

#### GET /v1/api/profile

```bash
curl http://localhost:8080/v1/api/profile \
  -H "Authorization: Bearer YOUR_TOKEN"

# Expected: 200 OK
# Response: {"success":true,"data":{"email":"...","role":"..."}}
```

#### GET /v1/admin/users

```bash
curl http://localhost:8080/v1/admin/users \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Expected: 200 OK
# Response: {"success":true,"data":{"users":[...],"total":10}}
```

#### GET /v1/admin/users/:id

```bash
curl http://localhost:8080/v1/admin/users/1 \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Expected: 200 OK
```

#### PUT /v1/admin/users/:id

```bash
curl -X PUT http://localhost:8080/v1/admin/users/1 \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer ADMIN_TOKEN" \
  -d '{"email":"updated@example.com"}'

# Expected: 200 OK
```

#### DELETE /v1/admin/users/:id

```bash
curl -X DELETE http://localhost:8080/v1/admin/users/1 \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Expected: 200 OK
```

### Admin Endpoints

#### GET /v1/admin/dashboard

```bash
curl http://localhost:8080/v1/admin/dashboard \
  -H "Authorization: Bearer ADMIN_TOKEN"

# Expected: 200 OK
```

#### POST /v1/super-admin/promote

```bash
curl -X POST http://localhost:8080/v1/super-admin/promote \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer SUPER_ADMIN_TOKEN" \
  -d '{"user_id":2}'

# Expected: 200 OK
```

---

## Manual Testing

### Using Postman

1. Import the API collection (create from endpoints above)
2. Set environment variables:
   - `base_url`: http://localhost:8080
   - `token`: (set after login)
3. Run collection tests

### Using cURL Scripts

Create a test script:

```bash
#!/bin/bash

BASE_URL="http://localhost:8080"

# Register
echo "Testing registration..."
curl -X POST $BASE_URL/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Login
echo "Testing login..."
TOKEN=$(curl -s -X POST $BASE_URL/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}' \
  | jq -r '.data.token')

echo "Token: $TOKEN"

# Get Profile
echo "Testing profile..."
curl $BASE_URL/v1/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

---

## Troubleshooting

### Common Issues

#### 1. Database Connection Failed

```bash
# Check MySQL is running
mysql -u root -p -e "SELECT 1"

# Verify test database exists
mysql -u root -p -e "SHOW DATABASES LIKE 'auth_service_test'"

# Recreate test database
mysql -u root -p < migrations/001_create_users_table.sql
```

#### 2. Tests Failing

```bash
# Clean test database
mysql -u root -p auth_service_test -e "
  DELETE FROM password_reset_tokens;
  DELETE FROM token_blacklist;
  DELETE FROM users;
"

# Run single test for debugging
go test ./tests/... -v -run TestRegisterUser
```

#### 3. Port Already in Use

```bash
# Find process using port 8080
# Linux/Mac
lsof -i :8080

# Windows
netstat -ano | findstr :8080

# Kill the process or change port in setup_test.go
```

#### 4. Import Errors

```bash
# Update dependencies
go mod tidy
go get github.com/stretchr/testify/assert

# Clear module cache
go clean -modcache
go mod download
```

### Debug Mode

```bash
# Run with verbose output
go test ./tests/... -v -test.v

# Run with race detector
go test ./tests/... -race

# Run with timeout
go test ./tests/... -timeout 30s
```

---

## Best Practices

1. **Isolation**: Each test should be independent
2. **Cleanup**: Always clean up test data
3. **Assertions**: Use descriptive assertions
4. **Coverage**: Maintain >80% coverage
5. **Speed**: Keep tests fast (<5s per test)
6. **Documentation**: Document complex test scenarios

---

## CI/CD Integration

### GitHub Actions

```yaml
name: Tests
on: [push, pull_request]
jobs:
  test:
    runs-on: ubuntu-latest
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: hjkl
          MYSQL_DATABASE: auth_service_test
        ports:
          - 3306:3306
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
      - run: go test ./tests/... -v -cover
```

---

## Support

For issues or questions:
- Check test output for error messages
- Review database state
- Verify environment configuration
- Check service logs in `logs/` directory
