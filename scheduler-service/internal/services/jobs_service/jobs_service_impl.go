package jobsservice

import (
	"context"

	"github.com/Blake2912/distributed-job-scheduler/common/contracts/jobs"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
)

type jobsService struct {
	jobRepository repository.JobRepository
}

func NewJobsService(jobRepository repository.JobRepository) JobsService {
	return &jobsService{
		jobRepository: jobRepository,
	}
}

// Create Jobs
func (j *jobsService) CreateJobs(ctx context.Context, payload []jobs.CreateJobsPayload) error {
	
	
	return nil
}

// View jobs
