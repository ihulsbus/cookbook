package config

import (
	hh "recipe-service/internal/handlers/http"
	rh "recipe-service/internal/handlers/rabbitmq"
	dr "recipe-service/internal/repositories/database"
	rr "recipe-service/internal/repositories/rabbitmq"
	s "recipe-service/internal/services"

	m "github.com/ihulsbus/cookbook/shared/models"

	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/keycloak"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	Configuration m.Config
	err           error

	Logger         *log.Logger = log.New()
	DatabaseClient *gorm.DB
	KeycloakModule *keycloak.KeycloakModule
	Cors           cors.Config
	RabbitMQClient *rmq.RabbitMQ

	// Repositories
	DatabaseRepository *dr.DatabaseRepository
	RabbitMQRepository *rr.RabbitMQRepository

	// Services
	RecipeService *s.RecipeService

	// Handlers
	HttpHandler     *hh.HttpHandlers
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

	if Configuration.Global.ListenPort == "" {
		Logger.Warn("Listen port is empty. Defaulting to 8080")
		Configuration.Global.ListenPort = "8080"
	}

	initDatabase()
	initCors()
	KeycloakModule, err = initOauth()
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
	RabbitMQRepository, err = rr.NewRabbitMQRepository(RabbitMQClient.Connection, "cookbook", Logger)
	if err != nil {
		Logger.Errorf("Error setting up RabbitMQ publisher: %v", err)
		Logger.Fatal("Encountered fatal error. Exiting.")
	}

	// Init services
	RecipeService = s.NewRecipeService(DatabaseRepository, RabbitMQRepository, Logger)

	// Init handlers
	HttpHandler = hh.NewHttpHandlers(RecipeService, Logger)
	RabbitMQHandler, err = rh.NewRabbitMQHandler(RecipeService, Logger)
	if err != nil {
		Logger.Errorf("Error setting up RabbitMQ Consumer: %v", err)
		Logger.Fatal("Encountered fatal error. Exiting.")
	}
	HealthHandler = healthh.NewHealthHandlers(DatabaseClient, Logger)
}
