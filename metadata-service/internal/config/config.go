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
	CategoryRepository *r.CategoryRepository
	TagRepository      *r.TagRepository

	// Services
	CategoryService *s.CategoryService
	TagService      *s.TagService

	// Handlers
	CategoryHandlers *h.CategoryHandlers
	TagHandlers      *h.TagHandlers
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

	// Init services
	CategoryService = s.NewCategoryService(CategoryRepository)
	TagService = s.NewTagService(TagRepository)

	// Init handlers
	CategoryHandlers = h.NewCategoryHandlers(CategoryService, Logger)
	TagHandlers = h.NewTagHandlers(TagService, Logger)
}
