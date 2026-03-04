package jobsservice

import (
	"context"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/contracts/jobs"
)

type JobsService interface {
	CreateJobs(ctx context.Context, payload []jobs.CreateJobsPayload) error
}
