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

		// Recipe metadata is the per-recipe assignment of tags/category/
		// etc, not a catalog lookup table itself - CRUD follows "full
		// CRUD on their own recipes", any authenticated user.
		recipe := metadata.Group("/recipe")
		{
			readMetadata := recipe.Group("")
			readMetadata.Use(c.OIDCModule.RequireAuth())
			{
				readMetadata.GET("", c.MetadataHandlers.GetAll)
				readMetadata.GET(":id", c.MetadataHandlers.Get)
			}
			createMetadata := recipe.Group("")
			createMetadata.Use(c.OIDCModule.RequireAuth())
			{
				createMetadata.POST(":id", c.MetadataHandlers.Create)
			}

			updateMetadata := recipe.Group("")
			updateMetadata.Use(c.OIDCModule.RequireAuth())
			{
				updateMetadata.PUT(":id", c.MetadataHandlers.Update)
			}

			deleteMetadata := recipe.Group("")
			deleteMetadata.Use(c.OIDCModule.RequireAuth())
			{
				deleteMetadata.DELETE(":id", c.MetadataHandlers.Delete)
			}
		}

		// Tag routes - metadata lookup table: any user reads, only
		// cookbook-admin writes.
		tag := metadata.Group("/tag")
		{
			readTag := tag.Group("")
			readTag.Use(c.OIDCModule.RequireAuth())
			{
				readTag.GET("", c.TagHandlers.GetAll)
				readTag.GET(":id", c.TagHandlers.Get)
			}

			createTag := tag.Group("")
			createTag.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				createTag.POST("", c.TagHandlers.Create)
			}

			updateTag := tag.Group("")
			updateTag.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				updateTag.PUT(":id", c.TagHandlers.Update)
			}

			deleteTag := tag.Group("")
			deleteTag.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				deleteTag.DELETE(":id", c.TagHandlers.Delete)
			}
		}

		// Category routes - metadata lookup table: any user reads, only
		// cookbook-admin writes.
		category := metadata.Group("/category")
		{
			readCategory := category.Group("")
			readCategory.Use(c.OIDCModule.RequireAuth())
			{
				readCategory.GET("", c.CategoryHandlers.GetAll)
				readCategory.GET(":id", c.CategoryHandlers.Get)
			}

			createCategory := category.Group("")
			createCategory.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				createCategory.POST("", c.CategoryHandlers.Create)
			}

			updateCategory := category.Group("")
			updateCategory.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				updateCategory.PUT(":id", c.CategoryHandlers.Update)
			}

			deleteCategory := category.Group("")
			deleteCategory.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				deleteCategory.DELETE(":id", c.CategoryHandlers.Delete)
			}
		}

		// CuisineType routes - metadata lookup table: any user reads,
		// only cookbook-admin writes.
		cuisineType := metadata.Group("/cuisinetype")
		{
			readCuisineType := cuisineType.Group("")
			readCuisineType.Use(c.OIDCModule.RequireAuth())
			{
				readCuisineType.GET("", c.CuisineTypeHandlers.GetAll)
				readCuisineType.GET(":id", c.CuisineTypeHandlers.Get)
			}

			createCuisineType := cuisineType.Group("")
			createCuisineType.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				createCuisineType.POST("", c.CuisineTypeHandlers.Create)
			}

			updateCuisineType := cuisineType.Group("")
			updateCuisineType.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				updateCuisineType.PUT(":id", c.CuisineTypeHandlers.Update)
			}

			deleteCuisineType := cuisineType.Group("")
			deleteCuisineType.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				deleteCuisineType.DELETE(":id", c.CuisineTypeHandlers.Delete)
			}
		}

		// DifficultyLevel routes - metadata lookup table: any user reads,
		// only cookbook-admin writes.
		DifficultyLevel := metadata.Group("/difficultylevel")
		{
			readDifficultyLevel := DifficultyLevel.Group("")
			readDifficultyLevel.Use(c.OIDCModule.RequireAuth())
			{
				readDifficultyLevel.GET("", c.DifficultyLevelHandlers.GetAll)
				readDifficultyLevel.GET(":id", c.DifficultyLevelHandlers.Get)
			}

			createDifficultyLevel := DifficultyLevel.Group("")
			createDifficultyLevel.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				createDifficultyLevel.POST("", c.DifficultyLevelHandlers.Create)
			}

			updateDifficultyLevel := DifficultyLevel.Group("")
			updateDifficultyLevel.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				updateDifficultyLevel.PUT(":id", c.DifficultyLevelHandlers.Update)
			}

			deleteDifficultyLevel := DifficultyLevel.Group("")
			deleteDifficultyLevel.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				deleteDifficultyLevel.DELETE(":id", c.DifficultyLevelHandlers.Delete)
			}
		}

		// PreparationTime routes - metadata lookup table: any user
		// reads, only cookbook-admin writes.
		PreparationTime := metadata.Group("/preparationtime")
		{
			readPreparationTime := PreparationTime.Group("")
			readPreparationTime.Use(c.OIDCModule.RequireAuth())
			{
				readPreparationTime.GET("", c.PreparationTimeHandler.GetAll)
				readPreparationTime.GET(":id", c.PreparationTimeHandler.Get)
			}

			createPreparationTime := PreparationTime.Group("")
			createPreparationTime.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				createPreparationTime.POST("", c.PreparationTimeHandler.Create)
			}

			updatePreparationTime := PreparationTime.Group("")
			updatePreparationTime.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				updatePreparationTime.PUT(":id", c.PreparationTimeHandler.Update)
			}

			deletePreparationTime := PreparationTime.Group("")
			deletePreparationTime.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				deletePreparationTime.DELETE(":id", c.PreparationTimeHandler.Delete)
			}
		}

		// Search is a read operation (filters/looks up metadata) - any
		// authenticated user.
		search := metadata.Group("/search")
		{
			all := search.Group("/all")
			all.Use(c.OIDCModule.RequireAuth())
			{
				all.GET("", c.SearchHandlers.GetAllMetadata)
			}

			createSearch := search.Group("")
			createSearch.Use(c.OIDCModule.RequireAuth())
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
