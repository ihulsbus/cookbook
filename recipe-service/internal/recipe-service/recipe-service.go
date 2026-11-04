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

	go c.OutboxRelay.Start(ctx)
	go internalHttpServer(ctx)

	httpServer(ctx)
}

// internalHttpServer serves recipe-service's internal-only API (KIT-28)
// on its own port, bound to its own gin.Engine with no auth/CORS
// middleware - it's for other services' shared/authz instances to call
// from inside the cluster, not for clients. This port must stay out of
// ingress.hosts[].paths in k8s/values.yaml; defense in depth beyond that
// (a NetworkPolicy restricting who can reach it) is tracked separately,
// since the shared Helm library this chart uses has no NetworkPolicy
// template today.
func internalHttpServer(ctx context.Context) {
	router := gin.New()
	router.Use(m.Logger(log))
	router.Use(gin.Recovery())

	internal := router.Group("/internal/v1")
	{
		internal.GET("/recipe/:id/access", c.InternalHandler.GetAccess)
	}

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + strconv.Itoa(c.Configuration.Global.InternalListenPort),
		WriteTimeout: 30 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		<-ctx.Done()
		log.Info("Stopping internal webserver")
		srv.Shutdown(ctx)
	}()

	log.Infof("recipe service internal API available on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
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
		// Every route here only requires a valid token; per-recipe
		// ownership/visibility enforcement (KIT-25) happens inside the
		// service layer, since it needs each recipe's owner_id and
		// visibility to decide, not just the caller's identity.
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
