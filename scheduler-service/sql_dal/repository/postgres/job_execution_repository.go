package postgres

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/constants/worker_constants"
	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/contracts"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
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

func (je *JobExecutionRepository) InsertNewJobExecutions(ctx context.Context, jobIdToStatusMap map[uint]contracts.JobExecutionCreationData) error {

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

		statusCases []string
		//retryCountCases []string

		args []any
	)

	for id, upd := range jobExecutionUpdates {
		ids = append(ids, id)

		statusCases = append(statusCases, "WHEN ? THEN ?")
		args = append(args, id, upd.Status)

		/*
			if upd.RetryCount != nil {
				retryCountCases = append(retryCountCases, "WHEN ? THEN ?")
				args = append(args, id, *upd.RetryCount)
			}
		*/
	}

	var setParts []string

	if len(statusCases) > 0 {
		setParts = append(setParts, fmt.Sprintf(
			"status = CASE id %s END",
			strings.Join(statusCases, " "),
		))
	}

	/*
		if len(retryCountCases) > 0 {
			setParts = append(setParts, fmt.Sprintf(
				"retry_count = CASE id %s END",
				strings.Join(retryCountCases, " "),
			))
		}
	*/

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

func createJobExecution(jobId uint, data contracts.JobExecutionCreationData) models.JobExecution {
	return models.JobExecution{
		JobID:              jobId,
		Status:             database_constants.JobExecutionStatus(data.Status),
		ScheduledAt:        data.ScheduledAt,
		CreatedAt:          time.Now().UTC(),
		MaxAttemptsAllowed: 1, // Keeping it as 1 for now tomorrow we can change this to make it configurable from service end
	}
}

func (r *JobExecutionRepository) GetJobAndMarkExecutionAsRunning(ctx context.Context) (*models.JobExecution, error) {

	tx := r.db.WithContext(ctx).Begin()

	var exec models.JobExecution

	err := tx.Preload("Job").
		Clauses(clause.Locking{Strength: "UPDATE", Options: "SKIP LOCKED"}).
		Where(`
			(status = ? AND scheduled_at <= NOW())
			OR
			(status = ? AND retry_at <= NOW())
		`,
			database_constants.Todo,
			database_constants.Retry,
		).
		Order("scheduled_at ASC").
		Limit(1).
		First(&exec).Error

	if err == gorm.ErrRecordNotFound {
		tx.Rollback()
		return nil, nil
	}

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	leaseExpiry := time.Now().Add(worker_constants.LeaseDuration)

	err = tx.Model(&models.JobExecution{}).
		Where("id = ?", exec.ID).
		Updates(map[string]interface{}{
			"status":        database_constants.Running,
			"lease_expiry":  leaseExpiry,
			"attempt_count": gorm.Expr("attempt_count + 1"),
		}).Error

	if err != nil {
		tx.Rollback()
		return nil, err
	}

	tx.Commit()

	return &exec, nil
}

func (r *JobExecutionRepository) UpdateJobExecutionStatus(ctx context.Context, execId uint, status database_constants.JobExecutionStatus) error {
	return r.db.Model(&models.JobExecution{}).
		Where("id = ?", execId).
		Update("status", status).Error
}

func (r *JobExecutionRepository) MarkExpiredLeasesAsRetry(ctx context.Context) error {
	return r.db.WithContext(ctx).
		Model(&models.JobExecution{}).
		Where("status = ? AND lease_expiry < NOW()", database_constants.Running).
		Updates(map[string]interface{}{
			"status":   "RETRY",
			"retry_at": time.Now(),
		}).Error
}

func (r *JobExecutionRepository) ExtendLease(ctx context.Context, jobExecId uint) error {
	leaseExpiry := time.Now().Add(worker_constants.LeaseDuration)
	return r.db.WithContext(ctx).
		Model(&models.JobExecution{}).
		Where("id = ? AND status = ? AND lease_expiry > NOW()", jobExecId, database_constants.Running).
		Updates(map[string]interface{}{
			"lease_expiry": leaseExpiry,
		}).Error
}
