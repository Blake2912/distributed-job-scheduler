package ttlexpiryconsumers

import (
	"context"
	"log"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/eventbus"
)

type ttlExpiryConsumer struct {
}

func NewTTLExpiryConsumer() TTLExpiryConsumer {
	return &ttlExpiryConsumer{}
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
			}
		}
	}()

}
