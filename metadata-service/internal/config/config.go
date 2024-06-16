package config

import (
	h "metadata-service/internal/handlers"
	m "metadata-service/internal/models"
	r "metadata-service/internal/repositories"
	s "metadata-service/internal/services"

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
	CategoryRepository        *r.CategoryRepository
	TagRepository             *r.TagRepository
	CuisineTypeRepository     *r.CuisineTypeRepository
	DifficultyLevelRepository *r.DifficultyLevelRepository
	PreparationTimeRepository *r.PreparationTimeRepository

	// Services
	CategoryService        *s.CategoryService
	TagService             *s.TagService
	CuisineTypeService     *s.CuisineTypeService
	DifficultyLevelService *s.DifficultyLevelService
	PreparationTimeService *s.PreparationTimeService

	// Handlers
	CategoryHandlers        *h.CategoryHandlers
	TagHandlers             *h.TagHandlers
	CuisineTypeHandlers     *h.CuisineTypeHandlers
	DifficultyLevelHandlers *h.DifficultyLevelHandlers
	PreparationTimeHandlers *h.PreparationTimeHandlers
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
	CategoryRepository = r.NewCategoryRepository(DatabaseClient)
	TagRepository = r.NewTagRepository(DatabaseClient)
	CuisineTypeRepository = r.NewCuisineTypeRepository(DatabaseClient)
	DifficultyLevelRepository = r.NewDifficultyLevelRepository(DatabaseClient)
	PreparationTimeRepository = r.NewPreparationTimeRepository(DatabaseClient)

	// Init services
	CategoryService = s.NewCategoryService(CategoryRepository)
	TagService = s.NewTagService(TagRepository)
	CuisineTypeService = s.NewCuisineTypeService(CuisineTypeRepository)
	DifficultyLevelService = s.NewDifficultyLevelService(DifficultyLevelRepository)
	PreparationTimeService = s.NewPreparationTimeService(PreparationTimeRepository)

	// Init handlers
	CategoryHandlers = h.NewCategoryHandlers(CategoryService, Logger)
	TagHandlers = h.NewTagHandlers(TagService, Logger)
	CuisineTypeHandlers = h.NewCuisineTypeHandlers(CuisineTypeService, Logger)
	DifficultyLevelHandlers = h.NewDifficultyLevelHandlers(DifficultyLevelService, Logger)
	PreparationTimeHandlers = h.NewPreparationTimeHandlers(PreparationTimeService, Logger)
}
