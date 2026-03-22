# Postman Testing Guide - Go Auth Service

Complete guide for testing the Go Auth Service API using Postman.

## Table of Contents

1. [Setup](#setup)
2. [Environment Variables](#environment-variables)
3. [Collection Structure](#collection-structure)
4. [Testing Workflows](#testing-workflows)
5. [Automated Tests](#automated-tests)
6. [Tips and Tricks](#tips-and-tricks)

---

## Setup

### 1. Import Collection

Create a new collection in Postman named "Go Auth Service".

### 2. Set Base URL

Create an environment with the following variables:

| Variable | Initial Value | Current Value |
|----------|---------------|---------------|
| base_url | http://localhost:8080 | http://localhost:8080 |
| token | | (auto-populated) |
| admin_token | | (auto-populated) |
| user_id | | (auto-populated) |

---

## Environment Variables

### Creating Environment

1. Click "Environments" in Postman
2. Click "+" to create new environment
3. Name it "Go Auth Service - Local"
4. Add variables:

```
base_url = http://localhost:8080
token = 
admin_token = 
super_admin_code = YOUR_SUPER_ADMIN_CODE
user_email = test@example.com
user_password = password123
```

### Using Variables

In requests, use `{{variable_name}}`:
```
{{base_url}}/v1/auth/login
```

---

## Collection Structure

### Folder 1: Authentication

#### 1.1 Register User

**Method:** POST  
**URL:** `{{base_url}}/v1/auth/register`  
**Body (JSON):**
```json
{
  "email": "{{user_email}}",
  "password": "{{user_password}}"
}
```

**Tests Script:**
```javascript
// Check status code
pm.test("Status code is 201", function () {
    pm.response.to.have.status(201);
});

// Check response structure
pm.test("Response has success field", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData).to.have.property('success');
    pm.expect(jsonData.success).to.be.true;
});

// Save user ID
pm.test("Save user ID", function () {
    var jsonData = pm.response.json();
    if (jsonData.data && jsonData.data.id) {
        pm.environment.set("user_id", jsonData.data.id);
    }
});
```

---

#### 1.2 Register Super Admin

**Method:** POST  
**URL:** `{{base_url}}/v1/auth/register`  
**Body (JSON):**
```json
{
  "email": "admin@example.com",
  "password": "adminpass123",
  "super_admin_code": "{{super_admin_code}}"
}
```

**Tests Script:**
```javascript
pm.test("Status code is 201", function () {
    pm.response.to.have.status(201);
});

pm.test("User role is super_admin", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data.role).to.eql("super_admin");
});
```

---

#### 1.3 Login User

**Method:** POST  
**URL:** `{{base_url}}/v1/auth/login`  
**Body (JSON):**
```json
{
  "email": "{{user_email}}",
  "password": "{{user_password}}"
}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response contains token", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('token');
    
    // Save token to environment
    pm.environment.set("token", jsonData.data.token);
});

pm.test("Token is not empty", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data.token).to.not.be.empty;
});
```

---

#### 1.4 Login Admin

**Method:** POST  
**URL:** `{{base_url}}/v1/auth/login`  
**Body (JSON):**
```json
{
  "email": "admin@example.com",
  "password": "adminpass123"
}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Save admin token", function () {
    var jsonData = pm.response.json();
    pm.environment.set("admin_token", jsonData.data.token);
});
```

---

#### 1.5 Validate Token

**Method:** POST  
**URL:** `{{base_url}}/v1/auth/validate`  
**Body (JSON):**
```json
{
  "token": "{{token}}"
}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Token is valid", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data.valid).to.be.true;
});

pm.test("Response contains user info", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('user_id');
    pm.expect(jsonData.data).to.have.property('email');
    pm.expect(jsonData.data).to.have.property('role');
});
```

---

#### 1.6 Refresh Token

**Method:** POST  
**URL:** `{{base_url}}/v1/auth/refresh`  
**Headers:**
```
Authorization: Bearer {{token}}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("New token received", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('token');
    
    // Update token
    pm.environment.set("token", jsonData.data.token);
});
```

---

#### 1.7 Get Token Info

**Method:** GET  
**URL:** `{{base_url}}/v1/auth/token-info`  
**Headers:**
```
Authorization: Bearer {{token}}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Token info contains required fields", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('user_id');
    pm.expect(jsonData.data).to.have.property('email');
    pm.expect(jsonData.data).to.have.property('role');
    pm.expect(jsonData.data).to.have.property('issued_at');
    pm.expect(jsonData.data).to.have.property('expires_at');
});
```

---

#### 1.8 Logout

**Method:** POST  
**URL:** `{{base_url}}/v1/auth/logout`  
**Headers:**
```
Authorization: Bearer {{token}}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Logout successful", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.success).to.be.true;
});
```

---

### Folder 2: Health Checks

#### 2.1 Basic Health

**Method:** GET  
**URL:** `{{base_url}}/health`

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Service is healthy", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data.status).to.eql("healthy");
});
```

---

#### 2.2 Readiness Probe

**Method:** GET  
**URL:** `{{base_url}}/health/ready`

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Service is ready", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data.status).to.eql("ready");
    pm.expect(jsonData.data.database).to.eql("healthy");
});
```

---

#### 2.3 Liveness Probe

**Method:** GET  
**URL:** `{{base_url}}/health/live`

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Service is alive", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data.status).to.eql("alive");
});
```

---

### Folder 3: Password Management

#### 3.1 Forgot Password

**Method:** POST  
**URL:** `{{base_url}}/v1/auth/forgot-password`  
**Body (JSON):**
```json
{
  "email": "{{user_email}}"
}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Request processed", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.success).to.be.true;
});
```

---

#### 3.2 Reset Password

**Method:** POST  
**URL:** `{{base_url}}/v1/auth/reset-password`  
**Body (JSON):**
```json
{
  "token": "RESET_TOKEN_HERE",
  "new_password": "newpassword123"
}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Password reset successful", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.success).to.be.true;
});
```

---

#### 3.3 Change Password

**Method:** POST  
**URL:** `{{base_url}}/v1/auth/change-password`  
**Headers:**
```
Authorization: Bearer {{token}}
```
**Body (JSON):**
```json
{
  "old_password": "{{user_password}}",
  "new_password": "newpassword456"
}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Password changed successfully", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.success).to.be.true;
});

// Update password in environment
pm.environment.set("user_password", "newpassword456");
```

---

### Folder 4: User Management

#### 4.1 Get Profile

**Method:** GET  
**URL:** `{{base_url}}/v1/api/profile`  
**Headers:**
```
Authorization: Bearer {{token}}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Profile contains user data", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('id');
    pm.expect(jsonData.data).to.have.property('email');
    pm.expect(jsonData.data).to.have.property('role');
});
```

---

#### 4.2 List Users (Admin)

**Method:** GET  
**URL:** `{{base_url}}/v1/admin/users?limit=10&offset=0`  
**Headers:**
```
Authorization: Bearer {{admin_token}}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Response contains users array", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('users');
    pm.expect(jsonData.data.users).to.be.an('array');
});

pm.test("Response contains pagination info", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('total');
    pm.expect(jsonData.data).to.have.property('limit');
    pm.expect(jsonData.data).to.have.property('offset');
});
```

---

#### 4.3 Get User by ID (Admin)

**Method:** GET  
**URL:** `{{base_url}}/v1/admin/users/{{user_id}}`  
**Headers:**
```
Authorization: Bearer {{admin_token}}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("User data returned", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('id');
    pm.expect(jsonData.data.id).to.eql(parseInt(pm.environment.get("user_id")));
});
```

---

#### 4.4 Update User (Admin)

**Method:** PUT  
**URL:** `{{base_url}}/v1/admin/users/{{user_id}}`  
**Headers:**
```
Authorization: Bearer {{admin_token}}
```
**Body (JSON):**
```json
{
  "email": "updated@example.com",
  "is_verified": true
}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("User updated successfully", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.success).to.be.true;
});
```

---

#### 4.5 Delete User (Admin)

**Method:** DELETE  
**URL:** `{{base_url}}/v1/admin/users/{{user_id}}`  
**Headers:**
```
Authorization: Bearer {{admin_token}}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("User deleted successfully", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.success).to.be.true;
});
```

---

### Folder 5: Admin Operations

#### 5.1 Admin Dashboard

**Method:** GET  
**URL:** `{{base_url}}/v1/admin/dashboard`  
**Headers:**
```
Authorization: Bearer {{admin_token}}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("Dashboard data returned", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.data).to.have.property('total_users');
});
```

---

#### 5.2 Promote User to Admin

**Method:** POST  
**URL:** `{{base_url}}/v1/super-admin/promote`  
**Headers:**
```
Authorization: Bearer {{admin_token}}
```
**Body (JSON):**
```json
{
  "user_id": {{user_id}}
}
```

**Tests Script:**
```javascript
pm.test("Status code is 200", function () {
    pm.response.to.have.status(200);
});

pm.test("User promoted successfully", function () {
    var jsonData = pm.response.json();
    pm.expect(jsonData.success).to.be.true;
    pm.expect(jsonData.data.new_role).to.eql("admin");
});
```

---

## Testing Workflows

### Complete User Flow

Run these requests in order:

1. Register User
2. Login User
3. Get Profile
4. Validate Token
5. Change Password
6. Login with new password
7. Logout

### Complete Admin Flow

Run these requests in order:

1. Register Super Admin
2. Login Admin
3. Register Regular User
4. List Users
5. Get User by ID
6. Admin Dashboard
7. Promote User to Admin

---

## Automated Tests

### Collection Runner

1. Click "Runner" in Postman
2. Select "Go Auth Service" collection
3. Select environment
4. Click "Run Go Auth Service"

### Pre-request Scripts

Add to collection level for automatic token refresh:

```javascript
// Check if token is about to expire
const token = pm.environment.get("token");
if (token) {
    // Decode JWT (simplified)
    const base64Url = token.split('.')[1];
    const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
    const jsonPayload = decodeURIComponent(atob(base64).split('').map(function(c) {
        return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
    }).join(''));
    
    const payload = JSON.parse(jsonPayload);
    const exp = payload.exp * 1000; // Convert to milliseconds
    const now = Date.now();
    
    // Refresh if expiring in less than 1 hour
    if (exp - now < 3600000) {
        console.log("Token expiring soon, refreshing...");
        // Trigger refresh endpoint
    }
}
```

---

## Tips and Tricks

### 1. Global Headers

Set at collection level:
```
Content-Type: application/json
Accept: application/json
```

### 2. Dynamic Variables

Use Postman's dynamic variables:
```json
{
  "email": "{{$randomEmail}}",
  "password": "{{$randomPassword}}"
}
```

### 3. Console Logging

Add to tests for debugging:
```javascript
console.log("Response:", pm.response.json());
console.log("Token:", pm.environment.get("token"));
```

### 4. Conditional Tests

```javascript
if (pm.response.code === 200) {
    pm.test("Success response", function () {
        var jsonData = pm.response.json();
        pm.expect(jsonData.success).to.be.true;
    });
} else {
    pm.test("Error response", function () {
        var jsonData = pm.response.json();
        pm.expect(jsonData).to.have.property('error');
    });
}
```

### 5. Response Time Tests

```javascript
pm.test("Response time is less than 500ms", function () {
    pm.expect(pm.response.responseTime).to.be.below(500);
});
```

### 6. Header Tests

```javascript
pm.test("Has Request ID header", function () {
    pm.response.to.have.header("X-Request-ID");
});

pm.test("Has Rate Limit headers", function () {
    pm.response.to.have.header("X-RateLimit-Limit");
    pm.response.to.have.header("X-RateLimit-Remaining");
});
```

---

## Export/Import Collection

### Export

1. Right-click collection
2. Select "Export"
3. Choose "Collection v2.1"
4. Save as `Go_Auth_Service.postman_collection.json`

### Import

1. Click "Import" in Postman
2. Select the JSON file
3. Collection will be imported with all requests

---

## Troubleshooting

### Token Not Saving

Check Tests script:
```javascript
pm.environment.set("token", jsonData.data.token);
```

### 401 Unauthorized

- Check token is set in environment
- Verify Authorization header format: `Bearer {{token}}`
- Check token hasn't expired

### Rate Limit Errors

- Wait for rate limit reset time
- Check `X-RateLimit-Reset` header
- Reduce request frequency

---

## Best Practices

1. **Use environments** for different stages (dev, staging, prod)
2. **Save tokens automatically** in test scripts
3. **Add descriptive test names** for clarity
4. **Use folders** to organize requests
5. **Document requests** with descriptions
6. **Share collections** with team via export
7. **Use variables** instead of hardcoded values
8. **Test error cases** not just success cases

---

## Support

For Postman-related issues:
- Check environment variables are set
- Verify base URL is correct
- Review test script errors in console
- Check request/response in Postman Console
