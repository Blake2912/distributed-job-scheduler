package config

import (
	"context"
	"fmt"
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

func Migrate(ctx context.Context, database *gorm.DB) error {
	if err := database.WithContext(ctx).AutoMigrate(allModels...); err != nil {
		return fmt.Errorf("Migration failed: %w", err)
	}

	log.Println("Database migrations have been applied successfully")
	return nil
}
