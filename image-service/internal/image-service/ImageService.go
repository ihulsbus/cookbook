package instructionservice

import (
	"context"
	c "image-service/internal/config"
	m "image-service/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	log = c.Logger
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
	router.Use(m.Logger(log))

	// Panic recovery
	router.Use(gin.Recovery())

	// Cors handler
	router.Use(cors.New(c.Cors))

	v2 := router.Group("/api/v2")
	{
		image := v2.Group("/images")
		{
			readImage := image.Group("")
			readImage.Use(c.KeycloakModule.Middleware("administrator"))
			{
				readImage.GET("", c.HttpHandler.FindAll)
				readImage.GET(":id", c.HttpHandler.Find)
				// readImage.GET(":entityType/:entityID", c.HttpHandler.Find)
			}

			createImage := image.Group("")
			createImage.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createImage.POST(":entityType/:entityID", c.HttpHandler.Create)
			}

			updateImage := image.Group("")
			updateImage.Use(c.KeycloakModule.Middleware("administrator"))
			{

				updateImage.PUT(":id", c.HttpHandler.Update)
			}

			adminImage := image.Group("")
			adminImage.Use(c.KeycloakModule.Middleware("administrator"))
			{
				adminImage.DELETE(":id", c.HttpHandler.Delete)
			}
		}

	}

	// Server startup
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8080",
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		<-ctx.Done()
		log.Info("Stopping webserver")
		srv.Shutdown(ctx)
	}()

	log.Info("instruction service available on port 8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
}
