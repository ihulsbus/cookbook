package instructionservice

import (
	"context"
	c "image-service/internal/config"
	m "image-service/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
)

var (
	log = c.Logger
)

func ImageService(ctx context.Context) {
	c.RabbitMQHandler.StartConsuming(c.RabbitMQClient, "images", "cookbook")
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
			readRecipe := image.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				readRecipe.GET("", c.ImageHandler.FindAll)
				readRecipe.GET(":id", c.ImageHandler.Find)
			}

			createRecipe := image.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				createRecipe.POST("", c.ImageHandler.Create)
			}

			updateRecipe := image.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{

				updateRecipe.PUT(":id", c.ImageHandler.Update)
			}

			adminRecipe := image.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				adminRecipe.DELETE(":id", c.ImageHandler.Delete)
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
		srv.Shutdown(ctx)
	}()

	log.Info("instruction service available on port 8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
}
