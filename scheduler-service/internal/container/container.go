package container

import (
	"context"
	"log"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/httpclient"
	imagehandler "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/handlers/image_handler"
	spawnworkersHandler "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/handlers/spawn_workers"
	imageservice "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/image_service"
	jobscheduling "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/job_scheduling"
	leader "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/leader_election"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/scheduler"
	spawnworkers "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/spawn_workers"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/pod_library/client"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/commands"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository/postgres"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Container struct {
	ImageHandler        *imagehandler.Handler
	SpawnWorkersHandler *spawnworkersHandler.SpawnWorkersHandler
}

func BuildContainer(db *gorm.DB, rdb *redis.Client, ctx context.Context, httpClient *httpclient.Client, k8sClient *client.K8sClient) *Container {
	// Build Common dependencies
	redisQueueCommands := commands.NewRedisQueueCommands(rdb)

	// Build Repositories
	imageRepo := postgres.NewImageRepository(db)
	jobRepository := postgres.NewJobRepository(db)
	jobExecutionRepository := postgres.NewJobExecutionRepository(db)

	// Build Services
	imageService := imageservice.NewImageService(imageRepo)
	spawnWorkerService := spawnworkers.NewSpawnWorkerService(httpClient, k8sClient)
	jobSchedulerService := jobscheduling.NewJobSchedulingService(redisQueueCommands, jobRepository, jobExecutionRepository)

	leaderElector := leader.New(
		rdb,
		"scheduler:leader", //redis key
		10*time.Second,     //placeholder TTL
	)

	sched := scheduler.New(jobSchedulerService)

	err := leaderElector.Run(ctx, func(leaderCtx context.Context) {
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
