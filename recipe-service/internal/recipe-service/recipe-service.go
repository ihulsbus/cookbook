package recipeservice

import (
	"context"
	"net/http"
	c "recipe-service/internal/config"
	m "recipe-service/internal/middleware"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
)

var (
	log = c.Logger
)

func RecipeService(ctx context.Context) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	// Logging
	router.Use(m.Logger(log))

	// Panic recovery
	router.Use(gin.Recovery())

	// Cors handler
	router.Use(cors.New(c.Cors))

	v1 := router.Group("/api/v2")
	{
		recipe := v1.Group("/recipes")
		{
			readRecipe := recipe.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				readRecipe.GET("", c.RecipeEndpoints.GetAll)
				readRecipe.GET(":id", c.RecipeEndpoints.Get)
			}

			createRecipe := recipe.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				createRecipe.POST("", c.RecipeEndpoints.Create)
			}

			updateRecipe := recipe.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{

				updateRecipe.PUT(":id", c.RecipeEndpoints.Update)
			}

			adminRecipe := recipe.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				adminRecipe.DELETE(":id", c.RecipeEndpoints.Delete)
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

	log.Info("recipe service available on port 8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
}
