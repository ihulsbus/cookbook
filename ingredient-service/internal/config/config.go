package config

import (
	"context"
	ah "ingredient-service/internal/handlers/amounts"
	ih "ingredient-service/internal/handlers/ingredients"
	uh "ingredient-service/internal/handlers/units"
	ar "ingredient-service/internal/repositories/amounts"
	ir "ingredient-service/internal/repositories/ingredients"
	ur "ingredient-service/internal/repositories/units"
	as "ingredient-service/internal/services/AmountService"
	is "ingredient-service/internal/services/IngredientService"
	us "ingredient-service/internal/services/UnitService"
	"time"

	az "github.com/ihulsbus/cookbook/shared/authz"
	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"
	hc "github.com/ihulsbus/cookbook/shared/httpclient"
	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/oidc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type ingredientConfig struct {
	m.Config `mapstructure:",squash"`
	// RecipeServiceInternal is recipe-service's internal-only API (never
	// through the ingress), used only by Authorizer's projection-miss fallback.
	RecipeServiceInternal m.ApiClient
}

var (
	Configuration ingredientConfig
	Ctx           context.Context
	err           error

	Logger     *log.Logger = log.New()
	OIDCModule *oidc.Module
	Cors       cors.Config

	DatabaseClient *gorm.DB
	RabbitMQClient *rmq.RabbitMQ

	AmountRepository     *ar.AmountRepository
	IngredientRepository *ir.IngredientRepository
	UnitRepository       *ur.UnitRepository

	AmountService     *as.AmountService
	IngredientService *is.IngredientService
	UnitService       *us.UnitService

	// Authorizer answers CanRead/CanWrite against a local projection kept
	// current by Projector's RabbitMQ consumer, falling back to
	// recipe-service's internal endpoint on a projection miss.
	Authorizer *az.Authorizer
	Projector  *az.Projector

	AmountHandlers     *ah.AmountHandlers
	IngredientHandlers *ih.IngredientHandlers
	UnitHandlers       *uh.UnitHandlers
	HealthHandler      *healthh.Handlers
)

func init() {
	Ctx = context.Background()

	initViper()
	initConfig()
	initLogging()
	OIDCModule, err = initOauth()
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

	// No synchronous call to another service happens at startup: recipe
	// access state arrives via events (Projector's consumer), with a
	// lazy, circuit-broken fallback only made at request time on a miss.
	RabbitMQClient, err = rmq.NewRabbitMQConnection(
		Configuration.RabbitMQ.Username,
		Configuration.RabbitMQ.Password,
		Configuration.RabbitMQ.Host,
		Logger,
	)

	AmountRepository = ar.NewAmountRepository(DatabaseClient)
	IngredientRepository = ir.NewIngredientRepository(DatabaseClient)
	UnitRepository = ur.NewUnitRepository(DatabaseClient)

	AmountService = as.NewAmountService(AmountRepository)
	IngredientService = is.NewIngredientService(IngredientRepository)
	UnitService = us.NewUnitService(UnitRepository)

	Projector = az.NewProjector(DatabaseClient)
	Authorizer = az.NewAuthorizer(DatabaseClient, Logger)
	fallbackClient := hc.NewInternalHTTPClient(5*time.Second, Logger, nil, nil)
	Authorizer.SetMissFallback(az.NewRecipeServiceFallback(fallbackClient, Configuration.RecipeServiceInternal.BaseURL, Projector))

	AmountHandlers = ah.NewAmountHandlers(AmountService, Logger)
	IngredientHandlers = ih.NewIngredientHandlers(IngredientService, Logger)
	UnitHandlers = uh.NewUnitHandlers(UnitService, Logger)
	HealthHandler = healthh.NewHealthHandlers(DatabaseClient, Logger)
}
