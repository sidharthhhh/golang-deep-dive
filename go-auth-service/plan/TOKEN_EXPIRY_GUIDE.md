# Token Expiry Feature Guide

## Overview
The authentication service now implements role-based token expiry with different durations for each role.

## Token Expiry Durations

| Role | Expiry Duration |
|------|----------------|
| User | 7 days |
| Admin | 7 days |
| Super Admin | 30 days |

## Features

### 1. Automatic Token Expiry
- Tokens are automatically generated with role-specific expiry times
- JWT tokens include expiration timestamp in claims
- Middleware validates token expiry on each request

### 2. Login Response Includes Expiry Info
When users login, they receive:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "email": "user@example.com",
  "role": "user",
  "expires_in": "7 days",
  "message": "Login successful"
}
```

### 3. Token Refresh Endpoint
Users can refresh their token before it expires using the `/auth/refresh` endpoint.

## API Usage

### 1. Register User (7-day token)
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

### 2. Register Admin (requires promotion by super admin, 7-day token)
```bash
# First register as user
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@example.com","password":"password123"}'

# Then super admin promotes to admin
curl -X POST http://localhost:8080/super-admin/promote \
  -H "Authorization: Bearer SUPER_ADMIN_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"user_id":2}'
```

### 3. Register Super Admin (30-day token)
```bash
curl -X POST http://localhost:8080/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "email":"superadmin@example.com",
    "password":"password123",
    "super_admin_code":"SUPER_SECRET_ADMIN_CODE_2024"
  }'
```

### 4. Login (returns token with expiry info)
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Response:
```json
{
  "email": "user@example.com",
  "expires_in": "7 days",
  "message": "Login successful",
  "role": "user",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

### 5. Refresh Token (before expiry)
```bash
curl -X POST http://localhost:8080/auth/refresh \
  -H "Authorization: Bearer YOUR_CURRENT_TOKEN"
```

Response:
```json
{
  "email": "user@example.com",
  "expires_in": "7 days",
  "message": "Token refreshed successfully",
  "role": "user",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

## Token Validation

### Automatic Validation
All protected routes automatically validate:
1. Token signature
2. Token expiry
3. Token format

### Expired Token Response
When a token expires, the API returns:
```json
{
  "error": "invalid or expired token"
}
```
Status Code: 401 Unauthorized

## Best Practices

### For Users (7-day tokens)
- Refresh token every 6 days to avoid expiry
- Store token securely (localStorage/sessionStorage)
- Handle 401 errors by redirecting to login

### For Admins (7-day tokens)
- Same as users
- Refresh before performing critical operations

### For Super Admins (30-day tokens)
- Longer expiry for convenience
- Still recommended to refresh periodically
- More secure due to limited super admin accounts

## Security Considerations

1. **Token Storage**: Never store tokens in plain text
2. **HTTPS Only**: Always use HTTPS in production
3. **Refresh Strategy**: Implement automatic refresh before expiry
4. **Logout**: Clear tokens on logout
5. **Super Admin Code**: Keep `SUPER_ADMIN_CODE` secret and rotate regularly

## Implementation Details

### JWT Claims Structure
```go
{
  "user_id": 1,
  "email": "user@example.com",
  "role": "user",
  "exp": 1234567890,  // Expiry timestamp
  "iat": 1234567890,  // Issued at
  "nbf": 1234567890   // Not before
}
```

### Token Generation
```go
// Automatically sets expiry based on role
token, err := utils.GenerateToken(userID, email, role, jwtSecret)
```

### Token Validation
```go
// Validates signature, expiry, and format
claims, err := utils.ValidateToken(tokenString, jwtSecret)
```

## Troubleshooting

### Token Expired
**Problem**: Getting 401 errors
**Solution**: Login again or use refresh endpoint

### Token Invalid
**Problem**: Token not accepted
**Solution**: Check token format, ensure "Bearer " prefix

### Refresh Failed
**Problem**: Cannot refresh token
**Solution**: Token might be expired, login again

## Configuration

Token expiry is configured in `internal/utils/jwt.go`:
```go
case "user":
    expiryDuration = 7 * 24 * time.Hour
case "admin":
    expiryDuration = 7 * 24 * time.Hour
case "super_admin":
    expiryDuration = 30 * 24 * time.Hour
```

To modify expiry durations, update these values and restart the service.
