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
// start - this service only publishes events, nothing consumes them
// yet.
func UserService(ctx context.Context) {
	httpServer(ctx)
}

func httpServer(ctx context.Context) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	router.Use(m.Logger(c.Logger))
	router.Use(gin.Recovery())
	router.Use(cors.New(c.Cors))

	// Health check endpoints - no auth required.
	health := router.Group("/health")
	{
		health.GET("/live", c.HealthHandler.Liveness)
		health.GET("/ready", c.HealthHandler.Readiness)
		health.GET("/startup", c.HealthHandler.Startup)
	}

	v2 := router.Group("/api/v2/user")
	{
		// User lookup: every route only requires a valid token - no
		// ownership to enforce, just an exact-match search and id
		// resolve.
		users := v2.Group("/users")
		users.Use(c.OIDCModule.RequireAuth())
		{
			users.GET("", c.HttpHandler.SearchUsers)
			users.GET(":id", c.HttpHandler.GetUser)
		}

		me := v2.Group("/me")
		me.Use(c.OIDCModule.RequireAuth())
		{
			me.GET("/households", c.HttpHandler.ListMyHouseholds)
		}

		// Self-service households: owner-only operations are enforced
		// in the service layer, not here.
		households := v2.Group("/households")
		households.Use(c.OIDCModule.RequireAuth())
		{
			households.POST("", c.HttpHandler.CreateHousehold)
			households.POST(":id/members", c.HttpHandler.AddMember)
			households.DELETE(":id/members/:userId", c.HttpHandler.RemoveMember)
			households.POST(":id/leave", c.HttpHandler.LeaveHousehold)
			households.DELETE(":id", c.HttpHandler.DeleteHousehold)
		}
	}

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
