package instructionservice

import (
	"context"
	c "instruction-service/internal/config"
	m "instruction-service/internal/middleware"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	log            = c.Logger
	validAudiences []string
)

func InstructionService(ctx context.Context) {
	err := c.RabbitMQHandler.StartConsuming(c.RabbitMQClient.Connection, "instructions", "cookbook")
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
		health.GET("/live", c.HealthHandlers.Liveness)
		health.GET("/ready", c.HealthHandlers.Readiness)
		health.GET("/startup", c.HealthHandlers.Startup)
	}

	// API versioning setup
	v2 := router.Group("/api/v2")
	{
		instruction := v2.Group("/instruction")
		{
			searchInstruction := instruction.Group("")
			searchInstruction.Use(c.KeycloakModule.Middleware("administrator"))
			{
				searchInstruction.GET("/search", c.SearchHandlers.SearchInstruction)
			}

			readInstruction := instruction.Group("")
			readInstruction.Use(c.KeycloakModule.Middleware("administrator"))
			{
				readInstruction.GET(":id", c.InstructionHandlers.Get)
			}

			createInstruction := instruction.Group("")
			createInstruction.Use(c.KeycloakModule.Middleware("administrator"))
			{
				createInstruction.POST(":id", c.InstructionHandlers.Create)
			}

			updateInstruction := instruction.Group("")
			updateInstruction.Use(c.KeycloakModule.Middleware("administrator"))
			{
				updateInstruction.PUT(":id", c.InstructionHandlers.Update)
			}

			deleteInstruction := instruction.Group("")
			deleteInstruction.Use(c.KeycloakModule.Middleware("administrator"))
			{
				deleteInstruction.DELETE(":id", c.InstructionHandlers.Delete)
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

	log.Infof("instruction service available on port %s", srv.Addr)
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
}
