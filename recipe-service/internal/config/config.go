package config

import (
	hh "recipe-service/internal/handlers/http"
	rh "recipe-service/internal/handlers/rabbitmq"
	dr "recipe-service/internal/repositories/database"
	s "recipe-service/internal/services"
	"time"

	m "github.com/ihulsbus/cookbook/shared/models"

	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/oidc"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	Configuration m.Config
	err           error

	Logger            *log.Logger = log.New()
	DatabaseClient    *gorm.DB
	OIDCModule        *oidc.Module
	Cors              cors.Config
	RabbitMQClient    *rmq.RabbitMQ
	RabbitMQPublisher *rmq.Publisher
	OutboxRelay       *rmq.Relay

	// Repositories
	DatabaseRepository *dr.DatabaseRepository

	// Services
	RecipeService *s.RecipeService

	// Handlers
	HttpHandler     *hh.HttpHandlers
	InternalHandler *hh.InternalHandlers
	RabbitMQHandler *rh.RabbitMQHandler
	HealthHandler   *healthh.Handlers
)

func init() {
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

	if Configuration.Global.InternalListenPort == 0 {
		Logger.Warn("Internal listen port is empty. Defaulting to 9090")
		Configuration.Global.InternalListenPort = 9090
	}

	initDatabase()
	initCors()
	OIDCModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}
	RabbitMQClient, err = rmq.NewRabbitMQConnection(
		Configuration.RabbitMQ.Username,
		Configuration.RabbitMQ.Password,
		Configuration.RabbitMQ.Host,
		Logger,
	)

	// Init repositories
	DatabaseRepository = dr.NewDatabaseRepository(DatabaseClient)

	// Recipe events are written to the outbox table in the same DB
	// transaction as the change (DatabaseRepository.Create/Update/Delete)
	// and relayed to RabbitMQ in the background - see OutboxRelay below.
	// RabbitMQPublisher is the Relay's only use of the connection; it is
	// not called directly by the service layer.
	RabbitMQPublisher, err = rmq.NewPublisher(RabbitMQClient.Connection, "cookbook")
	if err != nil {
		Logger.Errorf("Error setting up RabbitMQ publisher: %v", err)
		Logger.Fatal("Encountered fatal error. Exiting.")
	}
	OutboxRelay = rmq.NewRelay(DatabaseClient, RabbitMQPublisher, Logger, 2*time.Second, 50)

	// Init services
	RecipeService = s.NewRecipeService(DatabaseRepository, Logger)

	// Init handlers
	HttpHandler = hh.NewHttpHandlers(RecipeService, Logger)
	InternalHandler = hh.NewInternalHandlers(RecipeService, Logger)
	RabbitMQHandler, err = rh.NewRabbitMQHandler(RecipeService, Logger)
	if err != nil {
		Logger.Errorf("Error setting up RabbitMQ Consumer: %v", err)
		Logger.Fatal("Encountered fatal error. Exiting.")
	}
	HealthHandler = healthh.NewHealthHandlers(DatabaseClient, Logger)
}
