package jobsservice

import (
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/api/contracts/jobs"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/models"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/sql_dal/repository"
)

type jobsService struct {
	jobRepository repository.JobRepository
}

func NewJobsService(jobRepository repository.JobRepository) JobsService {
	return &jobsService{
		jobRepository: jobRepository,
	}
}

// Create Jobs
func (j *jobsService) CreateJobs(ctx context.Context, payload []jobs.CreateJobsPayload) error {

	jobsToCreate := make([]models.Jobs, 0, len(payload))

	for _, jobFromPayload := range payload {
		metadata, err := prepareMetaDataForInsert(*jobFromPayload.IsRecurringJob)

		if err != nil {
			return err
		}

		job := models.Jobs{
			Name:                    jobFromPayload.Name,
			Type:                    jobFromPayload.Type,
			Config:                  jobFromPayload.Config,
			NextRunAt:               toNullTime(jobFromPayload.TriggerTime.UTC()),
			ShouldRetryAfterBackoff: *jobFromPayload.ShouldRetryAfterBackoff,
			Metadata:                metadata,
			Enabled:                 true,
		}

		jobsToCreate = append(jobsToCreate, job)
	}

	if len(jobsToCreate) == 0 {
		log.Println("No jobs found to create skipping create")
		return nil
	}

	return j.jobRepository.CreateJobs(ctx, jobsToCreate)
}

func prepareMetaDataForInsert(isRecurringJob bool) (string, error) {
	metaDataInContract := contracts.JobMetaDataContract{
		IsRecurringJob: isRecurringJob,
	}

	b, err := json.Marshal(metaDataInContract)
	inString := string(b)

	return inString, err
}

func toNullTime(t time.Time) sql.NullTime {
	if t.IsZero() {
		return sql.NullTime{Valid: false}
	}

	// strip date, keep only time
	timeOnly := time.Date(0, 1, 1,
		t.Hour(),
		t.Minute(),
		t.Second(),
		t.Nanosecond(),
		time.UTC,
	)

	return sql.NullTime{
		Time:  timeOnly,
		Valid: true,
	}
}

// View jobs
