package instructionservice

import (
	"context"
	c "instruction-service/internal/config"
	m "instruction-service/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	log            = c.Logger
	validAudiences []string
)

func InstructionService(ctx context.Context) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	// Logging
	router.Use(m.Logger(log))

	// Panic recovery
	router.Use(gin.Recovery())

	// Cors handler
	router.Use(cors.New(c.Cors))

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
		Addr:         ":8080",
		WriteTimeout: 300 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	go func() {
		<-ctx.Done()
		srv.Shutdown(ctx)
	}()

	log.Info("instruction service available on port 8080")
	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		log.Error(err)
	}
}
