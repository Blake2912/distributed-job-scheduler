package container

import (
	"log"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/httpclient"
	imagehandler "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/handlers/image_handler"
	spawnworkersHandler "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/handlers/spawn_workers"
	imageservice "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/image_service"
	spawnworkers "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/spawn_workers"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/pod_library/client"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository/postgres"
	"gorm.io/gorm"
)

type Container struct {
	ImageHandler        *imagehandler.Handler
	SpawnWorkersHandler *spawnworkersHandler.SpawnWorkersHandler
}

func BuildContainer(db *gorm.DB) *Container {
	// Build Repositories
	imageRepo := postgres.NewImageRepository(db)

	httpClient := httpclient.New(120 * time.Minute)
	k8sClient, err := client.New()
	if err != nil {
		log.Fatal(err)
	}

	// Build Services
	imageService := imageservice.NewImageService(imageRepo)
	spawnWorkerService := spawnworkers.NewSpawnWorkerService(*httpClient, *k8sClient)

	// Build Handlers and return them
	imageHandler := imagehandler.New(imageService)
	spawnWorkersHandler := spawnworkersHandler.NewSpawnWokersHandler(spawnWorkerService)

	return &Container{
		ImageHandler:        imageHandler,
		SpawnWorkersHandler: spawnWorkersHandler,
	}
}
