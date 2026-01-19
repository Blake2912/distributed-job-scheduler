package jobscheduling

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/common/helpers"
	"github.com/Blake2912/distributed-job-scheduler/common/infra_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/commands/queries"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/commands/queues"
	dalcontracts "github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/contracts"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
)

type jobSchedulingService struct {
	redisQueueCommands     queues.RedisQueueCommands
	jobRepository          repository.JobRepository
	jobExecutionRepository repository.JobExecutionRepository
	redisQueries           queries.Queries
}

func NewJobSchedulingService(
	redisQueueCommands queues.RedisQueueCommands,
	jobRepository repository.JobRepository,
	jobExecutionRepository repository.JobExecutionRepository,
	redisQueries queries.Queries,
) JobSchedulingService {
	return &jobSchedulingService{
		redisQueueCommands:     redisQueueCommands,
		jobRepository:          jobRepository,
		jobExecutionRepository: jobExecutionRepository,
		redisQueries:           redisQueries,
	}
}

func (j *jobSchedulingService) ScheduleJobs(ctx context.Context) error {

	t := time.Now().In(time.UTC)

	timeOnly := time.Date(0, 1, 1,
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		time.UTC,
	)
	// Pick jobs to process
	validJobs, error := j.jobRepository.GetJobsToSchedule(ctx, timeOnly)
	if error != nil {
		log.Printf("An error occurred while quering jobs data %s \n", error.Error())
		return error
	}

	if len(validJobs) == 0 {
		log.Println("No jobs found to schedule returning.")
		return nil
	}
	// Pick existing executions
	validJobIds := make([]uint, 0, len(validJobs))
	validJobsMap := make(map[uint]models.Jobs)

	for _, job := range validJobs {
		validJobIds = append(validJobIds, job.ID)
		validJobsMap[job.ID] = job
	}

	latestExecutions, error := j.jobExecutionRepository.GetLatestJobExecutions(ctx, validJobIds)
	if error != nil {
		log.Printf("An error occurred while quering jobs execution data %s \n", error.Error())
		return error
	}

	latestExecutionsMap := make(map[uint]models.JobExecution)

	for _, exec := range latestExecutions {
		latestExecutionsMap[exec.JobID] = exec
	}

	newJobsToExecute := make(map[uint]database_constants.JobExecutionStatus, len(validJobIds))
	newJobsIdsToPushToQueue := make([]uint, 0, len(validJobIds))

	today := time.Now().In(time.UTC)
	jobExecutionIdToStatusMap := make(map[uint]dalcontracts.JobExecutionUpdate)

	for _, job := range validJobIds {
		jobExecution, execFound := latestExecutionsMap[job]
		jobInfo := validJobsMap[job]

		var jobMetaData contracts.JobMetaDataContract

		err := json.Unmarshal([]byte(jobInfo.Metadata), &jobMetaData)

		if err != nil {
			log.Printf("An error occurred while deserializing Metadata field in jobs %s | So skipping scheduling of job: %d", err.Error(), job)
			continue
		}

		if execFound {
			createdAt := jobExecution.CreatedAt

			// Job already executed for the day so skip it
			if today.Year() == createdAt.Year() && today.YearDay() == createdAt.YearDay() {

				if jobExecution.Status == database_constants.Error && jobExecution.RetryCount > 0 {
					newRetryCount := int(jobExecution.RetryCount - 1)

					jobExecutionIdToStatusMap[jobExecution.ID] = dalcontracts.JobExecutionUpdate{
						Status:     database_constants.Retry,
						RetryCount: &newRetryCount,
					}

					if err = j.jobExecutionRepository.UpdateJobExecutions(ctx, jobExecutionIdToStatusMap); err != nil {
						log.Printf("An error occurred while updating job execution data %s", err.Error())
						continue
					}

					newJobsIdsToPushToQueue = append(newJobsIdsToPushToQueue, job)
					continue
				}

				if jobExecution.Status == database_constants.Todo {
					newJobsIdsToPushToQueue = append(newJobsIdsToPushToQueue, job)
					continue
				}

				if jobExecution.Status == database_constants.Retry && jobExecution.RetryCount > 0 {

					newRetryCount := int(jobExecution.RetryCount - 1)

					jobExecutionIdToStatusMap[jobExecution.ID] = dalcontracts.JobExecutionUpdate{
						Status:     database_constants.Retry,
						RetryCount: &newRetryCount,
					}

					if err = j.jobExecutionRepository.UpdateJobExecutions(ctx, jobExecutionIdToStatusMap); err != nil {
						log.Printf("An error occurred while updating job execution data %s", err.Error())
						continue
					}

					newJobsIdsToPushToQueue = append(newJobsIdsToPushToQueue, job)
					continue
				}

				if jobExecution.Status == database_constants.Running {
					executionIdInString := strconv.FormatUint(uint64(jobExecution.ID), 10)
					workerInfo, err := j.redisQueries.GetValue(ctx, executionIdInString)
					if err != nil {
						log.Printf("An error occurred while fetching the key information %s", err.Error())
						continue
					}

					healthCheckKey := helpers.BuildHealthCheckKey(workerInfo, executionIdInString)
					isWorkerRunning, err := j.redisQueries.KeyExists(ctx, healthCheckKey)

					if err != nil {
						log.Printf("An error occurred while fetching the healthCheck info %s", err.Error())
						continue
					}

					if isWorkerRunning {
						log.Printf("Worker %s is already running for this job execution %s so ignoring to reschedule", workerInfo, executionIdInString)
						continue
					}

					newRetryCount := int(jobExecution.RetryCount - 1)

					jobExecutionIdToStatusMap[jobExecution.ID] = dalcontracts.JobExecutionUpdate{
						Status:     database_constants.Retry,
						RetryCount: &newRetryCount,
					}

					if err = j.jobExecutionRepository.UpdateJobExecutions(ctx, jobExecutionIdToStatusMap); err != nil {
						log.Printf("An error occurred while updating job execution data %s", err.Error())
						continue
					}

					deletedCount, err := j.redisQueries.DeleteKey(ctx, executionIdInString)
					if err != nil {
						log.Printf("An error occurred while deleting the job to worker mapping %s", err)
						continue
					}
					if deletedCount == 0 {
						log.Printf("No keys were deleted either the key %s never existed or the key was deleted earlier", executionIdInString)
					}

					newJobsIdsToPushToQueue = append(newJobsIdsToPushToQueue, job)
					continue
				}

				continue
			}

			// Schedule only if the job is a recurring job
			if jobMetaData.IsRecurringJob {
				newJobsToExecute[job] = database_constants.Todo
				newJobsIdsToPushToQueue = append(newJobsIdsToPushToQueue, job)
			}

		} else {
			newJobsToExecute[job] = database_constants.Todo
			newJobsIdsToPushToQueue = append(newJobsIdsToPushToQueue, job)
		}
	}

	// Push the jobs to queue
	log.Printf("Starting to push the jobs into redis queue  %+v", newJobsIdsToPushToQueue)
	for _, job := range newJobsIdsToPushToQueue {
		jobInStr := strconv.FormatUint(uint64(job), 10)
		queueErr := j.redisQueueCommands.LEnqueue(ctx, infra_constants.JobsQueue, jobInStr)
		if queueErr != nil {
			log.Printf("An error occurred while pushing the job %d to redis queue %s | Will retry to en-queue job in the next poll", job, queueErr)
			// Removing from insert so that db and redis are consistent
			delete(newJobsToExecute, job)
		}
	}
	log.Printf("Completed pushing jobs to queue now inserting executions into database")

	return j.jobExecutionRepository.InsertNewJobExecutions(ctx, newJobsToExecute)
}
