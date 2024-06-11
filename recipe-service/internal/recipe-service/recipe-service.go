package recipeservice

import (
	"context"
	"net/http"
	c "recipe-service/internal/config"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	ginglog "github.com/szuecs/gin-glog"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
)

var (
	log = c.Logger

	config = ginkeycloak.BuilderConfig{
		Service: "",
		Url:     "",
		Realm:   "",
	}
)

func RecipeService(ctx context.Context) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	// Logging
	router.Use(ginglog.Logger(3 * time.Second))

	// Panic recovery
	router.Use(gin.Recovery())

	// Cors handler
	router.Use(cors.New(c.Cors))

	v1 := router.Group("/api/v1")
	{
		recipe := v1.Group("/recipe")
		{
			readRecipe := recipe.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(config).RestrictButForRole("administrator").Build())
			{
				readRecipe.GET("", c.RecipeEndpoints.GetAll)
				readRecipe.GET(":id", c.RecipeEndpoints.Get)
			}

			createRecipe := recipe.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(config).RestrictButForRole("administrator").Build())
			{
				createRecipe.POST("", c.RecipeEndpoints.Create)
			}

			updateRecipe := recipe.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(config).RestrictButForRole("administrator").Build())
			{

				updateRecipe.PUT(":id", c.RecipeEndpoints.Update)
			}

			adminRecipe := recipe.Group("")
			readRecipe.Use(ginkeycloak.NewAccessBuilder(config).RestrictButForRole("administrator").Build())
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

	log.Info("server available on port 8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
}
