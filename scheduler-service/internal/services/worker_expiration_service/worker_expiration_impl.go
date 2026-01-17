package workerexpirationservice

import (
	"context"
	"log"
	"strconv"
	"strings"

	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/state"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/commands/queries"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
)

type workerExpiration struct {
	jobExectionRepository repository.JobExecutionRepository
	redisQuery            queries.Queries
}

func NewWorkerExpiry(jobExectionRepository repository.JobExecutionRepository) WorkerExpiration {
	return &workerExpiration{
		jobExectionRepository: jobExectionRepository,
	}
}

func (w *workerExpiration) HandleWorkerExpiry(ctx context.Context, expiredKey string) {

	if !state.IsLeader() {
		log.Printf("Worker expiry was picked but for key %s but ignore as its not leader", expiredKey)
		return
	}

	if !strings.Contains(expiredKey, database_constants.HEALTH_CHECK_KEY_IDENTIFIER) {
		log.Printf("Ingoring %s key as no operation required here", expiredKey)
	}

	splitString := strings.Split(expiredKey, "_")
	workerId := splitString[1]
	jobExecutionId := splitString[2]

	log.Printf("Worker expiry picked up for worker id: %s and Job execution id: %s", workerId, jobExecutionId)

	jobExecutionIdInInt, err := strconv.ParseUint(jobExecutionId, 10, 64)

	if err != nil {
		log.Printf("An error occurred while converting the jobExecutionId to integer %s", err.Error())
		return
	}

	jobExecutionInfo, err := w.jobExectionRepository.GetJobExecutionInfoWithExecutionId(ctx, uint(jobExecutionIdInInt))

	if err != nil {
		log.Printf("An error occurred while fetching the job execution information for %d. Error: %s", jobExecutionIdInInt, err.Error())
		return
	}

	if jobExecutionInfo.Status == database_constants.Completed {
		log.Printf("The jobId %d has already completed so ignoring execution", jobExecutionInfo.JobID)
		return
	}

	jobExecutionIdToStatusMap := make(map[uint]database_constants.JobExecutionStatus)

	jobExecutionIdToStatusMap[jobExecutionInfo.ID] = database_constants.Retry

	err = w.jobExectionRepository.UpdateJobExecutionStatus(ctx, jobExecutionIdToStatusMap)

	if err != nil {
		log.Printf("Error occurred while updating job execution status %s", err.Error())
	}
}
