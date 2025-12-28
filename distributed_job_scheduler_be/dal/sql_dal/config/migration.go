package config

import (
	"log"

	datamodels "example.com/dal/sql_dal/datamodels"
	"gorm.io/gorm"
)

// Put in all of the data models here for migrating the same, maintain ordering
var allModels = []any{
	&datamodels.Jobs{},
	&datamodels.ImageInformation{},
	&datamodels.JobExecution{},
}

func Migrate(database *gorm.DB) error {
	err := database.AutoMigrate(allModels...)
	if err != nil {
		log.Fatal("Migration failed:", err)
		return err
	}

	return nil
}
