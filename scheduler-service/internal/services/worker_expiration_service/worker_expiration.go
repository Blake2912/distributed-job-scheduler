package workerexpirationservice

import "context"

type WorkerExpiration interface {
	HandleWorkerExpiry(ctx context.Context, expiredKey string)
}
