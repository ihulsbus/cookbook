package config

import (
	ah "ingredient-service/internal/handlers/amounts"
	ih "ingredient-service/internal/handlers/ingredients"
	rh "ingredient-service/internal/handlers/rabbitmq"
	uh "ingredient-service/internal/handlers/units"
	ar "ingredient-service/internal/repositories/amounts"
	ir "ingredient-service/internal/repositories/ingredients"
	rr "ingredient-service/internal/repositories/rabbitmq"
	ur "ingredient-service/internal/repositories/units"
	as "ingredient-service/internal/services/amounts"
	is "ingredient-service/internal/services/ingredients"
	us "ingredient-service/internal/services/units"

	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"
	m "github.com/ihulsbus/cookbook/shared/models"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/keycloak"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	Configuration m.Config
	err           error

	Logger         *log.Logger = log.New()
	DatabaseClient *gorm.DB
	KeycloakModule *keycloak.KeycloakModule
	Cors           cors.Config
	RabbitMQClient *rmq.RabbitMQ

	// Repositories
	AmountRepository     *ar.AmountRepository
	IngredientRepository *ir.IngredientRepository
	UnitRepository       *ur.UnitRepository
	RabbitMQRepository   *rr.Repository

	// Services
	AmountService     *as.AmountService
	IngredientService *is.IngredientService
	UnitService       *us.UnitService

	// Handlers
	AmountHandlers     *ah.AmountHandlers
	IngredientHandlers *ih.IngredientHandlers
	UnitHandlers       *uh.UnitHandlers
	HealthHandler      *healthh.Handlers
	RabbitMQHandler    *rh.RabbitMQHandler
)

func init() {
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

	if Configuration.Global.ListenPort == "" {
		Logger.Warn("Listen port is empty. Defaulting to 8080")
		Configuration.Global.ListenPort = "8080"
	}

	initDatabase()
	initUnits()
	initCors()

	// Init repositories
	AmountRepository = ar.NewAmountRepository(DatabaseClient)
	IngredientRepository = ir.NewIngredientRepository(DatabaseClient)
	UnitRepository = ur.NewUnitRepository(DatabaseClient)

	// Init services
	AmountService = as.NewAmountService(AmountRepository)
	IngredientService = is.NewIngredientService(IngredientRepository)
	UnitService = us.NewUnitService(UnitRepository)

	// Init handlers
	AmountHandlers = ah.NewAmountHandlers(AmountService, Logger)
	IngredientHandlers = ih.NewIngredientHandlers(IngredientService, Logger)
	UnitHandlers = uh.NewUnitHandlers(UnitService, Logger)
	HealthHandler = healthh.NewHealthHandlers(DatabaseClient, Logger)
}
