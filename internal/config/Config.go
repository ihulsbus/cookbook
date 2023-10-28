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
	INIT_NOK = "> Init NOK"
	INIT_OK  = "> Init OK"
)

var (
	// Generic
	envBinds []string = []string{
		"debug",
		"auth0_domain",
		"auth0_clientid",
		"auth0_audience",
		"database_host",
		"database_port",
		"database_name",
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
	S3Repository         *r.S3Repository
	IngredientRepository *r.IngredientRepository
	TagRepository        *r.TagRepository
	CategoryRepository   *r.CategoryRepository

	// Services
	RecipeService     *s.RecipeService
	ImageService      *s.ImageService
	IngredientService *s.IngredientService
	TagService        *s.TagService
	CategoryService   *s.CategoryService

	// Handlers
	RecipeHandlers     *h.RecipeHandlers
	IngredientHandlers *h.IngredientHandlers
	TagHandlers        *h.TagHandlers
	CategoryHandlers   *h.CategoryHandlers

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
			Logger.Fatal(INIT_NOK)
		}
	}

	// global
	Configuration.Global.Debug = viper.GetBool("debug")

	// auth
	Configuration.Auth0.Domain = viper.GetString("auth0_domain")
	Configuration.Auth0.ClientID = viper.GetString("auth0_clientid")
	Configuration.Auth0.Audience = viper.GetString("auth0_audience")

	// database
	Configuration.Database.Host = viper.GetString("database_host")
	Configuration.Database.Port = viper.GetInt("database_port")
	Configuration.Database.Database = viper.GetString("database_name")
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
	Configuration.Cors.AllowedOrigins = append(Configuration.Cors.AllowedOrigins, "https://cookbook.hulsbus.be", "https://gourmedy.com", "http://localhost:4200")
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
		Logger.Fatal(INIT_NOK)
	}

	Logger.Info("> init S3")
	// Init S3 session
	Configuration.S3ClientSession = connectS3(Configuration.S3.Endpoint, Configuration.S3.AWSAccessSecret, Configuration.S3.AWSAccessKey, "us-east-1")

	Logger.Info("> init application")
	// Init middleware
	Middleware = mi.NewMiddleware(&Configuration.Auth0, Logger)

	// Init repositories
	RecipeRepository = r.NewRecipeRepository(Configuration.DatabaseClient)
	S3Repository = r.NewS3Repository(Configuration.DatabaseClient, Configuration.S3, Configuration.S3ClientSession, Logger)
	IngredientRepository = r.NewIngredientRepository(Configuration.DatabaseClient)
	TagRepository = r.NewTagRepository(Configuration.DatabaseClient)
	CategoryRepository = r.NewCategoryRepository(Configuration.DatabaseClient)

	// Init services
	RecipeService = s.NewRecipeService(RecipeRepository)
	ImageService = s.NewImageService(S3Repository, Logger)
	IngredientService = s.NewIngredientService(IngredientRepository)
	TagService = s.NewTagService(TagRepository)
	CategoryService = s.NewCategoryService(CategoryRepository)

	// Init handlers
	RecipeHandlers = h.NewRecipeHandlers(RecipeService, ImageService, Logger)
	IngredientHandlers = h.NewIngredientHandlers(IngredientService, Logger)
	TagHandlers = h.NewTagHandlers(TagService, Logger)
	CategoryHandlers = h.NewCategoryHandlers(CategoryService, Logger)

	// Init endpoints
	RecipeEndpoints = e.NewRecipeEndpoints(RecipeHandlers, &Middleware.AuthMW)
	IngredientEndpoints = e.NewIngredientEndpoints(IngredientHandlers)
	TagEndpoints = e.NewTagEndpoints(TagHandlers)
	CategoryEndpoints = e.NewCategoryEndpoints(CategoryHandlers)

	Logger.Info(INIT_OK)
}
