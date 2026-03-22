# Logout Feature Guide

## Overview
The authentication service now includes a complete logout feature with token blacklisting. When users logout, their JWT tokens are added to a blacklist, preventing further use even if the token hasn't expired.

## How It Works

### 1. Token Blacklisting
- Each JWT token has a unique ID (JTI - JWT ID)
- On logout, the token's JTI is added to the `token_blacklist` table
- Middleware checks the blacklist on every request
- Blacklisted tokens are rejected with 401 Unauthorized

### 2. Automatic Cleanup
- Expired tokens are automatically cleaned from the blacklist
- Reduces database size and improves performance
- Can be run as a scheduled job

## Database Schema

```sql
CREATE TABLE token_blacklist (
    id INT AUTO_INCREMENT PRIMARY KEY,
    token_jti VARCHAR(255) NOT NULL UNIQUE,
    user_id INT NOT NULL,
    blacklisted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP NOT NULL,
    INDEX idx_token_jti (token_jti),
    INDEX idx_user_id (user_id),
    INDEX idx_expires_at (expires_at)
);
```

## API Usage

### 1. Login (Get Token)
```bash
curl -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}'
```

Response:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "email": "user@example.com",
  "role": "user",
  "expires_in": "7 days",
  "message": "Login successful"
}
```

### 2. Use Token (Access Protected Routes)
```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_TOKEN"
```

### 3. Logout (Blacklist Token)
```bash
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer YOUR_TOKEN"
```

Response:
```json
{
  "message": "Logged out successfully"
}
```

### 4. Try Using Logged Out Token (Will Fail)
```bash
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer LOGGED_OUT_TOKEN"
```

Response:
```json
{
  "error": "token has been revoked"
}
```
Status Code: 401 Unauthorized

## Features

### ✅ Secure Logout
- Tokens are immediately invalidated
- Cannot be reused after logout
- Prevents token theft attacks

### ✅ Token Tracking
- Each token has a unique JTI
- Tracks which user owns which token
- Stores token expiry for cleanup

### ✅ Automatic Cleanup
- Expired tokens are removed from blacklist
- Keeps database clean and performant
- Can be scheduled as a cron job

### ✅ Multi-Device Support
- Each login creates a new token with unique JTI
- Logout only affects the specific token used
- Users can be logged in on multiple devices

## Implementation Details

### JWT Token Structure
```json
{
  "user_id": 1,
  "email": "user@example.com",
  "role": "user",
  "jti": "a1b2c3d4e5f6...",  // Unique token ID
  "exp": 1234567890,
  "iat": 1234567890,
  "nbf": 1234567890
}
```

### Logout Flow
1. User sends logout request with token
2. Server extracts JTI from token
3. JTI is added to blacklist with expiry time
4. Server responds with success message
5. Future requests with that token are rejected

### Middleware Check
```go
// On every protected route request:
1. Extract token from Authorization header
2. Validate token signature and expiry
3. Check if token JTI is in blacklist
4. If blacklisted → reject with 401
5. If not blacklisted → allow request
```

## Security Considerations

### 1. Token Theft Protection
- Even if a token is stolen, user can logout to invalidate it
- Stolen tokens become useless after logout

### 2. Session Management
- Users can logout from specific devices
- Each device has its own token
- Logout doesn't affect other devices

### 3. Token Expiry
- Tokens still expire based on role (7/30 days)
- Blacklist entries are cleaned after expiry
- No need to store expired tokens forever

## Best Practices

### For Frontend Applications

1. **Store Token Securely**
```javascript
// Store in localStorage or sessionStorage
localStorage.setItem('auth_token', token);
```

2. **Clear Token on Logout**
```javascript
async function logout() {
  const token = localStorage.getItem('auth_token');
  
  // Call logout API
  await fetch('http://localhost:8080/auth/logout', {
    method: 'POST',
    headers: {
      'Authorization': `Bearer ${token}`
    }
  });
  
  // Clear local storage
  localStorage.removeItem('auth_token');
  
  // Redirect to login
  window.location.href = '/login';
}
```

3. **Handle Revoked Tokens**
```javascript
// Intercept 401 errors
axios.interceptors.response.use(
  response => response,
  error => {
    if (error.response.status === 401) {
      // Token expired or revoked
      localStorage.removeItem('auth_token');
      window.location.href = '/login';
    }
    return Promise.reject(error);
  }
);
```

### For Backend

1. **Run Cleanup Periodically**
```go
// Add to main.go or separate service
go func() {
    ticker := time.NewTicker(24 * time.Hour)
    for range ticker.C {
        tokenRepo.CleanupExpiredTokens(context.Background())
    }
}()
```

2. **Monitor Blacklist Size**
```sql
SELECT COUNT(*) FROM token_blacklist;
```

3. **Logout All User Sessions**
```sql
-- Add to repository if needed
INSERT INTO token_blacklist (token_jti, user_id, expires_at)
SELECT token_jti, user_id, expires_at 
FROM active_tokens 
WHERE user_id = ?;
```

## Testing

### 1. Test Normal Logout
```bash
# Login
TOKEN=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' \
  | jq -r '.token')

# Use token (should work)
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN"

# Logout
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer $TOKEN"

# Try using token again (should fail)
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN"
```

### 2. Test Multi-Device
```bash
# Login from device 1
TOKEN1=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' \
  | jq -r '.token')

# Login from device 2
TOKEN2=$(curl -s -X POST http://localhost:8080/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"user@example.com","password":"password123"}' \
  | jq -r '.token')

# Logout device 1
curl -X POST http://localhost:8080/auth/logout \
  -H "Authorization: Bearer $TOKEN1"

# Device 1 token should fail
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN1"

# Device 2 token should still work
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer $TOKEN2"
```

## Troubleshooting

### Token Still Works After Logout
**Problem**: Token not being blacklisted
**Solution**: 
- Check database connection
- Verify token_blacklist table exists
- Check JTI is being extracted correctly

### Performance Issues
**Problem**: Slow authentication
**Solution**:
- Run cleanup to remove expired tokens
- Add database indexes (already included)
- Consider Redis for blacklist (future enhancement)

### Logout Fails
**Problem**: Cannot logout
**Solution**:
- Verify token is valid
- Check Authorization header format
- Ensure token_blacklist table has write permissions

## Future Enhancements

1. **Redis Blacklist**: Use Redis for faster lookups
2. **Logout All Devices**: Add endpoint to logout from all devices
3. **Session Management**: View and manage active sessions
4. **Suspicious Activity**: Auto-logout on suspicious activity
5. **Token Refresh on Logout**: Optionally refresh instead of logout

## Migration

Run the migration to create the blacklist table:
```bash
mysql -u root -p auth_service < migrations/003_create_token_blacklist.sql
```

## Summary

The logout feature provides:
- ✅ Secure token invalidation
- ✅ Multi-device support
- ✅ Automatic cleanup
- ✅ Token theft protection
- ✅ Simple API integration

Users can now safely logout, and their tokens are immediately invalidated!
