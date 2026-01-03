package migration

import (
	"context"
	"fmt"
	"log"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
	"gorm.io/gorm"
)

// Put in all of the data models here for migrating the same, maintain ordering
var allModels = []any{
	&models.Jobs{},
	&models.ImageInformation{},
	&models.JobExecution{},
}

func Migrate(ctx context.Context, database *gorm.DB) error {
	if err := database.WithContext(ctx).AutoMigrate(allModels...); err != nil {
		return fmt.Errorf("Migration failed: %w", err)
	}

	log.Println("Database migrations have been applied successfully")
	return nil
}
