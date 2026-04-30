package config

import (
	"context"
	ah "ingredient-service/internal/handlers/amounts"
	ih "ingredient-service/internal/handlers/ingredients"
	rh "ingredient-service/internal/handlers/rabbitmq"
	uh "ingredient-service/internal/handlers/units"
	ar "ingredient-service/internal/repositories/amounts"
	ir "ingredient-service/internal/repositories/ingredients"
	ur "ingredient-service/internal/repositories/units"
	as "ingredient-service/internal/services/AmountService"
	cs "ingredient-service/internal/services/CacheService"
	is "ingredient-service/internal/services/IngredientService"
	us "ingredient-service/internal/services/UnitService"
	"time"

	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"
	hc "github.com/ihulsbus/cookbook/shared/httpclient"
	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"
	rc "github.com/ihulsbus/cookbook/shared/recipeclient"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/keycloak"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ingredientConfig struct {
	m.Config     `mapstructure:",squash"`
	RecipeClient m.ApiClient
}

var (
	Configuration ingredientConfig
	Ctx           context.Context
	err           error

	Logger         *log.Logger = log.New()
	KeycloakModule *keycloak.KeycloakModule
	Cors           cors.Config

	// Clients
	DatabaseClient *gorm.DB
	HttpClient     *hc.HTTPClient
	RabbitMQClient *rmq.RabbitMQ
	RecipeClient   *rc.RecipeAPIClient

	// Repositories
	AmountRepository     *ar.AmountRepository
	IngredientRepository *ir.IngredientRepository
	UnitRepository       *ur.UnitRepository

	// Services
	AmountService     *as.AmountService
	IngredientService *is.IngredientService
	UnitService       *us.UnitService
	CacheService      *cs.CacheService

	// Handlers
	AmountHandlers     *ah.AmountHandlers
	IngredientHandlers *ih.IngredientHandlers
	UnitHandlers       *uh.UnitHandlers
	HealthHandler      *healthh.Handlers
	RabbitMQHandler    *rh.RabbitMQHandler
)

func init() {
	Ctx = context.Background()

	initViper()
	initConfig()
	initLogging()
	KeycloakModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}

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

	initDatabase()
	initUnits()
	initCors()

	// Init clients
	HttpClient = hc.NewHTTPClient(5*time.Second, Configuration.Oauth.Url, Configuration.Oauth.Realm, Configuration.Oauth.ClientID, Configuration.Oauth.ClientSecret, Logger)
	RecipeClient, err = rc.NewRecipeAPIClient(Configuration.RecipeClient.BaseURL, HttpClient)
	if err != nil {
		Logger.Panicf("error initialising recipe client: %v", err)
	}
	RabbitMQClient, err = rmq.NewRabbitMQConnection(
		Configuration.RabbitMQ.Username,
		Configuration.RabbitMQ.Password,
		Configuration.RabbitMQ.Host,
		Logger,
	)

	// Init repositories
	AmountRepository = ar.NewAmountRepository(DatabaseClient)
	IngredientRepository = ir.NewIngredientRepository(DatabaseClient)
	UnitRepository = ur.NewUnitRepository(DatabaseClient)

	// Init services
	Logger.Debugf("Setting up cache service on port %d", Configuration.Global.ListenPort+1000)
	CacheService, err = cs.NewCacheService(Ctx, RecipeClient, Logger, Configuration.Global.ListenPort+1000)
	if err != nil {
		Logger.Fatalf("Error setting up cache: %v", err)
	}
	AmountService = as.NewAmountService(AmountRepository)
	IngredientService = is.NewIngredientService(IngredientRepository)
	UnitService = us.NewUnitService(UnitRepository)

	// Init handlers
	AmountHandlers = ah.NewAmountHandlers(AmountService, Logger)
	IngredientHandlers = ih.NewIngredientHandlers(IngredientService, Logger)
	UnitHandlers = uh.NewUnitHandlers(UnitService, Logger)
	HealthHandler = healthh.NewHealthHandlers(DatabaseClient, Logger)
	RabbitMQHandler, err = rh.NewRabbitMQHandler(CacheService, &Ctx, Logger)
	if err != nil {
		Logger.Fatalf("Error setting up RabbitMQ Consumer: %v", err)
	}
}
