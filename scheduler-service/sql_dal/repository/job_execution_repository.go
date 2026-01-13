package repository

import (
	"context"

	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
)

type JobExecutionRepository interface {
	GetLatestJobExecutions(ctx context.Context, jobIds []uint) ([]models.JobExecution, error)
	InsertNewJobExecutions(ctx context.Context, jobIdToStatusMap map[uint]database_constants.JobExecutionStatus) error
}
