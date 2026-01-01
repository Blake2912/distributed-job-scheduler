package container

import (
	imagehandler "example.com/dal/internal/api/handlers/image_handler"
	imageservice "example.com/dal/internal/service/image_service"
	"example.com/dal/sql_dal/repository/postgres"
	"gorm.io/gorm"
)

type Container struct {
	ImageHandler *imagehandler.Handler
}

func BuildContainer(db *gorm.DB) *Container {
	// Build Repositories
	imageRepo := postgres.NewImageRepository(db)

	// Build Services
	imageService := imageservice.NewImageService(imageRepo)

	// Build Handlers and return them
	imageHandler := imagehandler.New(imageService)

	return &Container{
		ImageHandler: imageHandler,
	}
}
