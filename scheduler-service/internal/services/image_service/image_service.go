package imageservice

import (
	"context"

	databaseconstants "github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
)

type ImageService interface {
	GetByTypeAndVersion(
		ctx context.Context,
		imageType databaseconstants.ImageType,
		version string,
	) (models.ImageInformation, error)
}
