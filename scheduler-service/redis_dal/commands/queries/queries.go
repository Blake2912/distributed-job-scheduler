package queries

import (
	"context"
	"time"
)

type Queries interface {
	CheckAndSetKeyWithTTL(ctx context.Context, key string, value string, ttl time.Duration) error
	DeleteKey(ctx context.Context, key string) (int64, error)
}
