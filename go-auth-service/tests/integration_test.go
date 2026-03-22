package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCompleteUserFlow tests the complete user workflow
func TestCompleteUserFlow(t *testing.T) {
	// Step 1: Register
	registerPayload := map[string]interface{}{
		"email":    "flowtest@example.com",
		"password": "password123",
	}
	body, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
	var registerResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &registerResponse)
	assert.True(t, registerResponse["success"].(bool))

	// Step 2: Login
	loginPayload := map[string]interface{}{
		"email":    "flowtest@example.com",
		"password": "password123",
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.True(t, loginResponse["success"].(bool))
	token := loginResponse["data"].(map[string]interface{})["token"].(string)
	assert.NotEmpty(t, token)

	// Step 3: Get Profile
	req, _ = http.NewRequest("GET", "/v1/api/profile", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var profileResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &profileResponse)
	assert.True(t, profileResponse["success"].(bool))
	profileData := profileResponse["data"].(map[string]interface{})
	assert.Equal(t, "flowtest@example.com", profileData["email"])

	// Step 4: Validate Token
	validatePayload := map[string]interface{}{
		"token": token,
	}
	body, _ = json.Marshal(validatePayload)
	req, _ = http.NewRequest("POST", "/v1/auth/validate", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var validateResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &validateResponse)
	assert.True(t, validateResponse["success"].(bool))
	validateData := validateResponse["data"].(map[string]interface{})
	assert.True(t, validateData["valid"].(bool))

	// Step 5: Change Password
	changePassPayload := map[string]interface{}{
		"old_password": "password123",
		"new_password": "newpassword456",
	}
	body, _ = json.Marshal(changePassPayload)
	req, _ = http.NewRequest("POST", "/v1/auth/change-password", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var changePassResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &changePassResponse)
	assert.True(t, changePassResponse["success"].(bool))

	// Step 6: Login with new password
	loginPayload = map[string]interface{}{
		"email":    "flowtest@example.com",
		"password": "newpassword456",
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var newLoginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &newLoginResponse)
	assert.True(t, newLoginResponse["success"].(bool))
	newToken := newLoginResponse["data"].(map[string]interface{})["token"].(string)
	assert.NotEmpty(t, newToken)

	// Step 7: Logout
	req, _ = http.NewRequest("POST", "/v1/auth/logout", nil)
	req.Header.Set("Authorization", "Bearer "+newToken)
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var logoutResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &logoutResponse)
	assert.True(t, logoutResponse["success"].(bool))
}

// TestCompleteAdminFlow tests the complete admin workflow
func TestCompleteAdminFlow(t *testing.T) {
	// Step 1: Register Super Admin
	registerPayload := map[string]interface{}{
		"email":            "superadmin@example.com",
		"password":         "adminpass123",
		"super_admin_code": "TEST_SUPER_ADMIN_CODE",
	}
	body, _ := json.Marshal(registerPayload)
	req, _ := http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Step 2: Login as Super Admin
	loginPayload := map[string]interface{}{
		"email":    "superadmin@example.com",
		"password": "adminpass123",
	}
	body, _ = json.Marshal(loginPayload)
	req, _ = http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var loginResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &loginResponse)
	adminToken := loginResponse["data"].(map[string]interface{})["token"].(string)

	// Step 3: Register Regular User
	registerPayload = map[string]interface{}{
		"email":    "regularuser@example.com",
		"password": "userpass123",
	}
	body, _ = json.Marshal(registerPayload)
	req, _ = http.NewRequest("POST", "/v1/auth/register", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)

	// Step 4: Admin Lists Users
	req, _ = http.NewRequest("GET", "/v1/admin/users", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var listResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &listResponse)
	assert.True(t, listResponse["success"].(bool))

	// Step 5: Admin Accesses Dashboard
	req, _ = http.NewRequest("GET", "/v1/admin/dashboard", nil)
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var dashboardResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &dashboardResponse)
	assert.True(t, dashboardResponse["success"].(bool))

	// Step 6: Promote User to Admin
	promotePayload := map[string]interface{}{
		"user_id": 2,
	}
	body, _ = json.Marshal(promotePayload)
	req, _ = http.NewRequest("POST", "/v1/super-admin/promote", bytes.NewBuffer(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+adminToken)
	w = httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var promoteResponse map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &promoteResponse)
	assert.True(t, promoteResponse["success"].(bool))
}

// TestRateLimiting tests rate limiting functionality
func TestRateLimiting(t *testing.T) {
	// This test would need to make multiple rapid requests
	// For now, we'll just verify the endpoint exists
	loginPayload := map[string]interface{}{
		"email":    "ratelimit@example.com",
		"password": "wrongpassword",
	}

	// Make multiple requests
	for i := 0; i < 3; i++ {
		body, _ := json.Marshal(loginPayload)
		req, _ := http.NewRequest("POST", "/v1/auth/login", bytes.NewBuffer(body))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()
		TestRouter.ServeHTTP(w, req)

		// Check for rate limit headers
		assert.NotEmpty(t, w.Header().Get("X-RateLimit-Limit"))
		assert.NotEmpty(t, w.Header().Get("X-RateLimit-Remaining"))
	}
}

// TestCORSHeaders tests CORS headers are present
func TestCORSHeaders(t *testing.T) {
	req, _ := http.NewRequest("OPTIONS", "/v1/auth/login", nil)
	req.Header.Set("Origin", "http://localhost:3000")
	req.Header.Set("Access-Control-Request-Method", "POST")

	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	// Check CORS headers
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Origin"))
	assert.NotEmpty(t, w.Header().Get("Access-Control-Allow-Methods"))
}

// TestRequestIDHeader tests that request ID is added to responses
func TestRequestIDHeader(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.NotEmpty(t, w.Header().Get("X-Request-ID"))
}
