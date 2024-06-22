package config

import (
	ch "metadata-service/internal/handlers/category"
	cuh "metadata-service/internal/handlers/cuisinetype"
	dh "metadata-service/internal/handlers/difficultylevel"
	ph "metadata-service/internal/handlers/preparationtime"
	sh "metadata-service/internal/handlers/search"
	th "metadata-service/internal/handlers/tag"

	m "metadata-service/internal/models"

	cr "metadata-service/internal/repositories/category"
	cur "metadata-service/internal/repositories/cuisinetype"
	dr "metadata-service/internal/repositories/difficultylevel"
	pr "metadata-service/internal/repositories/preparationtime"
	sr "metadata-service/internal/repositories/search"
	tr "metadata-service/internal/repositories/tag"

	cs "metadata-service/internal/services/category"
	cus "metadata-service/internal/services/cuisinetype"
	ds "metadata-service/internal/services/difficultylevel"
	ps "metadata-service/internal/services/preparationtime"
	ss "metadata-service/internal/services/search"
	ts "metadata-service/internal/services/tag"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-contrib/cors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

var (
	Configuration m.Config

	Logger         *log.Logger = log.New()
	DatabaseClient *gorm.DB
	Cors           cors.Config

	// Repositories
	CategoryRepository        *cr.CategoryRepository
	CuisineTypeRepository     *cur.CuisineTypeRepository
	DifficultyLevelRepository *dr.DifficultyLevelRepository
	PreparationTimeRepository *pr.PreparationTimeRepository
	SearchRepository          *sr.SearcRepository
	TagRepository             *tr.TagRepository

	// Services
	CategoryService        *cs.CategoryService
	CuisineTypeService     *cus.CuisineTypeService
	DifficultyLevelService *ds.DifficultyLevelService
	PreparationTimeService *ps.PreparationTimeService
	SearchService          *ss.SearchService
	TagService             *ts.TagService

	// Handlers
	CategoryHandlers        *ch.CategoryHandlers
	CuisineTypeHandlers     *cuh.CuisineTypeHandlers
	DifficultyLevelHandlers *dh.DifficultyLevelHandlers
	PreparationTimeHandlers *ph.PreparationTimeHandlers
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
	initCors()

	// Init repositories
	CategoryRepository = cr.NewCategoryRepository(DatabaseClient)
	CuisineTypeRepository = cur.NewCuisineTypeRepository(DatabaseClient)
	DifficultyLevelRepository = dr.NewDifficultyLevelRepository(DatabaseClient)
	PreparationTimeRepository = pr.NewPreparationTimeRepository(DatabaseClient)
	SearchRepository = sr.NewSearchRepository(DatabaseClient)
	TagRepository = tr.NewTagRepository(DatabaseClient)

	// Init services
	CategoryService = cs.NewCategoryService(CategoryRepository)
	CuisineTypeService = cus.NewCuisineTypeService(CuisineTypeRepository)
	DifficultyLevelService = ds.NewDifficultyLevelService(DifficultyLevelRepository)
	PreparationTimeService = ps.NewPreparationTimeService(PreparationTimeRepository)
	SearchService = ss.NewSearchService(SearchRepository)
	TagService = ts.NewTagService(TagRepository)

	// Init handlers
	CategoryHandlers = ch.NewCategoryHandlers(CategoryService, Logger)
	CuisineTypeHandlers = cuh.NewCuisineTypeHandlers(CuisineTypeService, Logger)
	DifficultyLevelHandlers = dh.NewDifficultyLevelHandlers(DifficultyLevelService, Logger)
	PreparationTimeHandlers = ph.NewPreparationTimeHandlers(PreparationTimeService, Logger)
	SearchHandlers = sh.NewSearchHandlers(SearchService, Logger)
	TagHandlers = th.NewTagHandlers(TagService, Logger)
}
