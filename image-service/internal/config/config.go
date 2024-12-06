package config

import (
	h "image-service/internal/handlers"
	m "image-service/internal/models"
	ir "image-service/internal/repositories/image"
	mq "image-service/internal/repositories/rabbitmq"
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

	Logger         *log.Logger = log.New()
	DatabaseClient *gorm.DB
	S3Client       *s3.S3
	Cors           cors.Config
	RabbitMQClient *rabbitmq.Conn

	// Repositories
	ImageRepository    *ir.ImageRepository
	S3Repository       *sr.S3Repository
	RabbitMQRepository *mq.RabbitMQRepository

	// Services
	ImageService *s.ImageService

	// Handlers
	ImageHandlers *h.ImageHandlers
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
	ImageRepository = ir.NewImageRepository(DatabaseClient)
	S3Repository = sr.NewS3Repository(S3Client, Logger, Configuration.S3.BucketName)
	RabbitMQRepository = mq.NewRabbitMQRepository(RabbitMQClient)

	// Init services
	ImageService = s.NewImageService(ImageRepository, S3Repository, Logger)

	// Init handlers
	ImageHandlers = h.NewImageHandlers(ImageService, Logger)
}
