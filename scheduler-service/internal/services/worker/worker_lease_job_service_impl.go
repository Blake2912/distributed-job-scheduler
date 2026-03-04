package worker

import (
	"context"
	"fmt"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/constants/worker_constants"
	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	dalcontracts "github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/contracts"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
)

type workerService struct {
	jobExecutionRepository repository.JobExecutionRepository
}

func NewWorkerLeaseJobService(jobExecutionRepository repository.JobExecutionRepository) WorkerLeaseJobService {
	return &workerService{
		jobExecutionRepository: jobExecutionRepository,
	}
}

func (svc *workerService) LeaseJob(ctx context.Context) (*contracts.JobToExecute, error) {

	exec, err := svc.jobExecutionRepository.GetJobAndMarkExecutionAsRunning(ctx)
	if err != nil {
		return nil, err
	}

	if exec == nil {
		return nil, nil
	}

	executionPayload := contracts.JobToExecute{
		JobExecutionID: exec.ID,
		JobType:        exec.Job.Type,
		Payload:        exec.Job.Config,
	}

	return &executionPayload, nil
}

func (svc *workerService) CompleteJobExecution(ctx context.Context, jobExecId uint, request contracts.ReportCompletionRequest) error {
	exec, err := svc.jobExecutionRepository.GetJobExecutionInfoWithExecutionId(ctx, jobExecId)
	if err != nil {
		return err
	}

	//safe check for idempotency, meaning the job has already been executed
	if exec.Status == database_constants.Completed || exec.Status == database_constants.Error {
		return nil
	}

	if exec.Status != database_constants.Running {
		return fmt.Errorf("invalid state transition")
	}

	if request.Status == worker_constants.SUCCESS {
		return svc.jobExecutionRepository.UpdateJobExecutionStatus(ctx, exec.ID, database_constants.Completed)
	}

	if request.Status == worker_constants.FAILED {
		if request.Retryable && exec.AttemptCount < exec.MaxAttemptsAllowed {
			backoff := calculateBackoff(int(exec.AttemptCount))
			retryAt := time.Now().Add(backoff)

			updatePayload := dalcontracts.JobExecutionUpdate{
				Status:  database_constants.Retry,
				RetryAt: retryAt,
			}

			return svc.jobExecutionRepository.UpdateJobExecutions(ctx, map[uint]dalcontracts.JobExecutionUpdate{exec.ID: updatePayload})
		}

		//update comments with error payload
		return svc.jobExecutionRepository.UpdateJobExecutionStatus(ctx, exec.ID, database_constants.Error)
	}

	return fmt.Errorf("invalid completion status")
}

func calculateBackoff(attempt int) time.Duration {
	base := time.Second
	max := 30 * time.Second

	delay := base * time.Duration(1<<uint(attempt-1))
	if delay > max {
		return max
	}

	return delay
}

func (svc *workerService) ExtendJobLease(ctx context.Context, jobExecId uint) error {
	return svc.jobExecutionRepository.ExtendLease(ctx, jobExecId)
}
