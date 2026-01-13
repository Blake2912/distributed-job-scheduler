package redissubscriber

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/eventbus"
	"github.com/redis/go-redis/v9"
)

func PublishRedisKeyExpiryEvent(rdb *redis.Client, ctx context.Context, bus *eventbus.EventBus[eventbus.TTLExpiredEvent]) {
	go func() {
		pubsub := rdb.Subscribe(ctx, "__keyevent@0__:expired")
		defer pubsub.Close()

		log.Println("TTL expiry listener started")

		for {
			select {
			case <-ctx.Done():
				log.Println("TTL expiry listener stopped")
				return
			default:
				msg, err := pubsub.ReceiveMessage(ctx)
				if err != nil {
					if ctx.Err() != nil {
						return
					}
					log.Println("Redis subscribe error:", err)
					time.Sleep(time.Second) // avoid tight loop on failure
					continue
				}
				fmt.Printf("Key expired: %s\n", msg.Payload)

				// Handle expiry event
				bus.Publish(eventbus.TTLExpiredEvent{
					Key: msg.Payload,
				})
			}
		}
	}()
}
