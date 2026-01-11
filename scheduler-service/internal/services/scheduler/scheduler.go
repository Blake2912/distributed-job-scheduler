package scheduler

import (
	"context"
	"log"
	"time"

	jobscheduling "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/job_scheduling"
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

	ticker := time.NewTicker(30 * time.Minute) // Poll once in every 30 min
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped (lost leadership)")
			return
		case <-ticker.C:
			log.Println("Polling DB for eligible jobs...")
			s.jobSchedulingService.ScheduleJobs(ctx)
		}
	}
}
