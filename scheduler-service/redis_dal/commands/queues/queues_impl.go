package queues

import (
	"context"
	"errors"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisQueueCommands struct {
	rdb *redis.Client
}

func NewRedisQueueCommands(rdb *redis.Client) RedisQueueCommands {
	return &redisQueueCommands{
		rdb: rdb,
	}
}

// Enqueue's from the left
func (r *redisQueueCommands) LEnqueue(ctx context.Context, queueName string, jobId string) error {
	return r.rdb.LPush(ctx, queueName, jobId).Err()
}

// Dequeue's the from the right, here it will block the queue and then get the result to avoid race conditions
func (r *redisQueueCommands) RDequeue(ctx context.Context, queueName string) (string, error) {
	data, error := r.rdb.BRPop(ctx, 0*time.Second, queueName).Result()
	if error != nil {
		return "", error
	}
	if len(data) < 2 {
		return "", errors.New("empty response")
	}
	return data[1], nil
}
