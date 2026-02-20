package contracts

import (
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/database_constants"
)

type JobExecutionUpdate struct {
	Status     database_constants.JobExecutionStatus
	RetryCount *int
}

type JobExecutionCreationData struct {
	Status      string
	ScheduledAt time.Time
}
