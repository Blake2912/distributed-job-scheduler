package spawnworkers

import "context"


type SpawnWorkersService interface {
	SpawnWorkers(ctx context.Context, noOfWorkers int) error
}