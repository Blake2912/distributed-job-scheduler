package repository

import (
	"context"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
)

type JobRepository interface {
	GetJobsToSchedule(ctx context.Context, currentTime time.Time) ([]models.Jobs, error)
	CreateJobs(ctx context.Context, jobs []models.Jobs) error
	GetJobById(ctx context.Context, id uint) (models.Jobs, error)
}
