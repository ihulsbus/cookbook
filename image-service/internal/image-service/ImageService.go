package imageservice

import (
	"context"
	c "image-service/internal/config"
	m "image-service/internal/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func ImageService(ctx context.Context) {
	// Keeps the recipe_access projection current; no synchronous call to
	// recipe-service happens at startup (KIT-79).
	consumer, err := c.Projector.StartConsumer(c.RabbitMQClient.Connection, "image-service-recipe-access", "cookbook")
	if err != nil {
		c.Logger.Fatalf("Startup of RabbitMQ Consumer encountered fatal error: %v", err.Error())
		return
	}
	defer consumer.Close()

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

	v2 := router.Group("/api/v2")
	{
		image := v2.Group("/images")
		{
			// Images attached to a recipe (EntityType "recipe") are
			// authorized against that recipe; images attached to
			// anything else (e.g. "ingredient") have no owner, so any
			// authenticated caller may act on them (KIT-29). The
			// resolve-then-authorize logic lives in HttpHandler, since
			// Find/Update/Delete take an IMAGE id and need a DB lookup
			// first to learn the owning recipe - RequireAuth is the only
			// thing enforceable declaratively here.
			//
			// FindAll (list every image) has no cheap per-row filter
			// against the projection, so it's admin-only instead.
			adminReadImage := image.Group("")
			adminReadImage.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				adminReadImage.GET("", c.HttpHandler.FindAll)
			}

			readImage := image.Group("")
			readImage.Use(c.OIDCModule.RequireAuth())
			{
				readImage.GET(":id", c.HttpHandler.Find)
				readImage.GET("/search", c.HttpHandler.SearchByRecipe)
			}

			createImage := image.Group("")
			createImage.Use(c.OIDCModule.RequireAuth())
			{
				createImage.POST(":entityType/:entityID", c.HttpHandler.Create)
			}

			updateImage := image.Group("")
			updateImage.Use(c.OIDCModule.RequireAuth())
			{

				updateImage.PUT(":id", c.HttpHandler.Update)
			}

			adminImage := image.Group("")
			adminImage.Use(c.OIDCModule.RequireAuth())
			{
				adminImage.DELETE(":id", c.HttpHandler.Delete)
			}
		}

	}

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

	c.Logger.Infof("image service available on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		c.Logger.Error(err)
	}
}
