package config

import (
	"context"
	hh "image-service/internal/handlers/http"
	dr "image-service/internal/repositories/database"
	rr "image-service/internal/repositories/rabbitmq"
	sr "image-service/internal/repositories/s3"
	is "image-service/internal/services/ImageService"
	"time"

	az "github.com/ihulsbus/cookbook/shared/authz"
	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"
	hc "github.com/ihulsbus/cookbook/shared/httpclient"
	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/oidc"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type imageConfig struct {
	m.Config `mapstructure:",squash"`
	// RecipeServiceInternal is used only by Authorizer's projection-miss fallback.
	RecipeServiceInternal m.ApiClient
}

var (
	Configuration imageConfig
	Ctx           context.Context
	err           error

	Logger     *log.Logger = log.New()
	OIDCModule *oidc.Module
	Cors       cors.Config

	S3Client       *s3.S3
	DatabaseClient *gorm.DB
	RabbitMQClient *rmq.RabbitMQ

	DatabaseRepository *dr.DatabaseRepository
	RabbitMQRepository *rr.RabbitMQRepository
	S3Repository       *sr.S3Repository

	ImageService *is.ImageService

	// Mirrors ingredient-service's config.go pattern.
	Authorizer *az.Authorizer
	Projector  *az.Projector

	HttpHandler   *hh.HttpHandler
	HealthHandler *healthh.Handlers
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

	initCors()

	OIDCModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}

	initDatabase()
	S3Client = initS3(
		Configuration.S3.Endpoint,
		Configuration.S3.AWSAccessSecret,
		Configuration.S3.AWSAccessKey,
		"us-east-1",
	)
	// RabbitMQClient.Connection is consumed by Projector's consumer (started in
	// ImageService.go); no synchronous call to another service happens at startup.
	RabbitMQClient, err = rmq.NewRabbitMQConnection(
		Configuration.RabbitMQ.Username,
		Configuration.RabbitMQ.Password,
		Configuration.RabbitMQ.Host,
		Logger,
	)

	DatabaseRepository = dr.NewDatabaseRepository(DatabaseClient)
	RabbitMQRepository, err = rr.NewRabbitMQRepository(RabbitMQClient.Connection, "cookbook", Logger)
	if err != nil {
		Logger.Fatalf("Error setting up RabbitMQ publisher: %v", err)
	}
	S3Repository = sr.NewS3Repository(S3Client, sr.S3Presigner{Client: S3Client}, Logger, Configuration.S3.BucketName)

	ImageService = is.NewImageService(DatabaseRepository, RabbitMQRepository, S3Repository, Logger)

	Projector = az.NewProjector(DatabaseClient)
	Authorizer = az.NewAuthorizer(DatabaseClient, Logger)
	fallbackClient := hc.NewInternalHTTPClient(5*time.Second, Logger, nil, nil)
	Authorizer.SetMissFallback(az.NewRecipeServiceFallback(fallbackClient, Configuration.RecipeServiceInternal.BaseURL, Projector))

	HttpHandler = hh.NewHttpHandler(ImageService, Authorizer, Logger)
	HealthHandler = healthh.NewHealthHandlers(DatabaseClient, Logger)
}
