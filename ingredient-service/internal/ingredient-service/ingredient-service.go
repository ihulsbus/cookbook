package ingredientservice

import (
	"context"
	c "ingredient-service/internal/config"
	m "ingredient-service/internal/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func IngredientService(ctx context.Context) {
	err := c.RabbitMQHandler.StartConsuming(c.RabbitMQClient.Connection, "ingredients", "cookbook")
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
	router.Use(m.Logger(c.Logger))

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

	v2 := router.Group("/api/v2")
	{
		amount := v2.Group("/amount")
		{
			readAmounts := amount.Group("")
			readAmounts.Use(c.KeycloakModule.Middleware("administrator"))
			{
				readAmounts.GET(":id", c.AmountHandlers.Find)
			}
			createAmounts := amount.Group("")
			createAmounts.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createAmounts.POST(":id", c.AmountHandlers.Create)
			}

			updateAmounts := amount.Group("")
			updateAmounts.Use(c.KeycloakModule.Middleware("administrator"))
			{
				updateAmounts.PUT(":id", c.AmountHandlers.Update)
			}

			deleteAmounts := amount.Group("")
			deleteAmounts.Use(c.KeycloakModule.Middleware("administrator"))
			{
				deleteAmounts.DELETE(":id", c.AmountHandlers.Delete)
			}
		}
		ingredient := v2.Group("/ingredient")
		{
			readIngredient := ingredient.Group("")
			readIngredient.Use(c.KeycloakModule.Middleware("administrator"))
			{
				readIngredient.GET("", c.IngredientHandlers.GetAll)
				readIngredient.GET(":id", c.IngredientHandlers.GetSingle)
			}

			createIngredient := ingredient.Group("")
			createIngredient.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createIngredient.POST("", c.IngredientHandlers.Create)
			}

			updateIngredient := ingredient.Group("")
			updateIngredient.Use(c.KeycloakModule.Middleware("administrator"))
			{
				updateIngredient.PUT(":id", c.IngredientHandlers.Update)
			}

			adminIngredient := ingredient.Group("")
			adminIngredient.Use(c.KeycloakModule.Middleware("administrator"))
			{
				adminIngredient.DELETE(":id", c.IngredientHandlers.Delete)
			}
		}
		unit := v2.Group("/unit")
		{
			readUnit := unit.Group("")
			readUnit.Use(c.KeycloakModule.Middleware("administrator"))
			{
				readUnit.GET("", c.UnitHandlers.GetAll)
				readUnit.GET(":id", c.UnitHandlers.GetSingle)
			}

			createUnit := unit.Group("")
			createUnit.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createUnit.POST("", c.UnitHandlers.Create)
			}

			updateUnit := unit.Group("")
			updateUnit.Use(c.KeycloakModule.Middleware("administrator"))
			{
				updateUnit.PUT(":id", c.UnitHandlers.Update)
			}

			deleteUnit := unit.Group("")
			deleteUnit.Use(c.KeycloakModule.Middleware("administrator"))
			{
				deleteUnit.DELETE(":id", c.UnitHandlers.Delete)
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
		srv.Shutdown(ctx)
	}()

	c.Logger.Infof("ingredient service available on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		c.Logger.Error(err)
	}
}
