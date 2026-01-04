package redisclient

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

func New(ctx context.Context) (*redis.Client, error) {
	addr, found := os.LookupEnv("REDIS_ADDR")
	if !found {
		addr = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return rdb, nil
}
