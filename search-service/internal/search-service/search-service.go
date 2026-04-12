package searchservice

import (
	"context"
	"net/http"
	c "search-service/internal/config"
	m "search-service/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	log = c.Logger
)

func IngredientService(ctx context.Context) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	// Logging
	router.Use(m.Logger(log))

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

	v2 := router.Group("/api/v2")
	{
		search := v2.Group("/search")
		{
			doSearch := search.Group("")
			doSearch.Use(c.KeycloakModule.Middleware("administrator"))
			{
				doSearch.GET("", c.SearchHandlers.Search)
			}
		}

	}

	// Server startup
	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + c.Configuration.Global.ListenPort,
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
	}()

	log.Info("search service available on port 8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
}
