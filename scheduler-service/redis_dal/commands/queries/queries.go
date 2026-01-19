package queries

import (
	"context"
	"time"
)

type Queries interface {
	CheckAndSetKeyWithTTL(ctx context.Context, key string, value string, ttl time.Duration) error
	DeleteKey(ctx context.Context, key string) (int64, error)
	KeyExists(ctx context.Context, key string) (bool, error)
	GetValue(ctx context.Context, key string) (string, error)
}
