package config

import (
	"context"
	rh "instruction-service/internal/handlers/RabbitmqHandlers"
	cs "instruction-service/internal/services/CacheService"
	"time"

	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"

	ih "instruction-service/internal/handlers/InstructionHandlers"
	sh "instruction-service/internal/handlers/SearchHandlers"

	ir "instruction-service/internal/repositories/instructions"
	sr "instruction-service/internal/repositories/search"

	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"
	hc "github.com/ihulsbus/cookbook/shared/httpclient"
	rc "github.com/ihulsbus/cookbook/shared/recipeclient"

	is "instruction-service/internal/services/InstructionService"
	ss "instruction-service/internal/services/SearchService"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/keycloak"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type instructionConfig struct {
	m.Config     `mapstructure:",squash"`
	RecipeClient m.ApiClient
}

var (
	err           error
	Ctx           context.Context
	Configuration instructionConfig

	Logger         = log.New()
	KeycloakModule *keycloak.KeycloakModule
	Cors           cors.Config

	// Clients
	DatabaseClient *gorm.DB
	HttpClient     *hc.HTTPClient
	RabbitMQClient *rmq.RabbitMQ
	RecipeClient   *rc.RecipeAPIClient

	// Repositories
	InstructionRepository *ir.InstructionRepository
	SearchRepository      *sr.SearchRepository

	// Services
	InstructionService *is.InstructionService
	SearchService      *ss.SearchService
	CacheService       *cs.CacheService

	// Handlers
	InstructionHandlers *ih.InstructionHandlers
	SearchHandlers      *sh.SearchHandlers
	HealthHandlers      *healthh.Handlers
	RabbitMQHandler     *rh.RabbitMQHandler
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
	KeycloakModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}

	// Init clients
	HttpClient = hc.NewHTTPClient(5*time.Second, Configuration.Oauth.Url, Configuration.Oauth.Realm, Configuration.Oauth.ClientID, Configuration.Oauth.ClientSecret, Logger)
	RecipeClient, err = rc.NewRecipeAPIClient(Configuration.RecipeClient.BaseURL, HttpClient)
	if err != nil {
		Logger.Panicf("error initialising recipe client: %v", err)
	}
	RabbitMQClient, err = rmq.NewRabbitMQConnection(
		Configuration.RabbitMQ.Username,
		Configuration.RabbitMQ.Password,
		Configuration.RabbitMQ.Host,
		Logger,
	)

	// Init repositories
	InstructionRepository = ir.NewInstructionRepository(DatabaseClient)
	SearchRepository = sr.NewSearchRepository(DatabaseClient)

	// Init services
	CacheService, err = cs.NewCacheService(Ctx, RecipeClient, Logger, Configuration.Global.ListenPort+1000)
	if err != nil {
		Logger.Fatalf("Error setting up cache: %v", err)
	}
	InstructionService = is.NewInstructionService(InstructionRepository, RecipeClient)
	SearchService = ss.NewSearchService(SearchRepository)

	// Init handlers
	InstructionHandlers = ih.NewInstructionHandlers(InstructionService, Logger)
	SearchHandlers = sh.NewSearchHandlers(SearchService, Logger)
	HealthHandlers = healthh.NewHealthHandlers(DatabaseClient, Logger)
	RabbitMQHandler, err = rh.NewRabbitMQHandler(CacheService, &Ctx, Logger)
	if err != nil {
		Logger.Fatalf("Error setting up RabbitMQ Consumer: %v", err)
	}
}
