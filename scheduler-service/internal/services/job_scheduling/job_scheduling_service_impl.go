package jobscheduling

import (
	"context"
	"encoding/json"
	"log"
	"strconv"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
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

	today := time.Now().UTC()

	/*
		timeOnly := time.Date(0, 1, 1,
			t.Hour(),
			t.Minute(),
			t.Second(),
			t.Nanosecond(),
			time.UTC,
		)*/
	// Pick jobs to process
	validJobs, error := j.jobRepository.GetJobsToSchedule(ctx, today)
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
		log.Printf("An error occurred while querying jobs execution data %s \n", error.Error())
		return error
	}

	latestExecutionsMap := make(map[uint]models.JobExecution)

	for _, exec := range latestExecutions {
		latestExecutionsMap[exec.JobID] = exec
	}

	newJobsToExecute := make(map[uint]dalcontracts.JobExecutionCreationData, len(validJobIds))

	for _, job := range validJobIds {
		jobExecution, execFound := latestExecutionsMap[job]
		jobInfo := validJobsMap[job]

		var jobMetaData contracts.JobMetaDataContract

		err := json.Unmarshal([]byte(jobInfo.Metadata), &jobMetaData)

		if err != nil {
			log.Printf("An error occurred while deserializing Metadata field in jobs %s | So skipping scheduling of job: %d", err.Error(), job)
			continue
		}

		scheduledAtTime := GetScheduledAtTime(jobInfo.NextRunAt)

		if execFound {
			createdAt := jobExecution.CreatedAt
			if today.Year() == createdAt.Year() && today.YearDay() == createdAt.YearDay() {
				//execution has already been created
				continue
			}

			//For recurring jobs create new execution
			if jobMetaData.IsRecurringJob {
				newJobsToExecute[job] = dalcontracts.JobExecutionCreationData{
					Status:      string(database_constants.Todo),
					ScheduledAt: scheduledAtTime,
				}
			}

		} else {
			newJobsToExecute[job] = dalcontracts.JobExecutionCreationData{
				Status:      string(database_constants.Todo),
				ScheduledAt: scheduledAtTime,
			}
		}
	}

	if len(newJobsToExecute) == 0 {
		log.Println("No new jobs to schedule.")
		return nil
	}

	log.Printf("Scheduling %s jobs.", strconv.Itoa(len(newJobsToExecute)))

	return j.jobExecutionRepository.InsertNewJobExecutions(ctx, newJobsToExecute)
}

func GetScheduledAtTime(t string) time.Time {
	runTime, err := time.Parse("15:04:05", t)
	if err != nil {
		log.Fatal(err)
	}

	now := time.Now()
	year, month, day := now.Date()

	hour, min, sec := runTime.Clock()
	combined := time.Date(year, month, day, hour, min, sec, 0, runTime.Location())

	return combined
}

func (j *jobSchedulingService) RecoverExpiredLeases(ctx context.Context) error {
	return j.jobExecutionRepository.MarkExpiredLeasesAsRetry(ctx)
}
