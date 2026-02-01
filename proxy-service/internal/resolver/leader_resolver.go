package resolver

import "context"

type LeaderResolver interface {
	GetLeader(ctx context.Context) (string, error)
}
