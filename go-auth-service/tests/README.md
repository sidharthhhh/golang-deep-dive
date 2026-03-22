# Go Auth Service - API Tests

Comprehensive test suite for all API endpoints in the Go Auth Service.

## Test Structure

```
tests/
├── setup_test.go          # Test environment setup and teardown
├── auth_test.go           # Authentication tests (register, login, logout, validate)
├── health_test.go         # Health check endpoint tests
├── password_test.go       # Password management tests (forgot, reset, change)
├── user_test.go           # User profile and management tests
├── admin_test.go          # Admin-specific endpoint tests
├── integration_test.go    # End-to-end workflow tests
├── run_tests.sh           # Linux/Mac test runner script
├── run_tests.bat          # Windows test runner script
└── README.md              # This file
```

## Prerequisites

1. **MySQL Database**: Ensure MySQL is running and accessible
2. **Test Database**: Create a test database named `auth_service_test`
3. **Environment Variables**: Set up test environment variables (handled by setup_test.go)
4. **Dependencies**: Install required Go packages

```bash
go get github.com/stretchr/testify/assert
go mod tidy
```

## Running Tests

### All Tests

```bash
# Linux/Mac
./tests/run_tests.sh

# Windows
tests\run_tests.bat

# Or directly with go
cd go-auth-service
go test ./tests/... -v
```

### Specific Test Files

```bash
# Auth tests only
go test ./tests/auth_test.go ./tests/setup_test.go -v

# Health tests only
go test ./tests/health_test.go ./tests/setup_test.go -v

# Integration tests only
go test ./tests/integration_test.go ./tests/setup_test.go -v
```

### Specific Test Functions

```bash
# Run a specific test
go test ./tests/... -v -run TestRegisterUser

# Run tests matching a pattern
go test ./tests/... -v -run "TestLogin.*"
```

### With Coverage

```bash
go test ./tests/... -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Test Categories

### 1. Authentication Tests (`auth_test.go`)

Tests for user authentication workflows:

- **TestRegisterUser**: User registration with various scenarios
  - Valid registration
  - Duplicate email
  - Invalid email format
  - Missing password

- **TestLogin**: User login functionality
  - Valid login
  - Invalid password
  - Non-existent user

- **TestLogout**: User logout functionality
  - Valid logout with token
  - Logout without token

- **TestTokenValidation**: Token validation endpoint
  - Valid token
  - Invalid token

### 2. Health Check Tests (`health_test.go`)

Tests for service health endpoints:

- **TestHealthCheck**: Basic health check (`/health`)
- **TestHealthReady**: Readiness probe (`/health/ready`)
- **TestHealthLive**: Liveness probe (`/health/live`)

### 3. Password Management Tests (`password_test.go`)

Tests for password-related operations:

- **TestForgotPassword**: Forgot password flow
  - Valid email
  - Non-existent email
  - Invalid email format

- **TestChangePassword**: Change password functionality
  - Valid password change
  - Wrong old password
  - No authentication token

### 4. User Management Tests (`user_test.go`)

Tests for user profile and management:

- **TestGetProfile**: Get user profile
  - With valid token
  - Without token

- **TestListUsers**: List all users (admin only)
  - Admin can list users
  - Unauthorized access

- **TestGetUserByID**: Get specific user by ID
  - Admin can get user
  - Unauthorized access

### 5. Admin Tests (`admin_test.go`)

Tests for admin-specific operations:

- **TestPromoteToAdmin**: Promote user to admin role
  - Super admin can promote
  - Unauthorized access

- **TestUpdateUser**: Update user information
  - Admin can update user
  - Unauthorized access

- **TestDeleteUser**: Delete user
  - Admin can delete user
  - Unauthorized access

- **TestAdminDashboard**: Admin dashboard access
  - Admin can access
  - Regular user cannot access
  - Unauthorized access

### 6. Integration Tests (`integration_test.go`)

End-to-end workflow tests:

- **TestCompleteUserFlow**: Complete user journey
  1. Register
  2. Login
  3. Get profile
  4. Validate token
  5. Change password
  6. Login with new password
  7. Logout

- **TestCompleteAdminFlow**: Complete admin journey
  1. Register super admin
  2. Login as super admin
  3. Register regular user
  4. List users
  5. Access dashboard
  6. Promote user to admin

- **TestRateLimiting**: Rate limiting functionality
- **TestCORSHeaders**: CORS headers presence
- **TestRequestIDHeader**: Request ID tracking

## Test Database Setup

The tests use a separate test database to avoid affecting production data.

### Create Test Database

```sql
CREATE DATABASE auth_service_test;
USE auth_service_test;

-- Run all migrations
SOURCE migrations/001_create_users_table.sql;
SOURCE migrations/002_add_role_column.sql;
SOURCE migrations/003_create_token_blacklist.sql;
SOURCE migrations/004_create_password_reset_tokens.sql;
```

### Automatic Cleanup

The test suite automatically cleans up test data after each test run:
- Deletes all password reset tokens
- Deletes all blacklisted tokens
- Deletes all test users

## Environment Variables

Test environment variables are set in `setup_test.go`:

```go
APP_ENV=test
APP_PORT=8080
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_USER=root
MYSQL_PASSWORD=hjkl
MYSQL_DB=auth_service_test
JWT_SECRET=test-secret-key-for-testing-purposes-only
SUPER_ADMIN_CODE=TEST_SUPER_ADMIN_CODE
CORS_ALLOWED_ORIGINS=http://localhost:3000
```

## Helper Functions

### Creating Test Users

```go
// Create admin user and get token
adminToken := createAdminUser(t)

// Create regular user
createRegularUser(t, "user@example.com")
```

## Expected Test Results

All tests should pass with proper database setup:

```
PASS: TestRegisterUser
PASS: TestLogin
PASS: TestLogout
PASS: TestTokenValidation
PASS: TestHealthCheck
PASS: TestHealthReady
PASS: TestHealthLive
PASS: TestForgotPassword
PASS: TestChangePassword
PASS: TestGetProfile
PASS: TestListUsers
PASS: TestGetUserByID
PASS: TestPromoteToAdmin
PASS: TestUpdateUser
PASS: TestDeleteUser
PASS: TestAdminDashboard
PASS: TestCompleteUserFlow
PASS: TestCompleteAdminFlow
PASS: TestRateLimiting
PASS: TestCORSHeaders
PASS: TestRequestIDHeader
```

## Troubleshooting

### Database Connection Issues

```bash
# Check MySQL is running
mysql -u root -p -e "SELECT 1"

# Verify test database exists
mysql -u root -p -e "SHOW DATABASES LIKE 'auth_service_test'"
```

### Test Failures

1. **Clean test database**: Drop and recreate the test database
2. **Check migrations**: Ensure all migrations are applied
3. **Verify environment**: Check environment variables in setup_test.go
4. **Run individual tests**: Isolate failing tests for debugging

### Port Conflicts

If port 8080 is in use, update `APP_PORT` in `setup_test.go`.

## CI/CD Integration

### GitHub Actions Example

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
        with:
          go-version: 1.21
      - run: go mod download
      - run: go test ./tests/... -v -cover
```

## Best Practices

1. **Isolation**: Each test should be independent
2. **Cleanup**: Always clean up test data
3. **Assertions**: Use clear, descriptive assertions
4. **Coverage**: Aim for >80% code coverage
5. **Speed**: Keep tests fast (<5 seconds per test)
6. **Readability**: Write clear test names and comments

## Contributing

When adding new endpoints:

1. Add corresponding test cases
2. Update this README
3. Ensure all tests pass
4. Maintain test coverage above 80%

## Support

For issues or questions:
- Check test output for detailed error messages
- Review test database state
- Verify environment configuration
- Check service logs
