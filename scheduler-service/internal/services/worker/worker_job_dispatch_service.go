package worker

import (
	"context"

	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
)

type WorkerJobDispatchService interface {
	DispatchNextJob(ctx context.Context) (*contracts.JobToExecute, error)
}
