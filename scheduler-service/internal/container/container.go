package container

import (
	"context"
	"log"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/httpclient"
	imagehandler "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/handlers/image_handler"
	spawnworkersHandler "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/handlers/spawn_workers"
	imageservice "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/image_service"
	leader "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/leader_election"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/scheduler"
	spawnworkers "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/spawn_workers"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/pod_library/client"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository/postgres"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	ImageHandler        *imagehandler.Handler
	SpawnWorkersHandler *spawnworkersHandler.SpawnWorkersHandler
}

func BuildContainer(db *gorm.DB, rdb *redis.Client, ctx context.Context) *Container {
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
	leaderElector := leader.New(
		rdb,
		"scheduler:leader", //redis key
		10*time.Second,     //placeholder TTL
	)

	sched := scheduler.New()

	err = leaderElector.Run(ctx, func(leaderCtx context.Context) {
		sched.Run(leaderCtx)
	})

	if err != nil {
		log.Fatal(err)
	}

	// Build Handlers and return them
	imageHandler := imagehandler.New(imageService)
	spawnWorkersHandler := spawnworkersHandler.NewSpawnWokersHandler(spawnWorkerService)

	return &Container{
		ImageHandler:        imageHandler,
		SpawnWorkersHandler: spawnWorkersHandler,
	}
}
