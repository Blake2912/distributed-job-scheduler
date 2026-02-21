package worker

import (
	"context"

	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
)

type WorkerLeaseJobService interface {
	LeaseJob(ctx context.Context) (*contracts.JobToExecute, error)
	CompleteJobExecution(ctx context.Context, jobExecId uint, request contracts.ReportCompletionRequest) error
}
