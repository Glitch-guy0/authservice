package health

import (
	"context"
	"net/http"
	"time"

	"github.com/Glitch-guy0/authService/modules/core"
	"github.com/Glitch-guy0/authService/modules/logger"
	"github.com/gin-gonic/gin"
)

type HealthHandler struct {
	healthService *HealthService
	logger        logger.Logger
}

func NewHealthHandler(appCtx *core.AppContext) *HealthHandler {
	config := DefaultHealthCheckConfig()
	healthService := NewHealthService(appCtx, config)

	handler := &HealthHandler{
		healthService: healthService,
		logger:        appCtx.GetLogger(),
	}

	handler.registerDefaultCheckers()
	return handler
}

func (hh *HealthHandler) registerDefaultCheckers() {
	hh.healthService.RegisterChecker(NewBasicHealthChecker("server", func(ctx context.Context) Check {
		return Check{
			Name:      "server",
			Status:    StatusHealthy,
			Message:   "Server is running",
			Timestamp: time.Now(),
		}
	}))

	hh.healthService.RegisterChecker(NewBasicHealthChecker("database", func(ctx context.Context) Check {
		return Check{
			Name:      "database",
			Status:    StatusHealthy,
			Message:   "Database connection not configured",
			Timestamp: time.Now(),
		}
	}))

	hh.healthService.RegisterChecker(NewBasicHealthChecker("logger", func(ctx context.Context) Check {
		return Check{
			Name:      "logger",
			Status:    StatusHealthy,
			Message:   "Logger is operational",
			Timestamp: time.Now(),
		}
	}))
}

func (hh *HealthHandler) HealthCheck(c *gin.Context) {
	ctx := c.Request.Context()
	health := hh.healthService.GetHealth(ctx)

	statusCode := http.StatusOK
	switch health.Status {
	case StatusUnhealthy:
		statusCode = http.StatusServiceUnavailable
	case StatusDegraded:
		statusCode = http.StatusOK // Still return 200 for degraded
	}

	c.JSON(statusCode, health)
}

func (hh *HealthHandler) LivenessProbe(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":    "ok",
		"timestamp": time.Now(),
	})
}

func (hh *HealthHandler) ReadinessProbe(c *gin.Context) {
	ctx := c.Request.Context()
	health := hh.healthService.GetHealth(ctx)

	if health.Status == StatusHealthy {
		c.JSON(http.StatusOK, gin.H{
			"status":    "ready",
			"timestamp": time.Now(),
		})
	} else {
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"status":    "not ready",
			"timestamp": time.Now(),
		})
	}
}

func (hh *HealthHandler) GetHealthService() *HealthService {
	return hh.healthService
}
