package instructionservice

import (
	"context"
	c "instruction-service/internal/config"
	m "instruction-service/internal/middleware"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/tbaehler/gin-keycloak/pkg/ginkeycloak"
)

var (
	log = c.Logger
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
	v1 := router.Group("/api/v1")
	{
		recipe := v1.Group("/recipe")
		{
			readInstruction := recipe.Group("")
			readInstruction.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				readInstruction.GET(":id/instruction", c.InstructionEndpoints.GetInstruction)
			}

			createInstruction := recipe.Group("")
			createInstruction.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				createInstruction.POST(":id/instruction", c.InstructionEndpoints.CreateInstruction)
			}

			updateInstruction := recipe.Group("")
			updateInstruction.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				updateInstruction.PUT(":id/instruction", c.InstructionEndpoints.UpdateInstruction)
			}

			adminInstruction := recipe.Group("")
			adminInstruction.Use(ginkeycloak.NewAccessBuilder(ginkeycloak.BuilderConfig(c.Configuration.Oauth)).RestrictButForRole("administrator").Build())
			{
				adminInstruction.DELETE(":id/instruction", c.InstructionEndpoints.DeleteInstruction)
			}
		}
	}

	// Server startup
	srv := &http.Server{
		Handler:      router,
		Addr:         ":8081",
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
