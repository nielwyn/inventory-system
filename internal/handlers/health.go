package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nielwyn/inventory-system/internal/database"
	"github.com/nielwyn/inventory-system/pkg/response"
)

// HealthHandler handles health check endpoints
type HealthHandler struct {
	db *database.Database
}

// NewHealthHandler creates a new health handler
func NewHealthHandler(db *database.Database) *HealthHandler {
	return &HealthHandler{db: db}
}

// Health handles basic health check
func (h *HealthHandler) Health(c *gin.Context) {
	response.Success(c, http.StatusOK, "Service is healthy", gin.H{
		"status": "ok",
	})
}

// Ready handles readiness check with database ping
func (h *HealthHandler) Ready(c *gin.Context) {
	// Check database connection
	if err := h.db.Health(); err != nil {
		response.Error(c, http.StatusServiceUnavailable, "Database is not ready")
		return
	}

	response.Success(c, http.StatusOK, "Service is ready", gin.H{
		"status":   "ok",
		"database": "connected",
	})
}
