package config

import (
	ch "metadata-service/internal/handlers/category"
	cuh "metadata-service/internal/handlers/cuisinetype"
	dh "metadata-service/internal/handlers/difficultylevel"
	sh "metadata-service/internal/handlers/search"
	th "metadata-service/internal/handlers/tag"

	m "metadata-service/internal/models"

	cr "metadata-service/internal/repositories/category"
	cur "metadata-service/internal/repositories/cuisinetype"
	dr "metadata-service/internal/repositories/difficultylevel"
	mr "metadata-service/internal/repositories/metadata"
	sr "metadata-service/internal/repositories/search"
	tr "metadata-service/internal/repositories/tag"

	cs "metadata-service/internal/services/category"
	cus "metadata-service/internal/services/cuisinetype"
	ds "metadata-service/internal/services/difficultylevel"
	ms "metadata-service/internal/services/metadata"
	ss "metadata-service/internal/services/search"
	ts "metadata-service/internal/services/tag"

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
	Cors           cors.Config
	KeycloakModule *keycloak.KeycloakModule

	// Repositories
	CategoryRepository        *cr.CategoryRepository
	CuisineTypeRepository     *cur.CuisineTypeRepository
	DifficultyLevelRepository *dr.DifficultyLevelRepository
	SearchRepository          *sr.SearchRepository
	TagRepository             *tr.TagRepository
	MetadataRepository        *mr.RecipeMetadataRepository

	// Services
	CategoryService        *cs.CategoryService
	CuisineTypeService     *cus.CuisineTypeService
	DifficultyLevelService *ds.DifficultyLevelService
	SearchService          *ss.SearchService
	TagService             *ts.TagService
	MetadataService        *ms.MetadataService

	// Handlers
	CategoryHandlers        *ch.CategoryHandlers
	CuisineTypeHandlers     *cuh.CuisineTypeHandlers
	DifficultyLevelHandlers *dh.DifficultyLevelHandlers
	SearchHandlers          *sh.SearchHandlers
	TagHandlers             *th.TagHandlers
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

	initDatabase()
	initCategories()
	initCuisineTypes()
	initDifficultyLevels()
	initCors()

	KeycloakModule, err = initOauth()
	if err != nil {
		Logger.Panicf("error initialising oauth: %v", err)
	}

	// Init repositories
	CategoryRepository = cr.NewCategoryRepository(DatabaseClient)
	CuisineTypeRepository = cur.NewCuisineTypeRepository(DatabaseClient)
	DifficultyLevelRepository = dr.NewDifficultyLevelRepository(DatabaseClient)
	SearchRepository = sr.NewSearchRepository(DatabaseClient)
	TagRepository = tr.NewTagRepository(DatabaseClient)
	MetadataRepository = mr.NewRecipeMetadataRepository(DatabaseClient)

	// Init services
	CategoryService = cs.NewCategoryService(CategoryRepository)
	CuisineTypeService = cus.NewCuisineTypeService(CuisineTypeRepository)
	DifficultyLevelService = ds.NewDifficultyLevelService(DifficultyLevelRepository)
	SearchService = ss.NewSearchService(SearchRepository)
	TagService = ts.NewTagService(TagRepository)
	MetadataService = ms.NewMetadataService(MetadataRepository)

	// Init handlers
	CategoryHandlers = ch.NewCategoryHandlers(CategoryService, Logger)
	CuisineTypeHandlers = cuh.NewCuisineTypeHandlers(CuisineTypeService, Logger)
	DifficultyLevelHandlers = dh.NewDifficultyLevelHandlers(DifficultyLevelService, Logger)
	SearchHandlers = sh.NewSearchHandlers(SearchService, Logger)
	TagHandlers = th.NewTagHandlers(TagService, Logger)
}
