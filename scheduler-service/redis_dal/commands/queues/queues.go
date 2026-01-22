package queues

import (
	"context"
)

type RedisQueueCommands interface {
	LEnqueue(ctx context.Context, queueName string, jobId string) error
	RDequeue(ctx context.Context, queueName string) (string, error)
}
