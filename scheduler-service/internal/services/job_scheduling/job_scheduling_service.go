package jobscheduling

import "context"

type JobSchedulingService interface {
	ScheduleJobs(ctx context.Context) error
}
