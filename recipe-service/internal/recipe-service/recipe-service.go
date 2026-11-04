package recipeservice

import (
	"context"
	"net/http"
	c "recipe-service/internal/config"
	m "recipe-service/internal/middleware"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	log = c.Logger
)

func RecipeService(ctx context.Context) {
	err := c.RabbitMQHandler.StartConsuming(c.RabbitMQClient.Connection, "recipe", "cookbook")
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

	v2 := router.Group("/api/v2")
	{
		// "full CRUD on their own recipes" - any authenticated user
		// (per-recipe ownership checks come from the ownership epic).
		recipe := v2.Group("/recipe")
		{
			readRecipe := recipe.Group("")
			readRecipe.Use(c.OIDCModule.RequireAuth())
			{
				readRecipe.GET("", c.HttpHandler.GetAll)
				readRecipe.GET(":id", c.HttpHandler.Get)
			}

			createRecipe := recipe.Group("")
			createRecipe.Use(c.OIDCModule.RequireAuth())
			{
				createRecipe.POST("", c.HttpHandler.Create)
			}

			updateRecipe := recipe.Group("")
			updateRecipe.Use(c.OIDCModule.RequireAuth())
			{

				updateRecipe.PUT(":id", c.HttpHandler.Update)
			}

			adminRecipe := recipe.Group("")
			adminRecipe.Use(c.OIDCModule.RequireAuth())
			{
				adminRecipe.DELETE(":id", c.HttpHandler.Delete)
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
		log.Info("Stopping webserver")
		srv.Shutdown(ctx)
	}()

	log.Infof("recipe service available on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
}
