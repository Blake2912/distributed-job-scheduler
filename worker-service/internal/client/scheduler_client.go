package client

import (
	"context"

	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
)

// SchedulerClient defines scheduler API contract
type SchedulerClient interface {
	LeaseJob(ctx context.Context) (*contracts.JobToExecute, error)
	ReportCompletion(ctx context.Context, jobId uint, status string, executionError string, retryable bool) error
	ExtendLease(ctx context.Context, execId uint) error
}

/*
Expected Scheduler Endpoints

POST /jobs/lease
POST /jobs/{id}/completion
*/
