package config

import (
	"search-service/internal/helpers"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
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
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
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
