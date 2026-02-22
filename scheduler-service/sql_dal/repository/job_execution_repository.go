package repository

import (
	"context"

	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/contracts"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
)

type JobExecutionRepository interface {
	GetLatestJobExecutions(ctx context.Context, jobIds []uint) ([]models.JobExecution, error)
	InsertNewJobExecutions(ctx context.Context, jobIdToStatusMap map[uint]contracts.JobExecutionCreationData) error
	GetJobExecutionInfoWithExecutionId(ctx context.Context, jobExecutionId uint) (models.JobExecution, error)
	UpdateJobExecutions(ctx context.Context, jobExecutionUpdates map[uint]contracts.JobExecutionUpdate) error
	GetJobAndMarkExecutionAsRunning(ctx context.Context) (*models.JobExecution, error)
	UpdateJobExecutionStatus(ctx context.Context, execId uint, status database_constants.JobExecutionStatus) error
	MarkExpiredLeasesAsRetry(ctx context.Context) error
}
