package metadataservice

import (
	"context"
	"errors"
	c "metadata-service/internal/config"
	m "metadata-service/internal/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	log = c.Logger
)

func MetadataService(ctx context.Context) {
	err := c.RabbitMQHandler.StartConsuming(c.RabbitMQClient.Connection, "metadata", "cookbook")
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

	// Health check endpoints (no auth required)
	health := router.Group("/health")
	{
		health.GET("/live", c.HealthHandler.Liveness)
		health.GET("/ready", c.HealthHandler.Readiness)
		health.GET("/startup", c.HealthHandler.Startup)
	}

	// API versioning setup
	v2 := router.Group("/api/v2")
	metadata := v2.Group("/metadata")
	{

		recipe := metadata.Group("/recipe")
		{
			readMetadata := recipe.Group("")
			readMetadata.Use(c.KeycloakModule.Middleware("administrator"))
			{
				readMetadata.GET("", c.MetadataHandlers.GetAll)
				readMetadata.GET(":id", c.MetadataHandlers.Get)
			}
			createMetadata := recipe.Group("")
			createMetadata.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createMetadata.POST(":id", c.MetadataHandlers.Create)
			}

			updateMetadata := recipe.Group("")
			updateMetadata.Use(c.KeycloakModule.Middleware("administrator"))
			{
				updateMetadata.PUT(":id", c.MetadataHandlers.Update)
			}

			deleteMetadata := recipe.Group("")
			deleteMetadata.Use(c.KeycloakModule.Middleware("administrator"))
			{
				deleteMetadata.DELETE(":id", c.MetadataHandlers.Delete)
			}
		}

		// Tag routes
		tag := metadata.Group("/tag")
		{
			readTag := tag.Group("")
			readTag.Use(c.KeycloakModule.Middleware("administrator"))
			{
				readTag.GET("", c.TagHandlers.GetAll)
				readTag.GET(":id", c.TagHandlers.Get)
			}

			createTag := tag.Group("")
			createTag.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createTag.POST("", c.TagHandlers.Create)
			}

			updateTag := tag.Group("")
			updateTag.Use(c.KeycloakModule.Middleware("administrator"))
			{
				updateTag.PUT(":id", c.TagHandlers.Update)
			}

			deleteTag := tag.Group("")
			deleteTag.Use(c.KeycloakModule.Middleware("administrator"))
			{
				deleteTag.DELETE(":id", c.TagHandlers.Delete)
			}
		}

		// Category routes
		category := metadata.Group("/category")
		{
			readCategory := category.Group("")
			readCategory.Use(c.KeycloakModule.Middleware("administrator"))
			{
				readCategory.GET("", c.CategoryHandlers.GetAll)
				readCategory.GET(":id", c.CategoryHandlers.Get)
			}

			createCategory := category.Group("")
			createCategory.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createCategory.POST("", c.CategoryHandlers.Create)
			}

			updateCategory := category.Group("")
			updateCategory.Use(c.KeycloakModule.Middleware("administrator"))
			{
				updateCategory.PUT(":id", c.CategoryHandlers.Update)
			}

			deleteCategory := category.Group("")
			deleteCategory.Use(c.KeycloakModule.Middleware("administrator"))
			{
				deleteCategory.DELETE(":id", c.CategoryHandlers.Delete)
			}
		}

		// CuisineType routes
		cuisineType := metadata.Group("/cuisinetype")
		{
			readCuisineType := cuisineType.Group("")
			readCuisineType.Use(c.KeycloakModule.Middleware("administrator"))
			{
				readCuisineType.GET("", c.CuisineTypeHandlers.GetAll)
				readCuisineType.GET(":id", c.CuisineTypeHandlers.Get)
			}

			createCuisineType := cuisineType.Group("")
			createCuisineType.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createCuisineType.POST("", c.CuisineTypeHandlers.Create)
			}

			updateCuisineType := cuisineType.Group("")
			updateCuisineType.Use(c.KeycloakModule.Middleware("administrator"))
			{
				updateCuisineType.PUT(":id", c.CuisineTypeHandlers.Update)
			}

			deleteCuisineType := cuisineType.Group("")
			deleteCuisineType.Use(c.KeycloakModule.Middleware("administrator"))
			{
				deleteCuisineType.DELETE(":id", c.CuisineTypeHandlers.Delete)
			}
		}

		// DifficultyLevel routes
		DifficultyLevel := metadata.Group("/difficultylevel")
		{
			readDifficultyLevel := DifficultyLevel.Group("")
			readDifficultyLevel.Use(c.KeycloakModule.Middleware("administrator"))
			{
				readDifficultyLevel.GET("", c.DifficultyLevelHandlers.GetAll)
				readDifficultyLevel.GET(":id", c.DifficultyLevelHandlers.Get)
			}

			createDifficultyLevel := DifficultyLevel.Group("")
			createDifficultyLevel.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createDifficultyLevel.POST("", c.DifficultyLevelHandlers.Create)
			}

			updateDifficultyLevel := DifficultyLevel.Group("")
			updateDifficultyLevel.Use(c.KeycloakModule.Middleware("administrator"))
			{
				updateDifficultyLevel.PUT(":id", c.DifficultyLevelHandlers.Update)
			}

			deleteDifficultyLevel := DifficultyLevel.Group("")
			deleteDifficultyLevel.Use(c.KeycloakModule.Middleware("administrator"))
			{
				deleteDifficultyLevel.DELETE(":id", c.DifficultyLevelHandlers.Delete)
			}
		}

		// PreparationTime routes
		PreparationTime := metadata.Group("/preparationtime")
		{
			readPreparationTime := PreparationTime.Group("")
			readPreparationTime.Use(c.KeycloakModule.Middleware("administrator"))
			{
				readPreparationTime.GET("", c.PreparationTimeHandlers.GetAll)
				readPreparationTime.GET(":id", c.PreparationTimeHandlers.Get)
			}

			createPreparationTime := PreparationTime.Group("")
			createPreparationTime.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createPreparationTime.POST("", c.PreparationTimeHandlers.Create)
			}

			updatePreparationTime := PreparationTime.Group("")
			updatePreparationTime.Use(c.KeycloakModule.Middleware("administrator"))
			{
				updatePreparationTime.PUT(":id", c.PreparationTimeHandlers.Update)
			}

			deletePreparationTime := PreparationTime.Group("")
			deletePreparationTime.Use(c.KeycloakModule.Middleware("administrator"))
			{
				deletePreparationTime.DELETE(":id", c.PreparationTimeHandlers.Delete)
			}
		}

		// Search routes
		search := metadata.Group("/search")
		{
			all := search.Group("/all")
			all.Use(c.KeycloakModule.Middleware("administrator"))
			{
				all.GET("", c.SearchHandlers.GetAllMetadata)
			}

			createSearch := search.Group("")
			createSearch.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createSearch.POST("", c.SearchHandlers.SearchMetadata)
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
		err := srv.Shutdown(ctx)
		if err != nil {
			return
		}
	}()

	log.Infof("metadata service available on port %s", srv.Addr)
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error(err)
	}
}
