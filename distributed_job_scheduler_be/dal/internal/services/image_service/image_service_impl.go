package imageservice

import (
	"context"

	databaseconstants "example.com/constants/database_constants"
	"example.com/dal/sql_dal/models"
	"example.com/dal/sql_dal/repository"
)

type service struct {
	imageRepo repository.ImageRepository
}

func NewImageService(
	imageRepo repository.ImageRepository,
) Service {
	return &service{
		imageRepo: imageRepo,
	}
}

func (s *service) GetByTypeAndVersion(
	ctx context.Context,
	imageType databaseconstants.ImageType,
	version string,
) (models.ImageInformation, error) {
	return s.imageRepo.FindByTypeAndVersion(ctx, imageType, version)
}
