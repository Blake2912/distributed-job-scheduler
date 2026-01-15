package container

import (
	"context"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/httpclient"
	imagehandler "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/handlers/image_handler"
	jobshandler "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/handlers/jobs_handler"
	spawnworkersHandler "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/handlers/spawn_workers"
	ttlexpiryconsumers "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/consumers/ttl_expiry_consumers"
	imageservice "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/image_service"
	jobscheduling "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/job_scheduling"
	jobsservice "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/jobs_service"
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
	LeaderElector       leader.LeaderElection
	Scheduler           *scheduler.Scheduler
	JobsHandler         *jobshandler.JobsHandler
	WorkerHealthCheck   ttlexpiryconsumers.TTLExpiryConsumer
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
	jobsService := jobsservice.NewJobsService(jobRepository)
	workerHealthCheckService := ttlexpiryconsumers.NewTTLExpiryConsumer()

	leaderElector := leader.New(
		rdb,
		"scheduler:leader", //redis key
		10*time.Second,     //placeholder TTL
	)

	sched := scheduler.New(jobSchedulerService)

	// Build Handlers and return them
	imageHandler := imagehandler.New(imageService)
	spawnWorkersHandler := spawnworkersHandler.NewSpawnWokersHandler(spawnWorkerService)
	jobsHandler := jobshandler.New(jobsService)

	return &Container{
		ImageHandler:        imageHandler,
		SpawnWorkersHandler: spawnWorkersHandler,
		LeaderElector:       leaderElector,
		Scheduler:           sched,
		JobsHandler:         jobsHandler,
		WorkerHealthCheck:   workerHealthCheckService,
	}
}
