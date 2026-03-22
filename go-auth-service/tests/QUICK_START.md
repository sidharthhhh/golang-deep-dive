# Quick Start - Testing Guide

Get started with testing the Go Auth Service in 5 minutes!

## Prerequisites

- ✅ Go 1.21+ installed
- ✅ MySQL running
- ✅ Test database created

## Setup (One-Time)

### 1. Create Test Database

```bash
mysql -u root -p
```

```sql
CREATE DATABASE auth_service_test;
USE auth_service_test;

SOURCE migrations/001_create_users_table.sql;
SOURCE migrations/002_add_role_column.sql;
SOURCE migrations/003_create_token_blacklist.sql;
SOURCE migrations/004_create_password_reset_tokens.sql;
```

### 2. Install Dependencies

```bash
cd go-auth-service
go mod download
go get github.com/stretchr/testify/assert
```

## Run Tests

### Option 1: Use Test Scripts (Recommended)

```bash
# Linux/Mac
chmod +x tests/run_tests.sh
./tests/run_tests.sh

# Windows
tests\run_tests.bat
```

### Option 2: Direct Go Command

```bash
go test ./tests/... -v
```

### Option 3: With Coverage

```bash
# Linux/Mac
./tests/run_tests.sh coverage

# Windows
tests\run_tests.bat coverage

# Or directly
go test ./tests/... -v -cover -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

## Expected Output

```
=== RUN   TestRegisterUser
--- PASS: TestRegisterUser (0.15s)
=== RUN   TestLogin
--- PASS: TestLogin (0.12s)
=== RUN   TestLogout
--- PASS: TestLogout (0.10s)
...
PASS
ok      github.com/sidharthhhh/go-auth-service/tests    2.456s
```

## Test Individual Endpoints

```bash
# Auth tests only
go test ./tests/auth_test.go ./tests/setup_test.go -v

# Health tests only
go test ./tests/health_test.go ./tests/setup_test.go -v

# Integration tests only
go test ./tests/integration_test.go ./tests/setup_test.go -v
```

## Troubleshooting

### Database Connection Error

```bash
# Check MySQL is running
mysql -u root -p -e "SELECT 1"

# Verify test database exists
mysql -u root -p -e "SHOW DATABASES LIKE 'auth_service_test'"
```

### Tests Failing

```bash
# Clean test data
mysql -u root -p auth_service_test -e "
  DELETE FROM password_reset_tokens;
  DELETE FROM token_blacklist;
  DELETE FROM users;
"

# Run again
go test ./tests/... -v
```

### Port Already in Use

Edit `tests/setup_test.go` and change `APP_PORT` from 8080 to another port.

## Next Steps

- 📖 Read [TESTING_GUIDE.md](../TESTING_GUIDE.md) for detailed documentation
- 📊 Check [TEST_SUMMARY.md](TEST_SUMMARY.md) for coverage details
- 📝 Review [README.md](README.md) for test structure

## Quick Reference

| Command | Description |
|---------|-------------|
| `go test ./tests/... -v` | Run all tests |
| `go test ./tests/... -cover` | Run with coverage |
| `go test ./tests/... -run TestName` | Run specific test |
| `./tests/run_tests.sh` | Run with script (Linux/Mac) |
| `tests\run_tests.bat` | Run with script (Windows) |

## Test Files

- `setup_test.go` - Test environment setup
- `auth_test.go` - Authentication tests
- `health_test.go` - Health check tests
- `password_test.go` - Password management tests
- `user_test.go` - User management tests
- `admin_test.go` - Admin operation tests
- `integration_test.go` - End-to-end workflow tests

## Support

Having issues? Check:
1. MySQL is running
2. Test database exists and has migrations
3. Dependencies are installed (`go mod tidy`)
4. Port 8080 is available

For more help, see [TESTING_GUIDE.md](../TESTING_GUIDE.md)
