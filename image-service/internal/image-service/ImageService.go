package instructionservice

import (
	"context"
	c "image-service/internal/config"
	"net/http"
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

func ImageService(ctx context.Context) {
	router := gin.New()
	gin.SetMode(gin.ReleaseMode)

	// Logging
	router.Use(ginglog.Logger(3 * time.Second))

	// Panic recovery
	router.Use(gin.Recovery())

	// Cors handler
	router.Use(cors.New(c.Cors))

	privateGroup := router.Group("/api")
	privateGroup.Use(ginkeycloak.NewAccessBuilder(config).
		RestrictButForRole("administrator").
		Build())

	privateGroup.GET("/privategroup", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "Hello from private for groups"})
	})

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
