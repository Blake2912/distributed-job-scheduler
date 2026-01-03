package imageservice

import (
	"context"

	databaseconstants "github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
)

type imageService struct {
	imageRepo repository.ImageRepository
}

func NewImageService(
	imageRepo repository.ImageRepository,
) ImageService {
	return &imageService{
		imageRepo: imageRepo,
	}
}

func (s *imageService) GetByTypeAndVersion(
	ctx context.Context,
	imageType databaseconstants.ImageType,
	version string,
) (models.ImageInformation, error) {
	return s.imageRepo.FindByTypeAndVersion(ctx, imageType, version)
}
