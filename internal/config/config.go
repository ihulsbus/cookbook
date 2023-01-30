package config

import (
	e "github.com/ihulsbus/cookbook/internal/endpoints"
	h "github.com/ihulsbus/cookbook/internal/handlers"
	mi "github.com/ihulsbus/cookbook/internal/middleware"
	m "github.com/ihulsbus/cookbook/internal/models"
	r "github.com/ihulsbus/cookbook/internal/repositories"
	s "github.com/ihulsbus/cookbook/internal/services"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	// Generic
	Configuration m.Config
	Logger        *log.Logger

	// Middleware
	Middleware *mi.Middleware

	// Repositories
	RecipeRepository     *r.RecipeRepository
	IngredientRepository *r.IngredientRepository
	S3Repository         *r.S3Repository

	// Services
	RecipeService     *s.RecipeService
	IngredientService *s.IngredientService
	ImageService      *s.ImageService

	// Handlers
	Handlers *h.Handlers

	// Endpoints
	Endpoints *e.Endpoints
)

func initViper() {
	viper.SetEnvPrefix("cbb")

	// global
	viper.BindEnv("debug")

	Configuration.Global.Debug = viper.GetBool("debug")

	// oidc
	viper.BindEnv("oidc_url")
	viper.BindEnv("oidc_clientid")

	Configuration.Oidc.URL = viper.GetString("oidc_url")
	Configuration.Oidc.ClientID = viper.GetString("oidc_clientid")
	Configuration.Oidc.SigningAlgs = append(Configuration.Oidc.SigningAlgs, "RS256")
	Configuration.Oidc.SkipClientIDCheck = true // static for now. figure out if configurability is needed in our case
	Configuration.Oidc.SkipExpiryCheck = true
	Configuration.Oidc.SkipIssuerCheck = true

	// database
	viper.BindEnv("database_host")
	viper.BindEnv("database_port")
	viper.BindEnv("database_database")
	viper.BindEnv("database_username")
	viper.BindEnv("database_password")
	viper.BindEnv("database_sslmode")
	viper.BindEnv("database_timezone")

	Configuration.Database.Host = viper.GetString("database_host")
	Configuration.Database.Port = viper.GetInt("database_port")
	Configuration.Database.Database = viper.GetString("database_database")
	Configuration.Database.Username = viper.GetString("database_username")
	Configuration.Database.Password = viper.GetString("database_password")
	Configuration.Database.SSLMode = viper.GetString("database_sslmode")
	Configuration.Database.Timezone = viper.GetString("database_timezone")

	// S3
	viper.BindEnv("s3_endpoint")
	viper.BindEnv("s3_key")
	viper.BindEnv("s3_secret")
	viper.BindEnv("s3_bucket")

	Configuration.S3.Endpoint = viper.GetString("s3_endpoint")
	Configuration.S3.AWSAccessKey = viper.GetString("s3_key")
	Configuration.S3.AWSAccessSecret = viper.GetString("s3_secret")
	Configuration.S3.BucketName = viper.GetString("s3_bucket")

}

func initCors() {
	Configuration.Cors.AllowedOrigins = append(Configuration.Cors.AllowedOrigins, "*")
	Configuration.Cors.AllowCredentials = false
	Configuration.Cors.AllowedHeaders = append(Configuration.Cors.AllowedHeaders,
		"Authorization",
		"Content-Type",
	)
	Configuration.Cors.AllowedMethods = append(Configuration.Cors.AllowedMethods,
		"GET",
		"POST",
		"PUT",
		"DELETE",
		"HEAD",
	)
	Configuration.Cors.Debug = Configuration.Global.Debug

}

func init() {
	// Init logger
	Logger = log.New()
	Logger.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// Init Viper
	initViper()

	// Init CORS rules
	initCors()

	if Configuration.Global.Debug {
		Logger.SetLevel(log.DebugLevel)
		Logger.Debugln("Enabled DEBUG logging level")
	}

	// Init Database
	Configuration.DatabaseClient = initDatabase(
		Configuration.Database.Host,
		Configuration.Database.Username,
		Configuration.Database.Password,
		Configuration.Database.Database,
		Configuration.Database.Port,
		Configuration.Database.SSLMode,
		Configuration.Database.Timezone,
	)

	// Init S3 session
	Configuration.S3ClientSession = connectS3(Configuration.S3.Endpoint, Configuration.S3.AWSAccessSecret, Configuration.S3.AWSAccessKey, "us-east-1")

	// Init middleware
	Middleware = mi.NewMiddleware(&Configuration.Oidc, Logger)

	// Init repositories
	RecipeRepository = r.NewRecipeRepository(Configuration.DatabaseClient, Logger)
	IngredientRepository = r.NewIngredientRepository(Configuration.DatabaseClient, Logger)
	S3Repository = r.NewS3Repository(Configuration.DatabaseClient, Configuration.S3, Configuration.S3ClientSession, Logger)

	// Init services
	RecipeService = s.NewRecipeService(RecipeRepository, Configuration.Global.ImageFolder, Logger)
	IngredientService = s.NewIngredientService(IngredientRepository, Logger)
	ImageService = s.NewImageService(S3Repository, Logger)

	// Init handlers
	Handlers = h.NewHandlers(RecipeService, IngredientService, ImageService, Logger)

	// Init endpoints
	Endpoints = e.NewEndpoints(Handlers)

}
