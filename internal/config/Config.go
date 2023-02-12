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

const (
	InitNOK = "> Init NOK"
	InitOK  = "> Init OK"
)

var (
	// Generic
	envBinds []string = []string{
		"debug",
		"oidc_url",
		"oidc_clientid",
		"database_host",
		"database_port",
		"database_database",
		"database_username",
		"database_password",
		"database_sslmode",
		"database_timezone",
		"s3_endpoint",
		"s3_key",
		"s3_secret",
		"s3_bucket",
	}
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
	RecipeHandlers     *h.RecipeHandlers
	IngredientHandlers *h.IngredientHandlers

	// Endpoints
	RecipeEndpoints     *e.RecipeEndpoints
	IngredientEndpoints *e.IngredientEndpoints
	TagEndpoints        *e.TagEndpoints
	CategoryEndpoints   *e.CategoryEndpoints
)

func initViper() {
	viper.SetEnvPrefix("cbb")

	for i := range envBinds {
		err := viper.BindEnv(envBinds[i])
		if err != nil {
			Logger.Errorf("error binding to env var '%s': %s", envBinds[i], err.Error())
			Logger.Fatal(InitNOK)
		}
	}

	// global
	Configuration.Global.Debug = viper.GetBool("debug")

	// oidc
	Configuration.Oidc.URL = viper.GetString("oidc_url")
	Configuration.Oidc.ClientID = viper.GetString("oidc_clientid")
	Configuration.Oidc.SigningAlgs = append(Configuration.Oidc.SigningAlgs, "RS256")
	Configuration.Oidc.SkipClientIDCheck = true // static for now. figure out if configurability is needed in our case
	Configuration.Oidc.SkipExpiryCheck = true
	Configuration.Oidc.SkipIssuerCheck = true

	// database
	Configuration.Database.Host = viper.GetString("database_host")
	Configuration.Database.Port = viper.GetInt("database_port")
	Configuration.Database.Database = viper.GetString("database_database")
	Configuration.Database.Username = viper.GetString("database_username")
	Configuration.Database.Password = viper.GetString("database_password")
	Configuration.Database.SSLMode = viper.GetString("database_sslmode")
	Configuration.Database.Timezone = viper.GetString("database_timezone")

	// S3
	Configuration.S3.Endpoint = viper.GetString("s3_endpoint")
	Configuration.S3.AWSAccessKey = viper.GetString("s3_key")
	Configuration.S3.AWSAccessSecret = viper.GetString("s3_secret")
	Configuration.S3.BucketName = viper.GetString("s3_bucket")

}

func initCors() {
	Configuration.Cors.AllowedOrigins = append(Configuration.Cors.AllowedOrigins, "https://cookbook.hulsbus.be")
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

	Logger.Info("> init config")
	// Init Viper
	initViper()

	Logger.Info("> init cors")
	// Init CORS rules
	initCors()

	if Configuration.Global.Debug {
		Logger.SetLevel(log.DebugLevel)
		Logger.Debugln("> Enabled DEBUG logging level")
	}

	Logger.Info("> init DB")
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

	if err := initUnits(); err != nil {
		Logger.Error(err)
		Logger.Fatal(InitNOK)
	}

	Logger.Info("> init S3")
	// Init S3 session
	Configuration.S3ClientSession = connectS3(Configuration.S3.Endpoint, Configuration.S3.AWSAccessSecret, Configuration.S3.AWSAccessKey, "us-east-1")

	Logger.Info("> init application")
	// Init middleware
	Middleware = mi.NewMiddleware(&Configuration.Oidc, Logger)

	// Init repositories
	RecipeRepository = r.NewRecipeRepository(Configuration.DatabaseClient)
	IngredientRepository = r.NewIngredientRepository(Configuration.DatabaseClient)
	S3Repository = r.NewS3Repository(Configuration.DatabaseClient, Configuration.S3, Configuration.S3ClientSession, Logger)

	// Init services
	RecipeService = s.NewRecipeService(RecipeRepository)
	IngredientService = s.NewIngredientService(IngredientRepository)
	ImageService = s.NewImageService(S3Repository, Logger)

	// Init handlers
	RecipeHandlers = h.NewRecipeHandlers(RecipeService, ImageService, Logger)
	IngredientHandlers = h.NewIngredientHandlers(IngredientService, Logger)

	// Init endpoints
	RecipeEndpoints = e.NewRecipeEndpoints(RecipeHandlers)
	IngredientEndpoints = e.NewIngredientEndpoints(IngredientHandlers)
	TagEndpoints = e.NewTagEndpoints()
	CategoryEndpoints = e.NewCategoryEndpoints()

	Logger.Info(InitOK)
}
