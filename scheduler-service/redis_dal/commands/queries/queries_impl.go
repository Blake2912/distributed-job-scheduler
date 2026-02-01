package queries

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type queries struct {
	rdb *redis.Client
}

func NewRedisQueries(rdb *redis.Client) Queries {
	return &queries{
		rdb: rdb,
	}
}

// Checks if key exists if found then resets the TTL else creates a new key with the value
func (q *queries) CheckAndSetKeyWithTTL(ctx context.Context, key string, value string, ttl time.Duration) error {

	exists, err := q.KeyExists(ctx, key)
	if err != nil {
		return err
	}
	if !exists {
		return q.setKeyWithTTL(ctx, key, value, ttl)
	} else {
		return q.ResetTTL(ctx, key, ttl)
	}
}

// Deletes a key from redis
func (q *queries) DeleteKey(ctx context.Context, key string) (int64, error) {
	return q.rdb.Del(ctx, key).Result()
}

// Returns true if any key exists in the redis database
func (q *queries) KeyExists(ctx context.Context, key string) (bool, error) {
	res, err := q.rdb.Exists(ctx, key).Result()
	return res != 0, err
}

// Gets the value for a particular key
func (q *queries) GetValue(ctx context.Context, key string) (string, error) {
	return q.rdb.Get(ctx, key).Result()
}

// Sets a new key in redis with a specified TTL
func (q *queries) setKeyWithTTL(ctx context.Context, key string, value string, ttl time.Duration) error {
	return q.rdb.Set(ctx, key, value, ttl).Err()
}

// Resets TTL for a given Key in redis
func (q *queries) ResetTTL(ctx context.Context, key string, ttl time.Duration) error {
	return q.rdb.Expire(ctx, key, ttl).Err()
}
