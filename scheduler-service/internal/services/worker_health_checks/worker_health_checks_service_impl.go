package workerhealthchecks

import (
	"context"
	"log"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/helpers"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/commands/queries"
)

type workerHealthChecks struct {
	redisQuery queries.Queries
}

func NewWorkerHealthCheck(redisQuery queries.Queries) WorkerHealthChecks {
	return &workerHealthChecks{
		redisQuery: redisQuery,
	}
}

// Health check for worker
func (w *workerHealthChecks) CheckHealth(ctx context.Context, workerId string, jobExecutionId string) error {

	healthCheckKey := helpers.BuildHealthCheckKey(workerId, jobExecutionId)
	ttl := 15 * time.Second

	err := w.redisQuery.CheckAndSetKeyWithTTL(ctx, healthCheckKey, jobExecutionId, ttl)

	if err != nil {
		log.Printf("In HealthCheck | An error occurred while setting redis key for %s with value %s with err %s", healthCheckKey, jobExecutionId, err.Error())
		return err
	}

	err = w.redisQuery.CheckAndSetKeyWithTTL(ctx, jobExecutionId, workerId, -2)

	if err != nil {
		log.Printf("In HealthCheck | An error occurred while setting redis key for %s with value %s with err %s", healthCheckKey, jobExecutionId, err.Error())
		return err
	}

	return nil
}

// Delete existing keys for a worker
func (w *workerHealthChecks) DeleteWorkerKeys(ctx context.Context, workerId string, jobExecutionId string) error {

	healthCheckKey := helpers.BuildHealthCheckKey(workerId, jobExecutionId)

	res, err := w.redisQuery.DeleteKey(ctx, healthCheckKey)

	if err != nil {
		return err
	}

	if res == 0 {
		log.Printf("Tried to delete %s but it doesn't exist", healthCheckKey)
	}

	res, err = w.redisQuery.DeleteKey(ctx, jobExecutionId)

	if err != nil {
		return err
	}

	if res == 0 {
		log.Printf("Tried to delete %s but it doesn't exist", jobExecutionId)
	}
	return nil

}
