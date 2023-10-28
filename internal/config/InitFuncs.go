package config

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	m "github.com/ihulsbus/cookbook/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
)

var (
	units = []m.Unit{
		// US units
		{ID: 1, FullName: "Teaspoon", ShortName: "tsp"},
		{ID: 2, FullName: "Tablespoon", ShortName: "tbsp"},
		{ID: 3, FullName: "Fluid Ounce", ShortName: "fl oz"},
		{ID: 4, FullName: "Ounce", ShortName: "oz"},
		{ID: 5, FullName: "Pound", ShortName: "lb"},
		{ID: 6, FullName: "Cup", ShortName: "c"},
		{ID: 7, FullName: "Pint", ShortName: "pt"},
		{ID: 8, FullName: "Quart", ShortName: "qt"},
		{ID: 9, FullName: "Gallon", ShortName: "gal"},
		// Metric units
		{ID: 10, FullName: "Milliliter", ShortName: "ml"},
		{ID: 11, FullName: "Deciliter", ShortName: "dl"},
		{ID: 12, FullName: "Liter", ShortName: "l"},
		{ID: 13, FullName: "Milligram", ShortName: "mg"},
		{ID: 14, FullName: "Gram", ShortName: "g"},
		{ID: 15, FullName: "Kilogram", ShortName: "kg"},
		// Generic units
		{ID: 16, FullName: "Pinch", ShortName: "pn"},
		{ID: 17, FullName: "Cloves", ShortName: "cloves"},
		{ID: 18, FullName: "Pieces", ShortName: "pcs"},
	}
)

func connectS3(endpoint string, secret, key, region string) *s3.S3 {
	sess, err := session.NewSession(
		&aws.Config{
			Credentials:      credentials.NewStaticCredentials(key, secret, ""),
			Endpoint:         aws.String(endpoint),
			Region:           aws.String(region),
			S3ForcePathStyle: aws.Bool(false),
		})
	if err != nil {
		panic(err)
	}

	client := s3.New(sess)
	return client
}

func initDatabase(host string, user string, password string, dbname string, port int, sslmode string, timezone string) *gorm.DB {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s", host, user, password, dbname, port, sslmode, timezone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	})

	if err != nil {
		Logger.Errorf("Unable to connect to the database. Exiting..\n%v\n", err)
		Logger.Fatal(INIT_NOK)
	}

	err = db.SetupJoinTable(&m.Recipe{}, "Ingredient", &m.RecipeIngredient{})
	if err != nil {
		Logger.Errorf("Error while creating RecipeIngredient join tables: %s", err.Error())
		Logger.Fatal(INIT_NOK)
	}

	err = db.AutoMigrate(
		&m.Recipe{},
		&m.Ingredient{},
		&m.Instruction{},
		&m.Category{},
		&m.Tag{},
		&m.Unit{},
		&m.Author{},
	)

	if err != nil {
		Logger.Errorf("Error while automigrating database: %s", err.Error())
		Logger.Fatal(INIT_NOK)
	}

	return db
}

func initUnits() error {
	if err := Configuration.DatabaseClient.Transaction(func(tx *gorm.DB) error {
		if err := tx.Clauses(clause.OnConflict{
			DoNothing: true,
		}).Create(units).Error; err != nil {
			return err
		}

		return nil

	}); err != nil {
		return err
	}

	return nil
}
