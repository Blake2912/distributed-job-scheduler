package container

import (
	"net/http"

	"github.com/Blake2912/distributed-job-scheduler/common/constants/worker_constants"
	"github.com/Blake2912/distributed-job-scheduler/worker-service/internal/app"
	"github.com/Blake2912/distributed-job-scheduler/worker-service/internal/client"
	"github.com/Blake2912/distributed-job-scheduler/worker-service/internal/executor"
	"github.com/Blake2912/distributed-job-scheduler/worker-service/internal/worker"
)

type Container struct {
	//App
	App *app.App

	//registry
	ExecutorRegistry executor.Registry

	// Core execution
	Worker *worker.Worker

	// Clients
	SchedulerClient client.SchedulerClient

	// Optional: handlers / executors
	JobExecutor executor.Executor
}

func BuildContainer(
	httpClient *http.Client,
	schedulerBaseURL string,
	workerCount int,
) *Container {
	//client
	schedulerClient := client.NewHTTPSchedulerClient(
		schedulerBaseURL,
		httpClient,
	)

	//executors
	jobExecutor := executor.NewWebhookExecutor(httpClient)

	//executor registry
	registry := executor.NewRegistry()

	err := registry.Register(worker_constants.HTTPWebhookExecutor, jobExecutor)
	if err != nil {
		panic(err)
	}

	//worker
	workerInstance := worker.New(
		schedulerClient,
		registry,
	)

	appInstance := app.New(workerInstance, workerCount)

	return &Container{
		App:             appInstance,
		Worker:          workerInstance,
		SchedulerClient: schedulerClient,
		JobExecutor:     jobExecutor,
	}
}
