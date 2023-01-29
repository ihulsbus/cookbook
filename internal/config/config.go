package config

import (
	e "github.com/ihulsbus/cookbook/internal/endpoints"
	h "github.com/ihulsbus/cookbook/internal/handlers"
	mi "github.com/ihulsbus/cookbook/internal/middleware"
	m "github.com/ihulsbus/cookbook/internal/models"
	r "github.com/ihulsbus/cookbook/internal/repositories"
	s "github.com/ihulsbus/cookbook/internal/services"
	u "github.com/ihulsbus/cookbook/internal/utils"
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

	// Services
	RecipeService     *s.RecipeService
	IngredientService *s.IngredientService

	// Handlers
	Handlers *h.Handlers

	// Endpoints
	Endpoints *e.Endpoints
)

func initViper() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/cookbook/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			Logger.Fatalf("Config file not found! Error was: %v", err)
		} else {
			Logger.Fatalf("Unknown config error occured. Error is: %v", err)
		}
	}

	err := viper.Unmarshal(&Configuration)
	if err != nil {
		Logger.Fatalf("Error unmarshalling config file: %v", err)
	}

	viper.WatchConfig()

	Logger.Infof("Using config file found at: %s", viper.GetViper().ConfigFileUsed())
	if Configuration.Global.Debug {
		Logger.SetLevel(log.DebugLevel)
		Logger.Debugln("Enabled DEBUG logging level")
	}

	return err
}

func init() {
	// Init logger
	Logger = log.New()
	Logger.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})

	// Init Viper
	if err := initViper(); err != nil {
		Logger.Fatalf("Error reading config: %v", err.Error())
	}

	// Init image folder
	err := u.InitFolder(Configuration.Global.ImageFolder)
	if err != nil {
		Logger.Fatal("Unable to create or detect image folder: %v", err)
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

	// Init middleware
	Middleware = mi.NewMiddleware(&Configuration.Oidc, Logger)

	// Init repositories
	RecipeRepository = r.NewRecipeRepository(Configuration.DatabaseClient, Logger)
	IngredientRepository = r.NewIngredientRepository(Configuration.DatabaseClient, Logger)

	// Init services
	RecipeService = s.NewRecipeService(RecipeRepository, Configuration.Global.ImageFolder, Logger)
	IngredientService = s.NewIngredientService(IngredientRepository, Logger)

	// Init handlers
	Handlers = h.NewHandlers(RecipeService, IngredientService, Logger)

	// Init endpoints
	Endpoints = e.NewEndpoints(Handlers)

}
