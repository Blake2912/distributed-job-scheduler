package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/contracts"
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

	jobExecutionsToInsert := make([]models.JobExecution, 0, len(jobIdToStatusMap))

	for jobId, status := range jobIdToStatusMap {
		jobExecutionsToInsert = append(jobExecutionsToInsert, createJobExecution(jobId, status))
	}

	return je.db.Create(&jobExecutionsToInsert).Error
}

func (je *JobExecutionRepository) GetJobExecutionInfoWithExecutionId(ctx context.Context, jobExecutionId uint) (models.JobExecution, error) {
	return gorm.G[models.JobExecution](je.db).
		Where("id = ?", jobExecutionId).
		First(ctx)
}

func (je *JobExecutionRepository) UpdateJobExecutions(ctx context.Context, jobExecutionUpdates map[uint]contracts.JobExecutionUpdate) error {

	if len(jobExecutionUpdates) == 0 {
		return nil
	}

	var (
		ids []uint

		statusCases     []string
		retryCountCases []string

		args []any
	)

	for id, upd := range jobExecutionUpdates {
		ids = append(ids, id)

		statusCases = append(statusCases, "WHEN ? THEN ?")
		args = append(args, id, upd.Status)

		if upd.RetryCount != nil {
			retryCountCases = append(retryCountCases, "WHEN ? THEN ?")
			args = append(args, id, *upd.RetryCount)
		}
	}

	var setParts []string

	if len(statusCases) > 0 {
		setParts = append(setParts, fmt.Sprintf(
			"status = CASE id %s END",
			strings.Join(statusCases, " "),
		))
	}

	if len(retryCountCases) > 0 {
		setParts = append(setParts, fmt.Sprintf(
			"retry_count = CASE id %s END",
			strings.Join(retryCountCases, " "),
		))
	}

	if len(setParts) == 0 {
		return nil
	}

	query := fmt.Sprintf(`
		UPDATE job_executions
		SET %s
		WHERE id IN ?;
	`, strings.Join(setParts, ", "))

	args = append(args, ids)

	return je.db.Exec(query, args...).Error
}

func createJobExecution(jobId uint, status database_constants.JobExecutionStatus) models.JobExecution {
	return models.JobExecution{
		JobID:      jobId,
		Status:     status,
		CreatedAt:  time.Now().UTC(),
		RetryCount: 1, // Keeping it as 1 for now tomorrow we can change this to make it configurable from service end
	}
}

func (r *JobExecutionRepository) MarkRunning(
	ctx context.Context,
	execID uint,
) error {

	result := r.db.WithContext(ctx).
		Model(&models.JobExecution{}).
		Where("id = ? AND status IN (?, ?)",
			execID,
			database_constants.Todo,
			database_constants.Retry,
		).
		Updates(map[string]interface{}{
			"status":     database_constants.Running,
			"started_at": time.Now(),
		})

	return result.Error
}
