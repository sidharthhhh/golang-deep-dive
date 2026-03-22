# API Documentation - Go Auth Service

Complete API reference for the Go Auth Service.

## Base URL

```
http://localhost:8080
```

## Table of Contents

1. [Authentication](#authentication)
2. [Health Checks](#health-checks)
3. [Password Management](#password-management)
4. [User Management](#user-management)
5. [Admin Operations](#admin-operations)
6. [Response Format](#response-format)
7. [Error Codes](#error-codes)

---

## Authentication

### Register User

Create a new user account.

**Endpoint:** `POST /v1/auth/register`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Register as Super Admin:**
```json
{
  "email": "admin@example.com",
  "password": "adminpass123",
  "super_admin_code": "YOUR_SUPER_ADMIN_CODE"
}
```

**Success Response (201 Created):**
```json
{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": 1,
    "email": "user@example.com",
    "role": "user",
    "is_verified": false,
    "created_at": "2024-01-08T12:00:00Z"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

**Error Response (409 Conflict):**
```json
{
  "success": false,
  "error": {
    "code": "CONFLICT",
    "message": "User already exists"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

### Login

Authenticate and receive a JWT token.

**Endpoint:** `POST /v1/auth/login`

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Success Response (200 OK):**
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
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

**Token Expiry:**
- User/Admin: 7 days
- Super Admin: 30 days

---

### Logout

Invalidate the current token.

**Endpoint:** `POST /v1/auth/logout`

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Logged out successfully",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

### Validate Token

Validate a JWT token (for microservices).

**Endpoint:** `POST /v1/auth/validate`

**Request Body:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Token is valid",
  "data": {
    "valid": true,
    "user_id": 1,
    "email": "user@example.com",
    "role": "admin",
    "permissions": [
      "users:read",
      "users:write",
      "admin:read"
    ],
    "expires_at": "2024-01-15T12:00:00Z",
    "issued_at": "2024-01-08T12:00:00Z"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

**Invalid Token Response (401 Unauthorized):**
```json
{
  "success": false,
  "error": {
    "code": "UNAUTHORIZED",
    "message": "Invalid or expired token"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

### Refresh Token

Get a new token before expiry.

**Endpoint:** `POST /v1/auth/refresh`

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Token refreshed successfully",
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "expires_in": "7 days"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

### Get Token Info

Get information about the current token.

**Endpoint:** `GET /v1/auth/token-info`

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "user_id": 1,
    "email": "user@example.com",
    "role": "user",
    "issued_at": "2024-01-08T12:00:00Z",
    "expires_at": "2024-01-15T12:00:00Z"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

## Health Checks

### Basic Health Check

Check if the service is running.

**Endpoint:** `GET /health`

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Service is healthy",
  "data": {
    "status": "healthy",
    "version": "1.0.0"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

### Readiness Probe

Check if the service is ready to accept traffic.

**Endpoint:** `GET /health/ready`

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Service is ready",
  "data": {
    "status": "ready",
    "database": "healthy",
    "version": "1.0.0"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

### Liveness Probe

Check if the service is alive.

**Endpoint:** `GET /health/live`

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Service is alive",
  "data": {
    "status": "alive",
    "uptime": "2h30m15s"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

## Password Management

### Forgot Password

Request a password reset token.

**Endpoint:** `POST /v1/auth/forgot-password`

**Request Body:**
```json
{
  "email": "user@example.com"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Password reset email sent",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

**Note:** For security, this endpoint always returns success even if the email doesn't exist.

---

### Reset Password

Reset password using a reset token.

**Endpoint:** `POST /v1/auth/reset-password`

**Request Body:**
```json
{
  "token": "reset_token_here",
  "new_password": "newpassword123"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Password reset successfully",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

### Change Password

Change password for authenticated user.

**Endpoint:** `POST /v1/auth/change-password`

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Request Body:**
```json
{
  "old_password": "oldpassword123",
  "new_password": "newpassword456"
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "Password changed successfully",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

## User Management

### Get Profile

Get current user's profile.

**Endpoint:** `GET /v1/api/profile`

**Headers:**
```
Authorization: Bearer YOUR_JWT_TOKEN
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "role": "user",
    "is_verified": false,
    "created_at": "2024-01-08T12:00:00Z",
    "updated_at": "2024-01-08T12:00:00Z"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

### List Users (Admin Only)

Get a list of all users.

**Endpoint:** `GET /v1/admin/users`

**Headers:**
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

**Query Parameters:**
- `limit` (optional): Number of users per page (default: 10)
- `offset` (optional): Pagination offset (default: 0)

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "users": [
      {
        "id": 1,
        "email": "user1@example.com",
        "role": "user",
        "is_verified": true,
        "created_at": "2024-01-08T12:00:00Z"
      },
      {
        "id": 2,
        "email": "user2@example.com",
        "role": "admin",
        "is_verified": true,
        "created_at": "2024-01-08T11:00:00Z"
      }
    ],
    "total": 25,
    "limit": 10,
    "offset": 0
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

### Get User by ID (Admin Only)

Get a specific user's details.

**Endpoint:** `GET /v1/admin/users/:id`

**Headers:**
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "id": 1,
    "email": "user@example.com",
    "role": "user",
    "is_verified": true,
    "created_at": "2024-01-08T12:00:00Z",
    "updated_at": "2024-01-08T12:00:00Z"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

### Update User (Admin Only)

Update user information.

**Endpoint:** `PUT /v1/admin/users/:id`

**Headers:**
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "is_verified": true
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "User updated successfully",
  "data": {
    "id": 1,
    "email": "newemail@example.com",
    "role": "user",
    "is_verified": true,
    "updated_at": "2024-01-08T12:30:00Z"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:30:00Z"
}
```

---

### Delete User (Admin Only)

Delete a user account.

**Endpoint:** `DELETE /v1/admin/users/:id`

**Headers:**
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "User deleted successfully",
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

## Admin Operations

### Admin Dashboard

Access admin dashboard.

**Endpoint:** `GET /v1/admin/dashboard`

**Headers:**
```
Authorization: Bearer ADMIN_JWT_TOKEN
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "data": {
    "total_users": 150,
    "active_users": 120,
    "admin_users": 5,
    "recent_registrations": 10
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

### Promote User to Admin (Super Admin Only)

Promote a user to admin role.

**Endpoint:** `POST /v1/super-admin/promote`

**Headers:**
```
Authorization: Bearer SUPER_ADMIN_JWT_TOKEN
```

**Request Body:**
```json
{
  "user_id": 5
}
```

**Success Response (200 OK):**
```json
{
  "success": true,
  "message": "User promoted to admin successfully",
  "data": {
    "user_id": 5,
    "new_role": "admin"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

## Response Format

All API responses follow this standard format:

### Success Response

```json
{
  "success": true,
  "message": "Operation successful",
  "data": { ... },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

### Error Response

```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error description"
  },
  "request_id": "550e8400-e29b-41d4-a716-446655440000",
  "timestamp": "2024-01-08T12:00:00Z"
}
```

---

## Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| BAD_REQUEST | 400 | Invalid request format or parameters |
| UNAUTHORIZED | 401 | Missing or invalid authentication |
| FORBIDDEN | 403 | Insufficient permissions |
| NOT_FOUND | 404 | Resource not found |
| CONFLICT | 409 | Resource already exists |
| RATE_LIMIT_EXCEEDED | 429 | Too many requests |
| INTERNAL_ERROR | 500 | Server error |
| TOKEN_EXPIRED | 401 | JWT token has expired |
| TOKEN_INVALID | 401 | JWT token is invalid |
| TOKEN_REVOKED | 401 | JWT token has been revoked |
| INVALID_PASSWORD | 401 | Password is incorrect |

---

## Rate Limiting

Rate limits are enforced per IP address:

| Endpoint | Limit |
|----------|-------|
| POST /v1/auth/login | 5 requests / 15 minutes |
| POST /v1/auth/register | 3 requests / hour |
| POST /v1/auth/validate | 1000 requests / minute |
| Other endpoints | 100 requests / minute |

### Rate Limit Headers

```
X-RateLimit-Limit: 5
X-RateLimit-Remaining: 3
X-RateLimit-Reset: 2024-01-08T12:15:00Z
```

---

## CORS

CORS is enabled for configured origins. Default allowed origins:
- http://localhost:3000
- http://localhost:3001

Configure via `CORS_ALLOWED_ORIGINS` environment variable.

---

## Request ID

Every request receives a unique ID for tracking:

```
X-Request-ID: 550e8400-e29b-41d4-a716-446655440000
```

Use this ID for debugging and support requests.

---

## Authentication Flow

### User Registration and Login

```
1. POST /v1/auth/register
   → Receive user data

2. POST /v1/auth/login
   → Receive JWT token

3. Use token in Authorization header
   → Authorization: Bearer YOUR_TOKEN

4. POST /v1/auth/logout (when done)
   → Token is blacklisted
```

### Token Validation (Microservices)

```
1. Receive request with token from client

2. POST /v1/auth/validate
   → Send token for validation

3. Check response
   → If valid: proceed with request
   → If invalid: return 401 Unauthorized
```

---

## Best Practices

1. **Store tokens securely** - Use httpOnly cookies or secure storage
2. **Refresh tokens before expiry** - Use /v1/auth/refresh endpoint
3. **Handle rate limits** - Implement exponential backoff
4. **Log request IDs** - Include X-Request-ID in logs
5. **Validate responses** - Always check `success` field
6. **Handle errors gracefully** - Use error codes for specific handling

---

## Support

For API issues or questions:
- Check error codes and messages
- Include request ID in support requests
- Review rate limit headers
- Check authentication token validity
