package config

import (
	hh "image-service/internal/handlers/http"
	rh "image-service/internal/handlers/rabbitmq"
	m "image-service/internal/models"
	dr "image-service/internal/repositories/database"
	rr "image-service/internal/repositories/rabbitmq"
	sr "image-service/internal/repositories/s3"
	s "image-service/internal/services"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
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
	S3Client       *s3.S3
	Cors           cors.Config
	RabbitMQClient *rmq.RabbitMQ

	// Repositories
	DatabaseRepository *dr.DatabaseRepository
	RabbitMQRepository *rr.RabbitMQRepository
	S3Repository       *sr.S3Repository

	// Services
	ImageService *s.ImageService

	// Handlers
	HttpHandler     *hh.HttpHandler
	RabbitMQHandler *rh.RabbitMQHandler
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

	initDatabase()
	S3Client = initS3(
		Configuration.S3.Endpoint,
		Configuration.S3.AWSAccessSecret,
		Configuration.S3.AWSAccessKey,
		"us-east-1",
	)
	initCors()
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

	S3Repository = sr.NewS3Repository(S3Client, Logger, Configuration.S3.BucketName)

	// Init services
	ImageService = s.NewImageService(DatabaseRepository, RabbitMQRepository, S3Repository, Logger)

	// Init handlers
	HttpHandler = hh.NewHttpHandler(ImageService, Logger)
	RabbitMQHandler, err = rh.NewRabbitMQHandler(ImageService, Logger)
	if err != nil {
		Logger.Errorf("Error setting up RabbitMQ Consumer: %v", err)
		Logger.Fatal("Encountered fatal error. Exiting.")
	}
}
