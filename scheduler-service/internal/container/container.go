package container

import (
	imagehandler "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/handlers/image_handler"
	imageservice "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/image_service"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository/postgres"
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
