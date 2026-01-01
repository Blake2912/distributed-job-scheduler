package imageservice

import (
	"context"

	databaseconstants "example.com/constants/database_constants"
	"example.com/dal/sql_dal/models"
)

type Service interface {
	GetByTypeAndVersion(
		ctx context.Context,
		imageType databaseconstants.ImageType,
		version string,
	) (models.ImageInformation, error)
}
