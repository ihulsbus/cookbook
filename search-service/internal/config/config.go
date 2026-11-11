package config

import (
	sh "search-service/internal/handlers"
	ir "search-service/internal/repositories/ingredients"
	mr "search-service/internal/repositories/metadata"
	rr "search-service/internal/repositories/recipes"
	ss "search-service/internal/services"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"
	"github.com/ihulsbus/cookbook/shared/oidc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

type searchConfig struct {
	m.Config          `mapstructure:",squash"`
	IngredientService m.ApiClient
	MetadataService   m.ApiClient
	RecipeService     m.ApiClient
}

var (
	Configuration searchConfig
	err           error

	Logger     *log.Logger = log.New()
	Cors       cors.Config
	OIDCModule *oidc.Module

	IngredientRepository *ir.IngredientRepository
	MetadataRepository   *mr.MetadataRepository
	RecipeRepository     *rr.RecipeRepository

	SearchService *ss.SearchService

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

	if Configuration.Global.ListenPort == 0 {
		Logger.Warn("Listen port is empty. Defaulting to 8080")
		Configuration.Global.ListenPort = 8080
	}

	initCors()
	OIDCModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}

	IngredientRepository = ir.NewIngredientRepository(Configuration.IngredientService.BaseURL)
	MetadataRepository = mr.NewMetadataRepository(Configuration.MetadataService.BaseURL)
	RecipeRepository = rr.NewRecipeRepository(Configuration.RecipeService.BaseURL)

	SearchService = ss.NewSearchService(RecipeRepository, IngredientRepository, MetadataRepository)

	SearchHandlers = sh.NewSearchHandlers(SearchService)
	HealthHandler = healthh.NewHealthHandlers(nil, Logger)
}
