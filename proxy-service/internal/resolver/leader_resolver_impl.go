package resolver

import (
	"context"
	"sync/atomic"
	"time"

	"github.com/Blake2912/distributed-job-scheduler/common/infra_constants"
	"github.com/Blake2912/distributed-job-scheduler/scheduler-service/redis_dal/commands/queries"
)

type leaderResolver struct {
	redisQueries queries.Queries
	leader       atomic.Value
	last         atomic.Int64
}

func NewLeaderResovler(redisQueries queries.Queries) LeaderResolver {
	r := &leaderResolver{
		redisQueries: redisQueries,
	}
	r.last.Store(0)

	return r
}

func (r *leaderResolver) GetLeader(ctx context.Context) (string, error) {
	if time.Now().Unix()-r.last.Load() < 2 {
		if l := r.leader.Load(); l != nil {
			return l.(string), nil
		}
	}
	leaderAddress, err := r.redisQueries.GetValue(ctx, infra_constants.LeaderAddressKey)

	if err == nil && leaderAddress != "" {
		r.leader.Store(leaderAddress)
		r.last.Store(time.Now().Unix())
	}
	return leaderAddress, err
}
