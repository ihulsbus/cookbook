package instructionservice

import (
	"context"
	c "metadata-service/internal/config"
	m "metadata-service/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
)

var (
	log = c.Logger
)

func MetadataService(ctx context.Context) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	// Logging
	router.Use(m.Logger(log))

	// Panic recovery
	router.Use(gin.Recovery())

	// Cors handler
	router.Use(cors.New(c.Cors))

	// API versioning setup
	v1 := router.Group("/api/v1")
	{
		tag := v1.Group("/tag")
		{
			readTag := tag.Group("")
			readTag.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				readTag.GET("", c.TagHandlers.GetAll)
				readTag.GET(":id", c.TagHandlers.Get)
			}

			createTag := tag.Group("")
			createTag.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				createTag.POST("", c.TagHandlers.Create)
			}

			updateTag := tag.Group("")
			updateTag.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				updateTag.PUT(":id", c.TagHandlers.Update)
			}

			deleteTag := tag.Group("")
			deleteTag.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				deleteTag.DELETE(":id", c.TagHandlers.Delete)
			}
		}

		category := v1.Group("/category")
		{
			readCategory := category.Group("")
			readCategory.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				readCategory.GET("", c.CategoryHandlers.GetAll)
				readCategory.GET(":id", c.CategoryHandlers.Get)
			}

			createCategory := category.Group("")
			createCategory.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				createCategory.POST("", c.CategoryHandlers.Create)
			}

			updateCategory := category.Group("")
			updateCategory.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				updateCategory.PUT(":id", c.CategoryHandlers.Update)
			}

			deleteCategory := category.Group("")
			deleteCategory.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				deleteCategory.DELETE(":id", c.CategoryHandlers.Delete)
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

	log.Info("metadata service available on port 8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
}
