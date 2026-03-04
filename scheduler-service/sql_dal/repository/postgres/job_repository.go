package postgres

import (
	"context"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
	"gorm.io/gorm"
)

type JobRepository struct {
	db *gorm.DB
}

func NewJobRepository(db *gorm.DB) repository.JobRepository {
	return &JobRepository{
		db: db,
	}
}

func (j *JobRepository) GetJobsToSchedule(ctx context.Context, currentTime time.Time) ([]models.Jobs, error) {
	currentTimeOfDay := currentTime.UTC().Format("15:04:05")

	return gorm.G[models.Jobs](j.db).
		Where("next_run_at <= ? AND enabled = ?", currentTimeOfDay, true).
		Find(ctx)
}

func (j *JobRepository) CreateJobs(ctx context.Context, jobs []models.Jobs) error {
	return j.db.Create(&jobs).Error
}

func (j *JobRepository) GetJobById(ctx context.Context, id uint) (models.Jobs, error) {
	return gorm.G[models.Jobs](j.db).
		Where("id = ?", id).
		First(ctx)
}
