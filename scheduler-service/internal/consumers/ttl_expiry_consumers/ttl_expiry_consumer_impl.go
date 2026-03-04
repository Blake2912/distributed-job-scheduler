package ttlexpiryconsumers

import (
	"context"
	"log"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/eventbus"
	workerexpirationservice "github.com/Blake2912/distributed-job-scheduler/scheduler-service/internal/services/worker_expiration_service"
)

type ttlExpiryConsumer struct {
	expirationSvc workerexpirationservice.WorkerExpiration
}

func NewTTLExpiryConsumer(expirationSvc workerexpirationservice.WorkerExpiration) TTLExpiryConsumer {
	return &ttlExpiryConsumer{
		expirationSvc: expirationSvc,
	}
}

func (t *ttlExpiryConsumer) StartTTLExpiryExecution(ctx context.Context, bus *eventbus.EventBus[eventbus.TTLExpiredEvent]) {

	id, ch := bus.Subscribe()
	go func() {
		defer bus.Unsubscribe(id)

		for {
			select {
			case <-ctx.Done():
				return
			case event := <-ch:
				log.Println("Executor received:", event.Key)
				t.expirationSvc.HandleWorkerExpiry(ctx, event.Key)
			}
		}
	}()

}
