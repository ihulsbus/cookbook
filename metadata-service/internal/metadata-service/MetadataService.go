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
	v1 := router.Group("/api/v2/metadata")
	{

		// Tag routes
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

		// Category routes
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

		// CuisineType routes
		cuisineType := v1.Group("/cuisinetype")
		{
			readCuisineType := cuisineType.Group("")
			readCuisineType.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				readCuisineType.GET("", c.CuisineTypeHandlers.GetAll)
				readCuisineType.GET(":id", c.CuisineTypeHandlers.Get)
			}

			createCuisineType := cuisineType.Group("")
			createCuisineType.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				createCuisineType.POST("", c.CuisineTypeHandlers.Create)
			}

			updateCuisineType := cuisineType.Group("")
			updateCuisineType.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				updateCuisineType.PUT(":id", c.CuisineTypeHandlers.Update)
			}

			deleteCuisineType := cuisineType.Group("")
			deleteCuisineType.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				deleteCuisineType.DELETE(":id", c.CuisineTypeHandlers.Delete)
			}
		}

		// DifficultyLevel routes
		DifficultyLevel := v1.Group("/difficultylevel")
		{
			readDifficultyLevel := DifficultyLevel.Group("")
			readDifficultyLevel.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				readDifficultyLevel.GET("", c.DifficultyLevelHandlers.GetAll)
				readDifficultyLevel.GET(":id", c.DifficultyLevelHandlers.Get)
			}

			createDifficultyLevel := DifficultyLevel.Group("")
			createDifficultyLevel.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				createDifficultyLevel.POST("", c.DifficultyLevelHandlers.Create)
			}

			updateDifficultyLevel := DifficultyLevel.Group("")
			updateDifficultyLevel.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				updateDifficultyLevel.PUT(":id", c.DifficultyLevelHandlers.Update)
			}

			deleteDifficultyLevel := DifficultyLevel.Group("")
			deleteDifficultyLevel.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				deleteDifficultyLevel.DELETE(":id", c.DifficultyLevelHandlers.Delete)
			}
		}

		// PreparationTime routes
		PreparationTime := v1.Group("/preparationtime")
		{
			readPreparationTime := PreparationTime.Group("")
			readPreparationTime.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				readPreparationTime.GET("", c.PreparationTimeHandlers.GetAll)
				readPreparationTime.GET(":id", c.PreparationTimeHandlers.Get)
			}

			createPreparationTime := PreparationTime.Group("")
			createPreparationTime.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				createPreparationTime.POST("", c.PreparationTimeHandlers.Create)
			}

			updatePreparationTime := PreparationTime.Group("")
			updatePreparationTime.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				updatePreparationTime.PUT(":id", c.PreparationTimeHandlers.Update)
			}

			deletePreparationTime := PreparationTime.Group("")
			deletePreparationTime.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				deletePreparationTime.DELETE(":id", c.PreparationTimeHandlers.Delete)
			}
		}

		// Search routes
		search := v1.Group("/search")
		{
			createSearch := search.Group("")
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
		srv.Shutdown(ctx)
	}()

	log.Info("metadata service available on port 8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
}
