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
	az "github.com/ihulsbus/cookbook/shared/authz"
)

func IngredientService(ctx context.Context) {
	// Keeps the recipe_access projection current; this is the only
	// cross-service data flow at startup, and it's event-driven (no
	// synchronous call to recipe-service).
	consumer, err := c.Projector.StartConsumer(c.RabbitMQClient.Connection, "ingredient-service-recipe-access", "cookbook")
	if err != nil {
		c.Logger.Fatalf("Startup of RabbitMQ Consumer encountered fatal error: %v", err.Error())
		return
	}
	defer consumer.Close()

	httpServer(ctx)
}

func httpServer(ctx context.Context) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	router.Use(m.Logger(c.Logger))
	router.Use(gin.Recovery())
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
		// Amounts are recipe-scoped: :id here is the RECIPE id, so
		// RequireRecipeAccess checks it directly with no extra lookup.
		amount := v2.Group("/amount")
		{
			readAmounts := amount.Group("")
			readAmounts.Use(c.OIDCModule.RequireAuth(), c.Authorizer.RequireRecipeAccess("id", az.PermissionRead))
			{
				readAmounts.GET(":id", c.AmountHandlers.Find)
			}
			createAmounts := amount.Group("")
			createAmounts.Use(c.OIDCModule.RequireAuth(), c.Authorizer.RequireRecipeAccess("id", az.PermissionWrite))
			{
				createAmounts.POST(":id", c.AmountHandlers.Create)
			}

			updateAmounts := amount.Group("")
			updateAmounts.Use(c.OIDCModule.RequireAuth(), c.Authorizer.RequireRecipeAccess("id", az.PermissionWrite))
			{
				updateAmounts.PUT(":id", c.AmountHandlers.Update)
			}

			deleteAmounts := amount.Group("")
			deleteAmounts.Use(c.OIDCModule.RequireAuth(), c.Authorizer.RequireRecipeAccess("id", az.PermissionWrite))
			{
				deleteAmounts.DELETE(":id", c.AmountHandlers.Delete)
			}
		}
		// Ingredients are a shared catalog: any user can read and create
		// (duplicate prevention guards the catalog instead of an admin
		// gate), but only cookbook-admin can update/delete entries.
		ingredient := v2.Group("/ingredient")
		{
			readIngredient := ingredient.Group("")
			readIngredient.Use(c.OIDCModule.RequireAuth())
			{
				readIngredient.GET("", c.IngredientHandlers.GetAll)
				// Static segment must be registered before :id so it isn't swallowed by the param route.
				readIngredient.GET("suggest", c.IngredientHandlers.Suggest)
				readIngredient.GET(":id", c.IngredientHandlers.GetSingle)
			}

			createIngredient := ingredient.Group("")
			createIngredient.Use(c.OIDCModule.RequireAuth())
			{
				createIngredient.POST("", c.IngredientHandlers.Create)
			}

			updateIngredient := ingredient.Group("")
			updateIngredient.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				updateIngredient.PUT(":id", c.IngredientHandlers.Update)
			}

			adminIngredient := ingredient.Group("")
			adminIngredient.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				adminIngredient.DELETE(":id", c.IngredientHandlers.Delete)
			}
		}
		// Units are catalog-admin-only for all writes; any user can read
		// them.
		unit := v2.Group("/unit")
		{
			readUnit := unit.Group("")
			readUnit.Use(c.OIDCModule.RequireAuth())
			{
				readUnit.GET("", c.UnitHandlers.GetAll)
				readUnit.GET(":id", c.UnitHandlers.GetSingle)
			}

			createUnit := unit.Group("")
			createUnit.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				createUnit.POST("", c.UnitHandlers.Create)
			}

			updateUnit := unit.Group("")
			updateUnit.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				updateUnit.PUT(":id", c.UnitHandlers.Update)
			}

			deleteUnit := unit.Group("")
			deleteUnit.Use(c.OIDCModule.RequireRole("cookbook-admin"))
			{
				deleteUnit.DELETE(":id", c.UnitHandlers.Delete)
			}
		}
	}

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
