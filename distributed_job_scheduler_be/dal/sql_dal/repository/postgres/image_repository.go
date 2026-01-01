package postgres

import (
	"context"

	databaseconstants "example.com/constants/database_constants"
	"example.com/dal/sql_dal/config"
	"example.com/dal/sql_dal/models"
	"example.com/dal/sql_dal/repository"
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
