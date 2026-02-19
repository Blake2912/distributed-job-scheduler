package worker

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/constants/worker_constants"
	"github.com/Blake2912/distributed-job-scheduler/common/contracts"
	"github.com/Blake2912/distributed-job-scheduler/worker-service/internal/client"
	"github.com/Blake2912/distributed-job-scheduler/worker-service/internal/executor"
)

type Worker struct {
	client   client.SchedulerClient
	registry executor.Registry
}

// Constructor
func New(client client.SchedulerClient, registry executor.Registry) *Worker {
	return &Worker{client: client, registry: registry}
}

func (w *Worker) Run(ctx context.Context) {
	log.Println("Worker started")

	backoff := 200 * time.Millisecond
	maxBackoff := 5 * time.Second

	for {
		select {
		case <-ctx.Done():
			log.Println("Worker shutting down")
			return
		default:
		}

		job, err := w.client.LeaseJob(ctx)

		if err != nil {
			log.Println("Failed to lease job:", err)
			backoff = nextBackoff(backoff, maxBackoff)
			sleepWithJitter(ctx, backoff)
			continue
		}

		if job == nil {
			backoff = nextBackoff(backoff, maxBackoff)
			sleepWithJitter(ctx, backoff)
			continue
		}

		// Reset backoff when job received
		backoff = 200 * time.Millisecond

		w.executeJob(ctx, job)
	}
}

func (w *Worker) executeJob(ctx context.Context, job *contracts.JobToExecute) {

	log.Printf("Executing job %s\n", job.JobID)

	//get executor for the job type
	exec, err := w.registry.Get(job.JobType)

	if err != nil {
		log.Println("Unknown job type:", job.JobType)

		reportErr := w.client.ReportCompletion(ctx, job.JobID, worker_constants.FAILED, err.Error(), false)
		if reportErr != nil {
			log.Println("Failed reporting failure:", reportErr)
		}
		return
	}

	err = exec.Execute(ctx, job.Payload)
	if err != nil {
		log.Printf("Job failed ID=%s Error=%v\n", job.JobID, err)

		isRetryable := executor.IsRetryable(err)

		// Scheduler decides retry logic
		_ = w.client.ReportCompletion(ctx, job.JobID, worker_constants.FAILED, err.Error(), isRetryable)
		return
	}

	log.Println("Job succeeded:", job.JobID)

	reportErr := w.client.ReportCompletion(ctx, job.JobID, worker_constants.SUCCESS, "", false)
	if reportErr != nil {
		log.Println("Failed reporting success:", reportErr)
	}
}

func nextBackoff(current, max time.Duration) time.Duration {
	next := current * 2
	if next > max {
		return max
	}
	return next
}

func sleepWithJitter(ctx context.Context, base time.Duration) {
	jitter := time.Duration(rand.Int63n(int64(base / 2)))
	wait := base + jitter

	select {
	case <-time.After(wait):
	case <-ctx.Done():
	}
}
