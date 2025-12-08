package health

import (
	"github.com/Glitch-guy0/authService/modules/core"
	"github.com/gin-gonic/gin"
)

// Controller handles health-related HTTP requests
type Controller struct {
	appCtx *core.AppContext
}

// NewController creates a new health controller
func NewController(appCtx *core.AppContext) *Controller {
	return &Controller{
		appCtx: appCtx,
	}
}

// RegisterRoutes registers all health routes
func (c *Controller) RegisterRoutes(router *gin.RouterGroup) {
	// Health check endpoints
	router.GET("", c.HealthCheck)
	router.GET("/check", c.HealthCheck)
	router.GET("/live", c.LivenessProbe)
	router.GET("/ready", c.ReadinessProbe)
}

// HealthCheck handles the health check endpoint
func (c *Controller) HealthCheck(ctx *gin.Context) {
	h := NewHealthHandler(c.appCtx)
	h.HealthCheck(ctx)
}

// LivenessProbe handles the liveness probe endpoint
func (c *Controller) LivenessProbe(ctx *gin.Context) {
	h := NewHealthHandler(c.appCtx)
	h.LivenessProbe(ctx)
}

// ReadinessProbe handles the readiness probe endpoint
func (c *Controller) ReadinessProbe(ctx *gin.Context) {
	h := NewHealthHandler(c.appCtx)
	h.ReadinessProbe(ctx)
}
