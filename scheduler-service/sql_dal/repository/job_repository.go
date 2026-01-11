package repository

import (
	"context"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
)

type JobRepository interface {
	GetJobsToSchedule(ctx context.Context, currentTime time.Time) ([]models.Jobs, error)
}
