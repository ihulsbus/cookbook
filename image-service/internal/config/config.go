package config

import (
	"context"
	hh "image-service/internal/handlers/http"
	rh "image-service/internal/handlers/rabbitmq"
	dr "image-service/internal/repositories/database"
	rr "image-service/internal/repositories/rabbitmq"
	sr "image-service/internal/repositories/s3"
	cs "image-service/internal/services/CacheService"
	is "image-service/internal/services/ImageService"
	"time"

	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"
	hc "github.com/ihulsbus/cookbook/shared/httpclient"
	m "github.com/ihulsbus/cookbook/shared/models"
	rc "github.com/ihulsbus/cookbook/shared/recipeclient"
	"github.com/olric-data/olric"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/keycloak"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type imageConfig struct {
	m.Config     `mapstructure:",squash"`
	RecipeClient m.ApiClient
}

var (
	Configuration imageConfig
	Ctx           context.Context
	err           error

	Logger         *log.Logger = log.New()
	KeycloakModule *keycloak.KeycloakModule
	Cors           cors.Config
	Cache          *olric.DMap

	// Clients
	HttpClient     *hc.HTTPClient
	S3Client       *s3.S3
	DatabaseClient *gorm.DB
	RabbitMQClient *rmq.RabbitMQ
	RecipeClient   *rc.RecipeAPIClient

	// Repositories
	DatabaseRepository *dr.DatabaseRepository
	RabbitMQRepository *rr.RabbitMQRepository
	S3Repository       *sr.S3Repository

	// Services
	CacheService *cs.CacheService
	ImageService *is.ImageService

	// Handlers
	HttpHandler     *hh.HttpHandler
	RabbitMQHandler *rh.RabbitMQHandler
	HealthHandler   *healthh.Handlers
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

	if Configuration.Global.ListenPort == "" {
		Logger.Warn("Listen port is empty. Defaulting to 8080")
		Configuration.Global.ListenPort = "8080"
	}

	initCors()
	if err != nil {
		Logger.Panicf("error initialising cache: %v", err)
	}

	KeycloakModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}

	// Init clients
	initDatabase()
	S3Client = initS3(
		Configuration.S3.Endpoint,
		Configuration.S3.AWSAccessSecret,
		Configuration.S3.AWSAccessKey,
		"us-east-1",
	)
	RabbitMQClient, err = rmq.NewRabbitMQConnection(
		Configuration.RabbitMQ.Username,
		Configuration.RabbitMQ.Password,
		Configuration.RabbitMQ.Host,
		Logger,
	)
	HttpClient = hc.NewHTTPClient(5*time.Second, Configuration.Oauth.Url, Configuration.Oauth.Realm, Configuration.Oauth.ClientID, Configuration.Oauth.ClientSecret, Logger)
	RecipeClient, err = rc.NewRecipeAPIClient(Configuration.RecipeClient.BaseURL, HttpClient)
	if err != nil {
		Logger.Panicf("error initialising recipe client: %v", err)
	}

	// Init repositories
	DatabaseRepository = dr.NewDatabaseRepository(DatabaseClient)
	RabbitMQRepository, err = rr.NewRabbitMQRepository(RabbitMQClient.Connection, "cookbook", Logger)
	if err != nil {
		Logger.Fatalf("Error setting up RabbitMQ publisher: %v", err)
	}
	S3Repository = sr.NewS3Repository(S3Client, Logger, Configuration.S3.BucketName)

	// Init services
	CacheService, err = cs.NewCacheService(Ctx, RecipeClient, Logger)
	if err != nil {
		Logger.Fatalf("Error setting up cache: %v", err)
	}
	ImageService = is.NewImageService(DatabaseRepository, RabbitMQRepository, S3Repository, Logger)

	// Init handlers
	HttpHandler = hh.NewHttpHandler(ImageService, Logger)
	RabbitMQHandler, err = rh.NewRabbitMQHandler(CacheService, &Ctx, Logger)
	if err != nil {
		Logger.Fatalf("Error setting up RabbitMQ Consumer: %v", err)
	}
	HealthHandler = healthh.NewHealthHandlers(DatabaseClient, Logger)
}
