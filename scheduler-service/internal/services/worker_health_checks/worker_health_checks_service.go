package workerhealthchecks

import "context"

type WorkerHealthChecks interface {
	CheckHealth(ctx context.Context, workerId string, jobExecutionId string) error
	DeleteWorkerKeys(ctx context.Context, workerId string, jobExecutionId string) error
}
