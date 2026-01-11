package jobscheduling

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
	"github.com/Blake2912/distributed-job-scheduler/common/infra_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/commands"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
)

type jobSchedulingService struct {
	redisQueueCommands     commands.RedisQueueCommands
	jobRepository          repository.JobRepository
	jobExecutionRepository repository.JobExecutionRepository
}

func NewJobSchedulingService(
	redisQueueCommands commands.RedisQueueCommands,
	jobRepository repository.JobRepository,
	jobExecutionRepository repository.JobExecutionRepository,
) JobSchedulingService {
	return &jobSchedulingService{
		redisQueueCommands:     redisQueueCommands,
		jobRepository:          jobRepository,
		jobExecutionRepository: jobExecutionRepository,
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
				if jobExecution.Status == database_constants.Error {
					newJobsIdsToPushToQueue = append(newJobsIdsToPushToQueue, job)
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
	pushedJobs := make([]uint, len(validJobIds))
	for _, job := range newJobsIdsToPushToQueue {
		jobInStr := strconv.FormatUint(uint64(job), 10)
		queueErr := j.redisQueueCommands.LEnqueue(ctx, infra_constants.JobsQueue, jobInStr)
		if queueErr != nil {
			log.Printf("An error occurred while pushing the job %d to redis queue %s | Will retry to en-queue job in the next poll", job, queueErr)
			// Removing from insert so that db and redis are consistent
			delete(newJobsToExecute, job)
		}
		pushedJobs = append(pushedJobs, job)
	}
	log.Printf("Completed pushing jobs to queue now inserting executions into database")

	return j.jobExecutionRepository.InsertNewJobExecutions(ctx, newJobsToExecute)
}
