package config

import (
	"context"
	"time"

	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"

	ih "instruction-service/internal/handlers/InstructionHandlers"
	sh "instruction-service/internal/handlers/SearchHandlers"

	ir "instruction-service/internal/repositories/instructions"
	sr "instruction-service/internal/repositories/search"

	az "github.com/ihulsbus/cookbook/shared/authz"
	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"
	hc "github.com/ihulsbus/cookbook/shared/httpclient"

	is "instruction-service/internal/services/InstructionService"
	ss "instruction-service/internal/services/SearchService"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/oidc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type instructionConfig struct {
	m.Config `mapstructure:",squash"`
	// RecipeServiceInternal is used only by Authorizer's projection-miss fallback.
	RecipeServiceInternal m.ApiClient
}

var (
	err           error
	Ctx           context.Context
	Configuration instructionConfig

	Logger     = log.New()
	OIDCModule *oidc.Module
	Cors       cors.Config

	DatabaseClient *gorm.DB
	RabbitMQClient *rmq.RabbitMQ

	InstructionRepository *ir.InstructionRepository
	SearchRepository      *sr.SearchRepository

	InstructionService *is.InstructionService
	SearchService      *ss.SearchService

	// Authorizer answers CanRead/CanWrite against a local projection kept
	// current by Projector's RabbitMQ consumer, with a lazy fallback to
	// recipe-service's internal endpoint on a miss.
	Authorizer *az.Authorizer
	Projector  *az.Projector

	InstructionHandlers *ih.InstructionHandlers
	SearchHandlers      *sh.SearchHandlers
	HealthHandlers      *healthh.Handlers
)

func init() {
	Ctx = context.Background()

	initViper()
	initConfig()
	initLogging()

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("config file changed: %s", e.Name)

		initConfig()
		initLogging()
	})

	if Configuration.Global.ListenPort == 0 {
		Logger.Warn("Listen port is empty. Defaulting to 8080")
		Configuration.Global.ListenPort = 8080
	}

	initDatabase()
	initCors()
	OIDCModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}

	// RabbitMQClient.Connection is used by Projector's consumer (started in
	// InstructionService.go) - no synchronous call to any other service
	// happens during startup.
	RabbitMQClient, err = rmq.NewRabbitMQConnection(
		Configuration.RabbitMQ.Username,
		Configuration.RabbitMQ.Password,
		Configuration.RabbitMQ.Host,
		Logger,
	)

	InstructionRepository = ir.NewInstructionRepository(DatabaseClient)
	SearchRepository = sr.NewSearchRepository(DatabaseClient)

	InstructionService = is.NewInstructionService(InstructionRepository)
	SearchService = ss.NewSearchService(SearchRepository)

	Projector = az.NewProjector(DatabaseClient)
	Authorizer = az.NewAuthorizer(DatabaseClient, Logger)
	fallbackClient := hc.NewInternalHTTPClient(5*time.Second, Logger, nil, nil)
	Authorizer.SetMissFallback(az.NewRecipeServiceFallback(fallbackClient, Configuration.RecipeServiceInternal.BaseURL, Projector))

	InstructionHandlers = ih.NewInstructionHandlers(InstructionService, Logger)
	SearchHandlers = sh.NewSearchHandlers(SearchService, Logger)
	HealthHandlers = healthh.NewHealthHandlers(DatabaseClient, Logger)
}
