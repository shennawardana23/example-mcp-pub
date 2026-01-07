package controller

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/shennawardana23/example-mcp-pub/internal/app/model"
	"github.com/shennawardana23/example-mcp-pub/internal/database"
)

// HealthController handles health check requests
type HealthController struct {
	db    *database.Database
	redis *database.RedisClient
}

// NewHealthController creates a new health controller
func NewHealthController(db *database.Database, redis *database.RedisClient) *HealthController {
	return &HealthController{
		db:    db,
		redis: redis,
	}
}

// Health handles health check requests
// @Summary Health check
// @Description Check the health of the API and its dependencies
// @Tags health
// @Produce json
// @Success 200 {object} model.HealthResponse
// @Failure 503 {object} model.HealthResponse
// @Router /health [get]
func (ctrl *HealthController) Health(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Second)
	defer cancel()

	services := make(map[string]string)

	// Check database
	if err := ctrl.db.Health(ctx); err != nil {
		services["database"] = "unhealthy: " + err.Error()
	} else {
		services["database"] = "healthy"
	}

	// Check Redis
	if ctrl.redis != nil {
		if err := ctrl.redis.Health(ctx); err != nil {
			services["redis"] = "unhealthy: " + err.Error()
		} else {
			services["redis"] = "healthy"
		}
	}

	// Determine overall status
	status := "ok"
	statusCode := http.StatusOK
	for _, serviceStatus := range services {
		if serviceStatus != "healthy" {
			status = "degraded"
			statusCode = http.StatusServiceUnavailable
			break
		}
	}

	response := model.HealthResponse{
		Status:    status,
		Timestamp: time.Now(),
		Services:  services,
	}

	c.JSON(statusCode, response)
}
