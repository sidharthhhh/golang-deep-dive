package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestHealthCheck tests the basic health check endpoint
func TestHealthCheck(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health", nil)
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response["success"].(bool))
	assert.Equal(t, "Service is healthy", response["message"])
}

// TestHealthReady tests the readiness probe endpoint
func TestHealthReady(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health/ready", nil)
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "ready", data["status"])
	assert.Equal(t, "healthy", data["database"])
}

// TestHealthLive tests the liveness probe endpoint
func TestHealthLive(t *testing.T) {
	req, _ := http.NewRequest("GET", "/health/live", nil)
	w := httptest.NewRecorder()
	TestRouter.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.True(t, response["success"].(bool))
	data := response["data"].(map[string]interface{})
	assert.Equal(t, "alive", data["status"])
}
