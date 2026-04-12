package config

import (
	sh "search-service/internal/handlers"
	m "search-service/internal/models"
	ir "search-service/internal/repositories/ingredients"
	mr "search-service/internal/repositories/metadata"
	rr "search-service/internal/repositories/recipes"
	ss "search-service/internal/services"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"
	"github.com/ihulsbus/cookbook/shared/keycloak"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Configuration m.Config
	err           error

	Logger         *log.Logger = log.New()
	Cors           cors.Config
	KeycloakModule *keycloak.KeycloakModule

	// Repositories
	IngredientRepository *ir.IngredientRepository
	MetadataRepository   *mr.MetadataRepository
	RecipeRepository     *rr.RecipeRepository

	// Services
	SearchService *ss.SearchService

	// Handlers
	SearchHandlers *sh.SearchHandlers
	HealthHandler  *healthh.Handlers
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

	if Configuration.Global.ListenPort == "" {
		Logger.Warn("Listen port is empty. Defaulting to 8080")
		Configuration.Global.ListenPort = "8080"
	}

	initCors()
	KeycloakModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}

	// Init repositories
	IngredientRepository = ir.NewIngredientRepository(Configuration.IngredientService.BaseURL)
	MetadataRepository = mr.NewMetadataRepository(Configuration.MetadataService.BaseURL)
	RecipeRepository = rr.NewRecipeRepository(Configuration.RecipeService.BaseURL)

	// Init services
	SearchService = ss.NewSearchService(RecipeRepository, IngredientRepository, MetadataRepository)

	// Init handlers
	SearchHandlers = sh.NewSearchHandlers(SearchService)
	HealthHandler = healthh.NewHealthHandlers(nil, Logger)
}
