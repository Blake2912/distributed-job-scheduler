package scheduler

import (
	"context"
	"log"
	"time"
)

// Placeholder scheduler logic for now
type Scheduler struct{}

func New() *Scheduler {
	return &Scheduler{}
}

func (s *Scheduler) Run(ctx context.Context) {
	log.Println("Scheduler started (leader)")

	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			log.Println("Scheduler stopped (lost leadership)")
			return
		case <-ticker.C:
			log.Println("Polling DB for eligible jobs...")
		}
	}
}
