package worker

import (
	"context"

	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
)

type workerService struct {
	jobExecutionRepository repository.JobExecutionRepository
}

func NewWorkerLeaseJobService(jobExecutionRepository repository.JobExecutionRepository) WorkerJobDispatchService {
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
		JobID:   exec.JobID,
		JobType: exec.Job.Type,
		Payload: exec.Job.Config,
	}

	return &executionPayload, nil
}
