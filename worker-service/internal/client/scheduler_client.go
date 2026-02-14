package client

import (
	"context"
)

// LeaseJobResponse represents job lease from scheduler
type LeaseJobResponse struct {
	JobID   string
	JobType string
	Payload string
}

type ReportCompletionRequest struct {
	Status    string
	Error     string
	Retryable bool
}

// SchedulerClient defines scheduler API contract
type SchedulerClient interface {
	LeaseJob(ctx context.Context) (*LeaseJobResponse, error)
	ReportCompletion(ctx context.Context, jobId string, status string, executionError string, retryable bool) error
}

/*
Expected Scheduler Endpoints

POST /jobs/lease
POST /jobs/{id}/completion
*/
