package userservice

import (
	"context"
	"net/http"
	"strconv"
	"time"

	c "user-service/internal/config"
	m "user-service/internal/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// UserService starts the HTTP server. There's no RabbitMQ consumer to
// start (user-service only publishes household.member.added/removed,
// KIT-74 - nothing in this service consumes events yet) and no business
// routes either: user and household lookup (KIT-33) isn't built here.
func UserService(ctx context.Context) {
	httpServer(ctx)
}

func httpServer(ctx context.Context) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	// Logging
	router.Use(m.Logger(c.Logger))

	// Panic recovery
	router.Use(gin.Recovery())

	// Cors handler
	router.Use(cors.New(c.Cors))

	// Health check endpoints (no auth required)
	health := router.Group("/health")
	{
		health.GET("/live", c.HealthHandler.Liveness)
		health.GET("/ready", c.HealthHandler.Readiness)
		health.GET("/startup", c.HealthHandler.Startup)
	}

	// /api/v2/user routes (user/household lookup, KIT-33) aren't built
	// yet - this service only exists so far to own the identity DB and
	// publish household events (KIT-74). The Traefik route for
	// /api/v2/user is already wired (docker-compose.yml, k8s ingress)
	// ahead of there being anything behind it to route to.

	// Server startup
	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + strconv.Itoa(c.Configuration.Global.ListenPort),
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		<-ctx.Done()
		c.Logger.Info("Stopping webserver")
		srv.Shutdown(ctx)
	}()

	c.Logger.Infof("user service available on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		c.Logger.Error(err)
	}
}
