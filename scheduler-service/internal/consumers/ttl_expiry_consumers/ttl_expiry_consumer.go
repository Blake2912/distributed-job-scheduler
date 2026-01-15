package ttlexpiryconsumers

import (
	"context"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/eventbus"
)

type TTLExpiryConsumer interface {
	StartTTLExpiryExecution(ctx context.Context, bus *eventbus.EventBus[eventbus.TTLExpiredEvent])
}
