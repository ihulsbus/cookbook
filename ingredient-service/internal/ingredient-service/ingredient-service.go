package ingredientservice

import (
	"context"
	c "ingredient-service/internal/config"
	m "ingredient-service/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
)

var (
	log = c.Logger
)

func IngredientService(ctx context.Context) {
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
		ingredient := v1.Group("/ingredient")
		{
			readIngredient := ingredient.Group("")
			readIngredient.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				readIngredient.GET("", c.IngredientHandlers.GetAll)
				readIngredient.GET(":id", c.IngredientHandlers.GetSingle)
			}

			createIngredient := ingredient.Group("")
			readIngredient.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				createIngredient.POST("", c.IngredientHandlers.Create)
			}

			updateIngredient := ingredient.Group("")
			readIngredient.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				updateIngredient.PUT(":id", c.IngredientHandlers.Update)
			}

			adminIngredient := ingredient.Group("")
			readIngredient.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				adminIngredient.DELETE(":id", c.IngredientHandlers.Delete)
			}
		}

		unit := v1.Group("/unit")
		{
			readUnit := unit.Group("")
			readUnit.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				readUnit.GET("", c.UnitHandlers.GetAll)
				readUnit.GET(":id", c.UnitHandlers.GetSingle)
			}

			createUnit := unit.Group("")
			createUnit.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				createUnit.POST("", c.UnitHandlers.Create)
			}

			updateUnit := unit.Group("")
			updateUnit.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				updateUnit.PUT(":id", c.UnitHandlers.Update)
			}

			deleteUnit := unit.Group("")
			deleteUnit.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				deleteUnit.DELETE(":id", c.UnitHandlers.Delete)
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

	log.Info("ingredient service available on port 8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
}
