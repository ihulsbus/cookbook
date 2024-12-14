package config

import (
	httpHandler "image-service/internal/handlers/http"
	rabbitMQHandler "image-service/internal/handlers/rabbitmq"
	m "image-service/internal/models"
	ir "image-service/internal/repositories/database"
	sr "image-service/internal/repositories/s3"
	s "image-service/internal/services"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/wagslane/go-rabbitmq"
	"gorm.io/gorm"
)

var (
	Configuration m.Config
	err           error

	Logger         *log.Logger = log.New()
	DatabaseClient *gorm.DB
	S3Client       *s3.S3
	Cors           cors.Config
	RabbitMQClient *rabbitmq.Conn

	// Repositories
	ImageRepository *ir.DatabaseRepository
	S3Repository    *sr.S3Repository

	// Services
	ImageService *s.ImageService

	// Handlers
	ImageHandler    *httpHandler.ImageHandlers
	RabbitMQHandler *rabbitMQHandler.RabbitMQConsumer
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
	initRabbitMQ(
		Configuration.RabbitMQ.Username,
		Configuration.RabbitMQ.Password,
		Configuration.RabbitMQ.Host,
	)

	// Init repositories
	ImageRepository = ir.NewDatabaseRepository(DatabaseClient)
	S3Repository = sr.NewS3Repository(S3Client, Logger, Configuration.S3.BucketName)

	// Init services
	ImageService = s.NewImageService(ImageRepository, S3Repository, Logger)

	// Init handlers
	ImageHandler = httpHandler.NewImageHandlers(ImageService, Logger)
	RabbitMQHandler, err = rabbitMQHandler.NewRabbitMQConsumer(ImageService, Logger)
	if err != nil {
		Logger.Errorf("Error setting up RabbitMQ Consumer: %v", err)
		Logger.Fatal("Encountered fatal error. Exiting.")
	}
}
