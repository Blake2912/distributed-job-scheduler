package scheduler

import (
	"context"
	"log"
	"log/slog"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/infra_constants"
	schedulerconfig "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/config"
	jobscheduling "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/job_scheduling"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/state"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/commands/queries"
)

// Placeholder scheduler logic for now
type Scheduler struct {
	jobSchedulingService jobscheduling.JobSchedulingService
	redisQueries         queries.Queries
}

func New(jobSchedulingService jobscheduling.JobSchedulingService, redisQueries queries.Queries) *Scheduler {
	return &Scheduler{
		jobSchedulingService: jobSchedulingService,
		redisQueries:         redisQueries,
	}
}

func (s *Scheduler) Run(ctx context.Context) {
	log.Println("Scheduler started (leader)")
	state.SetLeader(true)

	ttl := time.Second * 35
	address := schedulerconfig.GetAdvertisedLeaderAddress()

	err := s.redisQueries.CheckAndSetKeyWithTTL(ctx, infra_constants.LeaderAddressKey, address, ttl)
	if err != nil {
		slog.Error("An error occurred while setting the leader address key to redis, leader will be unreachable | %s", "err", err.Error())
	}

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			state.SetLeader(false)
			cleanUpContext, cancel := context.WithTimeout(context.Background(), 3*time.Second)
			defer cancel()
			_, err := s.redisQueries.DeleteKey(cleanUpContext, infra_constants.LeaderAddressKey)

			if err != nil {
				slog.Error("An error occurred while deleting the key for leader address %s", "err", err.Error())
			}

			log.Println("Scheduler stopped (lost leadership)")
			return
		case <-ticker.C:
			log.Println("Polling DB for eligible jobs...")
			err := s.redisQueries.ResetTTL(ctx, infra_constants.LeaderAddressKey, ttl)
			if err != nil {
				slog.Error("An error occurred while resetting ttl key for leader address", "err", err.Error())
			}
			s.jobSchedulingService.ScheduleJobs(ctx)
		}
	}
}

func (s *Scheduler) StartExpiryRecoveryLoop(ctx context.Context) {

	ticker := time.NewTicker(60 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			s.jobSchedulingService.RecoverExpiredLeases(ctx)
		}
	}
}
