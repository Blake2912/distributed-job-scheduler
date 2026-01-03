package repository

import (
	"context"

	databaseconstants "github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
)

type ImageRepository interface {
	FindByTypeAndVersion(ctx context.Context, imageType databaseconstants.ImageType, imageVersion string) (models.ImageInformation, error)
}
