package worker

import (
	"context"
	"strconv"

	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/common/infra_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/commands/queries"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/commands/queues"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
)

type workerService struct {
	redisQueueCommands     queues.RedisQueueCommands
	jobRepository          repository.JobRepository
	jobExecutionRepository repository.JobExecutionRepository
	redisQueries           queries.Queries
}

func NewWorkerJobDispatchService(redisQueueCommands queues.RedisQueueCommands,
	jobRepository repository.JobRepository,
	jobExecutionRepository repository.JobExecutionRepository,
	redisQueries queries.Queries,
) WorkerJobDispatchService {
	return &workerService{
		redisQueueCommands:     redisQueueCommands,
		jobRepository:          jobRepository,
		jobExecutionRepository: jobExecutionRepository,
		redisQueries:           redisQueries,
	}
}

func (svc *workerService) DispatchNextJob(ctx context.Context) (*contracts.JobToExecute, error) {

	jobIdString, err := svc.redisQueueCommands.RDequeue(ctx, infra_constants.JobsQueue)

	if err != nil {
		return nil, err
	}

	if jobIdString == "" {
		return nil, nil
	}

	jobId, err := strconv.ParseUint(jobIdString, 10, 64)

	if err != nil {
		return nil, err
	}

	exec, err := svc.jobExecutionRepository.GetLatestJobExecutions(ctx, []uint{uint(jobId)})
	if err != nil {
		return nil, err
	}

	if len(exec) == 0 {
		//stale entry in queue, ignore
		return nil, nil
	}

	job, err := svc.jobRepository.GetJobById(ctx, uint(jobId))
	if err != nil {
		return nil, err
	}

	executionPayload := contracts.JobToExecute{
		JobID:   job.ID,
		JobType: job.Type,
		Payload: job.Config,
	}

	for _, execution := range exec {
		if execution.Status != database_constants.Todo &&
			execution.Status != database_constants.Retry {
			//stale entry in queue, ignore
			continue
		}

		//state transition
		err := svc.jobExecutionRepository.MarkRunning(ctx, execution.ID)
		if err != nil {
			return nil, err
		}
	}

	return &executionPayload, nil
}
