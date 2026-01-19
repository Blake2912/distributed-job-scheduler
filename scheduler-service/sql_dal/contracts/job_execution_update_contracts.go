package contracts

import "github.com/Blake2912/distributed-job-scheduler/common/database_constants"

type JobExecutionUpdate struct {
	Status     database_constants.JobExecutionStatus
	RetryCount *int
}
