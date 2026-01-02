package imageservice

import (
	"context"

	databaseconstants "example.com/constants/database_constants"
	"example.com/dal/sql_dal/models"
	"example.com/dal/sql_dal/repository"
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
