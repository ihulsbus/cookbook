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
	// Keeps the recipe_access projection current without a synchronous call to recipe-service.
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

	router.Use(m.Logger(c.Logger))
	router.Use(gin.Recovery())
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
			// Resolve-then-authorize (DB lookup of the owning recipe) lives in
			// HttpHandler since Find/Update/Delete only have an image id to start from.
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
