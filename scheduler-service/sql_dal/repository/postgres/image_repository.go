package postgres

import (
	"context"

	databaseconstants "github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/config"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
	"gorm.io/gorm"
)

type ImageRepository struct {
	db *gorm.DB
}

func NewImageRepository(db *gorm.DB) repository.ImageRepository {
	return &ImageRepository{db: db}
}

func (*ImageRepository) FindByTypeAndVersion(ctx context.Context, imageType databaseconstants.ImageType, imageVersion string) (models.ImageInformation, error) {
	return gorm.G[models.ImageInformation](config.DB).
		Where(&models.ImageInformation{
			ImageType:    imageType,
			ImageVersion: imageVersion,
		}).
		First(ctx)
}
