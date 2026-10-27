package health

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type Handlers struct {
	db     *gorm.DB
	logger LoggerInterface
}

type LoggerInterface interface {
	Debugf(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
}

// Response represents the health check response
type Response struct {
	Status      string            `json:"status"`
	Version     string            `json:"version,omitempty"`
	Timestamp   string            `json:"timestamp"`
	Checks      map[string]string `json:"checks,omitempty"`
	Description string            `json:"description,omitempty"`
}

func NewHealthHandlers(db *gorm.DB, logger LoggerInterface) *Handlers {
	var database *gorm.DB
	if db != nil {
		database = db
	} else {
		database = nil
	}
	return &Handlers{
		db:     database,
		logger: logger,
	}
}

// Liveness checks if the service is running
// This should be lightweight and only check if the app is alive
// Kubernetes will restart the pod if this fails
func (h *Handlers) Liveness(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, Response{
		Status:      "UP",
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Description: "Service is alive",
	})
}

// Readiness checks if the service is ready to accept traffic
// This checks dependencies like a database, external services, etc.
// Kubernetes will remove the pod from load balancer if this fails
func (h *Handlers) Readiness(ctx *gin.Context) {
	checks := make(map[string]string)
	overallStatus := "UP"

	// Check database connection
	if h.db != nil {
		if err := h.checkDatabase(); err != nil {
			checks["database"] = "DOWN: " + err.Error()
			overallStatus = "DOWN"
			h.logger.Warnf("Readiness check failed: database is down - %v", err)
		} else {
			checks["database"] = "UP"
		}
	}

	// Add more dependency checks here as needed
	// Example: RabbitMQ, Redis, external APIs, etc.

	statusCode := http.StatusOK
	if overallStatus == "DOWN" {
		statusCode = http.StatusServiceUnavailable
	}

	ctx.JSON(statusCode, Response{
		Status:      overallStatus,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Checks:      checks,
		Description: "Service readiness check",
	})
}

// Startup checks if the service has finished starting up
// This is useful for slow-starting services
// Kubernetes will wait longer before checking liveness/readiness
func (h *Handlers) Startup(ctx *gin.Context) {
	checks := make(map[string]string)
	overallStatus := "UP"

	// Check database connection
	if err := h.checkDatabase(); err != nil {
		checks["database"] = "DOWN: " + err.Error()
		overallStatus = "DOWN"
	} else {
		checks["database"] = "UP"
	}

	// Add startup-specific checks here
	// Example: check if migrations ran, cache is warmed up, etc.

	statusCode := http.StatusOK
	if overallStatus == "DOWN" {
		statusCode = http.StatusServiceUnavailable
	}

	ctx.JSON(statusCode, Response{
		Status:      overallStatus,
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Checks:      checks,
		Description: "Service startup check",
	})
}

// checkDatabase verifies database connectivity
func (h *Handlers) checkDatabase() error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sqlDB, err := h.db.DB()
	if err != nil {
		return err
	}

	return sqlDB.PingContext(ctx)
}
