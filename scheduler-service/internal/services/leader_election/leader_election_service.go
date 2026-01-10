package leader_election

import "context"

type LeaderElection interface {
	// Run blocks and manages leadership.
	// It calls onStartLeadership when leadership is acquired
	// and cancels the context when leadership is lost.
	Run(ctx context.Context, onStartLeadership func(ctx context.Context)) error
}
