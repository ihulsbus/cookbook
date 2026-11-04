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
	err := c.RabbitMQHandler.StartConsuming(c.RabbitMQClient.Connection, "images", "cookbook")
	if err != nil {
		c.Logger.Fatalf("Startup of RabbitMQ Consumer encountered fatal error: %v", err.Error())
		return
	}
	defer c.RabbitMQHandler.StopConsuming()

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
			// Images are owned by whatever recipe/ingredient they're
			// attached to - any authenticated user, same as the recipe
			// CRUD they ride along with (per-recipe ownership checks are
			// out of scope here, see the ownership epic).
			readImage := image.Group("")
			readImage.Use(c.OIDCModule.RequireAuth())
			{
				readImage.GET("", c.HttpHandler.FindAll)
				readImage.GET(":id", c.HttpHandler.Find)
				readImage.GET("/search", c.HttpHandler.SearchByRecipe)
				// readImage.GET(":entityType/:entityID", c.HttpHandler.Find)
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
