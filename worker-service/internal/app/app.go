package app

import (
	"context"
	"log"
	"sync"

	"github.com/Blake2912/distributed-job-scheduler/worker-service/internal/worker"
)

type App struct {
	worker *worker.Worker
	count  int
}

func New(worker *worker.Worker, count int) *App {
	return &App{
		worker: worker,
		count:  count,
	}
}

func (a *App) Run(ctx context.Context) {

	log.Printf("Starting %d workers\n", a.count)

	var wg sync.WaitGroup

	for i := 0; i < a.count; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			log.Printf("Worker-%d started\n", id)
			a.worker.Run(ctx)
		}(i)
	}

	wg.Wait()
}
