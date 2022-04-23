package config

import (
	"github.com/ihulsbus/cookbook/internal/repositories"
	"github.com/ihulsbus/cookbook/internal/services"
	u "github.com/ihulsbus/cookbook/internal/utils"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

var (
	Configuration        Config
	RecipeRepository     *repositories.RecipeRepository
	IngredientRepository *repositories.IngredientRepository
	RecipeService        *services.RecipeService
	IngredientService    *services.IngredientService
)

func initViper() error {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("/etc/cookbook/")
	viper.AddConfigPath(".")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Fatalf("Config file not found! Error was: %v", err)
		} else {
			log.Fatalf("Unknown config error occured. Error is: %v", err)
		}
	}

	err := viper.Unmarshal(&Configuration)
	if err != nil {
		log.Fatalf("Error unmarshalling config file: %v", err)
	}

	viper.WatchConfig()

	log.Infof("Using config file found at: %s", viper.GetViper().ConfigFileUsed())
	if Configuration.Global.Debug {
		log.SetLevel(log.DebugLevel)
		log.Debugln("Enabled DEBUG logging level")
	}

	return err
}

func init() {
	if err := initViper(); err != nil {
		log.Fatalf("Error reading config: %v", err.Error())
	}

	// Init image folder
	err := u.InitFolder(Configuration.Global.ImageFolder)
	if err != nil {
		log.Fatal("Unable to create or detect image folder: %v", err)
	}

	// Init Database
	Configuration.DatabaseClient = initDatabase(Configuration.Database.Host, Configuration.Database.Username, Configuration.Database.Password, Configuration.Database.Database, Configuration.Database.Port, Configuration.Database.SSLMode, Configuration.Database.Timezone)

	// Init repositories
	RecipeRepository = repositories.NewRecipeRepository(Configuration.DatabaseClient)
	IngredientRepository = repositories.NewIngredientRepository(Configuration.DatabaseClient)

	// Init services
	RecipeService = services.NewRecipeService(RecipeRepository, Configuration.Global.ImageFolder)
	IngredientService = services.NewIngredientService(IngredientRepository)

}
