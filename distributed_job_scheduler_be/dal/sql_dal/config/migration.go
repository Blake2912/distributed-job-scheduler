package config

import (
	"log"

	datamodels "example.com/sql_dal/data_models"
	"gorm.io/gorm"
)

// Put in all of the data models here for migrating the same, maintain ordering
var allModels = []any{
	&datamodels.Jobs{},
}

func Migrate(database *gorm.DB) error {
	err := database.AutoMigrate(allModels...)
	if err != nil {
		log.Fatal("Migration failed:", err)
		return err
	}

	return nil
}
