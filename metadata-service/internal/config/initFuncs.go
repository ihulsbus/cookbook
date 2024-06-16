package config

import (
	"fmt"
	"metadata-service/internal/helpers"
	m "metadata-service/internal/models"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initLogging() {
	log.Info("setting up the logging framework")

	Logger.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	logLevels := helpers.SetupLogLevels()

	if i, found := logLevels[strings.ToUpper(Configuration.Global.LogLevel)]; found {
		Logger.SetLevel(i)
		Logger.Infof("loglevel set to %s", strings.ToUpper(Logger.Level.String()))

	} else {
		Logger.Warn("no or invalid loglevel specified. Assuming default value of INFO. \n valid loglevels are: PANIC FATAL ERROR WARN INFO DEBUG TRACE")
		Logger.SetLevel(logLevels["INFO"])
	}
}

func initViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/config")
}

func initConfig() {
	Logger.Info("loading config")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Logger.Fatalf("config file not found: %v", err)
		} else {
			Logger.Fatalf("unknown error occured while reading config. error: %v", err)
		}
	}

	if err := viper.Unmarshal(&Configuration); err != nil {
		Logger.Fatalf("error unmarshaling config: %v", err)
	}

	Logger.Info("config file loaded")
}

func initDatabase() {
	Logger.Info("connecting to the database")
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
		Logger.Fatalf("Unable to connect to the database. Exiting..\n%v\n", err)
	}

	Logger.Info("performing database migrations")
	if err := DatabaseClient.AutoMigrate(
		&m.Tag{},
		&m.Category{},
		&m.CuisineType{},
		&m.PreparationTime{},
		&m.DifficultyLevel{},
		&m.RecipeTag{},
		&m.RecipeCategory{},
		&m.RecipeCuisineType{},
		&m.RecipePreparationTime{},
		&m.RecipeDifficultyLevel{},
	); err != nil {
		Logger.Fatalf("Error while automigrating database: %s", err.Error())
	}

	Logger.Info("connected!")
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
