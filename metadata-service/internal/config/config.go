package config

import (
	"context"
	"time"

	ch "metadata-service/internal/handlers/category"
	cuh "metadata-service/internal/handlers/cuisinetype"
	dh "metadata-service/internal/handlers/difficultylevel"
	mh "metadata-service/internal/handlers/metadata"
	ph "metadata-service/internal/handlers/preparationtimes"
	sh "metadata-service/internal/handlers/search"
	th "metadata-service/internal/handlers/tag"

	az "github.com/ihulsbus/cookbook/shared/authz"
	healthh "github.com/ihulsbus/cookbook/shared/healthchecks"
	hc "github.com/ihulsbus/cookbook/shared/httpclient"
	rmq "github.com/ihulsbus/cookbook/shared/rabbitmq"

	cr "metadata-service/internal/repositories/category"
	cur "metadata-service/internal/repositories/cuisinetype"
	dr "metadata-service/internal/repositories/difficultylevel"
	mr "metadata-service/internal/repositories/metadata"
	pr "metadata-service/internal/repositories/preparationtime"
	sr "metadata-service/internal/repositories/search"
	tr "metadata-service/internal/repositories/tag"

	ps "metadata-service/internal/services/PreparationTimes"
	cs "metadata-service/internal/services/category"
	cus "metadata-service/internal/services/cuisinetype"
	ds "metadata-service/internal/services/difficultylevel"
	ms "metadata-service/internal/services/metadata"
	ss "metadata-service/internal/services/search"
	ts "metadata-service/internal/services/tag"

	m "github.com/ihulsbus/cookbook/shared/models"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	"github.com/ihulsbus/cookbook/shared/oidc"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type metadataConfig struct {
	m.Config `mapstructure:",squash"`
	// RecipeServiceInternal is recipe-service's internal-only API (KIT-28),
	// used only by Authorizer's projection-miss fallback.
	RecipeServiceInternal m.ApiClient
}

var (
	Configuration metadataConfig
	err           error
	Ctx           context.Context

	Logger     = log.New()
	Cors       cors.Config
	OIDCModule *oidc.Module

	// Clients
	DatabaseClient *gorm.DB
	RabbitMQClient *rmq.RabbitMQ

	// Repositories
	CategoryRepository        *cr.CategoryRepository
	CuisineTypeRepository     *cur.CuisineTypeRepository
	DifficultyLevelRepository *dr.DifficultyLevelRepository
	SearchRepository          *sr.SearchRepository
	TagRepository             *tr.TagRepository
	MetadataRepository        *mr.RecipeMetadataRepository
	PreparationRepository     *pr.PreparationTimeRepository

	// Services
	CategoryService        *cs.CategoryService
	CuisineTypeService     *cus.CuisineTypeService
	DifficultyLevelService *ds.DifficultyLevelService
	SearchService          *ss.SearchService
	TagService             *ts.TagService
	MetadataService        *ms.MetadataService
	PreparationTimeService *ps.PreparationTime

	// Recipe access: see ingredient-service's config.go for the same
	// pattern.
	Authorizer *az.Authorizer
	Projector  *az.Projector

	// Handlers
	CategoryHandlers        *ch.CategoryHandlers
	CuisineTypeHandlers     *cuh.CuisineTypeHandlers
	DifficultyLevelHandlers *dh.DifficultyLevelHandlers
	SearchHandlers          *sh.SearchHandlers
	TagHandlers             *th.TagHandlers
	MetadataHandlers        *mh.MetadataHandlers
	HealthHandler           *healthh.Handlers
	PreparationTimeHandler  *ph.PreparationTime
)

func init() {
	Ctx = context.Background()

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

	initDatabase()
	initCategories()
	initCuisineTypes()
	initDifficultyLevels()
	initCors()

	OIDCModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}

	// RabbitMQClient.Connection is used by Projector's consumer (started
	// in MetadataService.go) - no synchronous call to any other service
	// happens during startup (KIT-79).
	RabbitMQClient, err = rmq.NewRabbitMQConnection(
		Configuration.RabbitMQ.Username,
		Configuration.RabbitMQ.Password,
		Configuration.RabbitMQ.Host,
		Logger,
	)

	// Init repositories
	CategoryRepository = cr.NewCategoryRepository(DatabaseClient)
	CuisineTypeRepository = cur.NewCuisineTypeRepository(DatabaseClient)
	DifficultyLevelRepository = dr.NewDifficultyLevelRepository(DatabaseClient)
	SearchRepository = sr.NewSearchRepository(DatabaseClient)
	TagRepository = tr.NewTagRepository(DatabaseClient)
	MetadataRepository = mr.NewRecipeMetadataRepository(DatabaseClient)
	PreparationRepository = pr.NewPreparationTimeRepository(DatabaseClient)

	// Init services
	CategoryService = cs.NewCategoryService(CategoryRepository)
	CuisineTypeService = cus.NewCuisineTypeService(CuisineTypeRepository)
	DifficultyLevelService = ds.NewDifficultyLevelService(DifficultyLevelRepository)
	SearchService = ss.NewSearchService(SearchRepository)
	TagService = ts.NewTagService(TagRepository)
	MetadataService = ms.NewMetadataService(MetadataRepository)
	PreparationTimeService = ps.NewPreparationTimeService(PreparationRepository)

	Projector = az.NewProjector(DatabaseClient)
	Authorizer = az.NewAuthorizer(DatabaseClient, Logger)
	fallbackClient := hc.NewInternalHTTPClient(5*time.Second, Logger, nil, nil)
	Authorizer.SetMissFallback(az.NewRecipeServiceFallback(fallbackClient, Configuration.RecipeServiceInternal.BaseURL, Projector))

	// Init handlers
	CategoryHandlers = ch.NewCategoryHandlers(CategoryService, Logger)
	CuisineTypeHandlers = cuh.NewCuisineTypeHandlers(CuisineTypeService, Logger)
	DifficultyLevelHandlers = dh.NewDifficultyLevelHandlers(DifficultyLevelService, Logger)
	SearchHandlers = sh.NewSearchHandlers(SearchService, Logger)
	TagHandlers = th.NewTagHandlers(TagService, Logger)
	MetadataHandlers = mh.NewMetadataHandlers(MetadataService, Logger)
	PreparationTimeHandler = ph.NewPreparationTimeHandlers(PreparationTimeService, Logger)
	HealthHandler = healthh.NewHealthHandlers(DatabaseClient, Logger)
}
