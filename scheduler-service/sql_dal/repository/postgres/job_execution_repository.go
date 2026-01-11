package postgres

import (
	"context"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
	"gorm.io/gorm"
)

type JobExecutionRepository struct {
	db *gorm.DB
}

func NewJobExecutionRepository(db *gorm.DB) repository.JobExecutionRepository {
	return &JobExecutionRepository{
		db: db,
	}
}

func (je *JobExecutionRepository) GetLatestJobExecutions(ctx context.Context, jobIds []uint) ([]models.JobExecution, error) {
	var executions []models.JobExecution

	err := je.db.WithContext(ctx).
		Raw(`
            SELECT DISTINCT ON (job_id)
                *
            FROM job_executions
            WHERE job_id IN ?
            ORDER BY job_id, created_at DESC
        `, jobIds).
		Scan(&executions).Error

	return executions, err
}

func (je *JobExecutionRepository) InsertNewJobExecutions(ctx context.Context, jobIdToStatusMap map[uint]database_constants.JobExecutionStatus) error {

	jobExecutionsToInsert := make([]models.JobExecution, len(jobIdToStatusMap))

	for jobId, status := range jobIdToStatusMap {
		jobExecutionsToInsert = append(jobExecutionsToInsert, createJobExecution(jobId, status))
	}
	
	return je.db.Create(&jobExecutionsToInsert).Error
}

func createJobExecution(jobId uint, status database_constants.JobExecutionStatus) models.JobExecution {
	return models.JobExecution{
		JobID:      jobId,
		Status:     status,
		CreatedAt:  time.Now().UTC(),
		RetryCount: 1, // Keeping it as 1 for now tomorrow we can change this to make it configurable from service end
	}
}
