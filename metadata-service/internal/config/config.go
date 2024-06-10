package config

import (
	"fmt"
	h "metadata-service/internal/helpers"
	m "metadata-service/internal/models"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	INIT_NOK = "> Init NOK"
	INIT_OK  = "> Init OK"
)

var (
	Configuration m.Config

	Logger         *log.Logger
	DatabaseClient *gorm.DB
	S3Client       *s3.S3
	Cors           cors.Config
)

func initLogging() {
	log.Info("> init logging")

	Logger = log.New()

	Logger.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	logLevels := h.SetupLogLevels()

	if i, found := logLevels[strings.ToUpper(Configuration.Global.LogLevel)]; found {
		Logger.SetLevel(i)

	} else {
		Logger.Warn("no or invalid loglevel specified. Assuming default value of INFO. \n valid loglevels are: PANIC FATAL ERROR WARN INFO DEBUG TRACE")
		Logger.SetLevel(logLevels["INFO"])
	}
}

func initViper() {
	log.Info("> init config")

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("config file not found: %v", err)
		} else {
			log.Fatalf("unknown error occured while reading config. error: %v", err)
		}
	}
	if err := viper.Unmarshal(&Configuration); err != nil {
		log.Fatalf("error unmarshaling config: %v", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Infof("config file changed: %s", e.Name)
	})
}

func initDatabase() {
	var err error

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		Configuration.Database.Host,
		Configuration.Database.Username,
		Configuration.Database.Password,
		Configuration.Database.Database,
		Configuration.Database.Port,
		Configuration.Database.SSLMode,
		Configuration.Database.Timezone)

	DatabaseClient, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		Logger.Errorf("Unable to connect to the database. Exiting..\n%v\n", err)
		Logger.Fatal(INIT_NOK)
	}

	if err := DatabaseClient.AutoMigrate(
		&m.Category{},
		&m.Tag{},
	); err != nil {
		Logger.Errorf("Error while automigrating database: %s", err.Error())
		Logger.Fatal(INIT_NOK)
	}

	Logger.Info(INIT_OK)
}

func initCors() {
	Cors = cors.Config{
		AllowOrigins:     Configuration.Cors.AllowedOrigins,
		AllowMethods:     Configuration.Cors.AllowedMethods,
		AllowHeaders:     Configuration.Cors.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: Configuration.Cors.AllowCredentials,
		MaxAge:           12 * time.Hour,
	}
}

func init() {
	initViper()
	initLogging()
	initDatabase()
	initCors()
}
