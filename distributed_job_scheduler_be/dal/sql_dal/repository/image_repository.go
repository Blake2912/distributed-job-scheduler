package repository

import (
	"context"

	databaseconstants "example.com/constants/database_constants"
	"example.com/dal/sql_dal/models"
)

type ImageRepository interface {
	FindByTypeAndVersion(ctx context.Context, imageType databaseconstants.ImageType, imageVersion string) (models.ImageInformation, error)
}
