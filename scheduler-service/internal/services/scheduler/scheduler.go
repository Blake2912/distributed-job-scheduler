package scheduler

import (
	"context"
	"log"
	"time"

	jobscheduling "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/job_scheduling"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/state"
)

// Placeholder scheduler logic for now
type Scheduler struct {
	jobSchedulingService jobscheduling.JobSchedulingService
}

func New(jobSchedulingService jobscheduling.JobSchedulingService) *Scheduler {
	return &Scheduler{
		jobSchedulingService: jobSchedulingService,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	log.Println("Scheduler started (leader)")
	state.SetLeader(true)

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			state.SetLeader(false)
			log.Println("Scheduler stopped (lost leadership)")
			return
		case <-ticker.C:
			log.Println("Polling DB for eligible jobs...")
			s.jobSchedulingService.ScheduleJobs(ctx)
		}
	}
}
