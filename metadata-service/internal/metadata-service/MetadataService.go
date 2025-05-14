package instructionservice

import (
	"context"
	"errors"
	c "metadata-service/internal/config"
	m "metadata-service/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
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
	v2 := router.Group("/api/v2/metadata")
	{

		// Tag routes
		tag := v2.Group("/tag")
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
		category := v2.Group("/category")
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
		cuisineType := v2.Group("/cuisinetype")
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
		DifficultyLevel := v2.Group("/difficultylevel")
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

		// Search routes
		search := v2.Group("/search")
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
		Addr:         ":8080",
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

	log.Info("metadata service available on port 8080")
	if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Error(err)
	}
}
