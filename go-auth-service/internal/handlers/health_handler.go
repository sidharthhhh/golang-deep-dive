package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sidharthhhh/go-auth-service/internal/pkg/response"
	"go.uber.org/zap"
)

// HealthHandler handles health check requests
type HealthHandler struct {
	db        *sql.DB
	logger    *zap.Logger
	startTime time.Time
	version   string
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *sql.DB, logger *zap.Logger, version string) *HealthHandler {
	return &HealthHandler{
		db:        db,
		logger:    logger,
		startTime: time.Now(),
		version:   version,
	}
}

// HealthResponse represents the health check response
type HealthResponse struct {
	Status    string `json:"status"`
	Version   string `json:"version"`
	Uptime    string `json:"uptime"`
	Timestamp string `json:"timestamp"`
}

// ReadinessResponse represents the readiness probe response
type ReadinessResponse struct {
	Status   string            `json:"status"`
	Checks   map[string]string `json:"checks"`
	Timestamp string           `json:"timestamp"`
}

// Health returns basic health status
// @Summary Health check
// @Description Returns basic health status of the service
// @Tags Health
// @Produce json
// @Success 200 {object} response.Response{data=HealthResponse}
// @Router /health [get]
func (h *HealthHandler) Health(c *gin.Context) {
	uptime := time.Since(h.startTime)

	data := HealthResponse{
		Status:    "healthy",
		Version:   h.version,
		Uptime:    uptime.String(),
		Timestamp: time.Now().Format(time.RFC3339),
	}

	response.Success(c, http.StatusOK, "Service is healthy", data)
}

// Ready checks if the service is ready to accept requests
// @Summary Readiness probe
// @Description Checks if service is ready (database connection, etc.)
// @Tags Health
// @Produce json
// @Success 200 {object} response.Response{data=ReadinessResponse}
// @Failure 503 {object} response.Response
// @Router /health/ready [get]
func (h *HealthHandler) Ready(c *gin.Context) {
	checks := make(map[string]string)

	// Check database connection
	if err := h.db.Ping(); err != nil {
		h.logger.Error("database ping failed", zap.Error(err))
		checks["database"] = "unavailable"

		data := ReadinessResponse{
			Status:    "not_ready",
			Checks:    checks,
			Timestamp: time.Now().Format(time.RFC3339),
		}

		c.JSON(http.StatusServiceUnavailable, response.Response{
			Success:   false,
			Message:   "Service is not ready",
			Data:      data,
			Timestamp: time.Now().Format(time.RFC3339),
		})
		return
	}

	checks["database"] = "connected"

	data := ReadinessResponse{
		Status:    "ready",
		Checks:    checks,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	response.Success(c, http.StatusOK, "Service is ready", data)
}

// Live checks if the service is alive
// @Summary Liveness probe
// @Description Checks if service is alive (for Kubernetes)
// @Tags Health
// @Produce json
// @Success 200 {object} response.Response
// @Router /health/live [get]
func (h *HealthHandler) Live(c *gin.Context) {
	response.Success(c, http.StatusOK, "Service is alive", gin.H{
		"status":    "alive",
		"timestamp": time.Now().Format(time.RFC3339),
	})
}
