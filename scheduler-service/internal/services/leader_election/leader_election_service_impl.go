package leader_election

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

type LeaderElector struct {
	rdb        *redis.Client
	key        string
	instanceID string
	ttl        time.Duration
	renewEvery time.Duration
}

func New(rdb *redis.Client, key string, ttl time.Duration) LeaderElection {
	return &LeaderElector{
		rdb:        rdb,
		key:        key,
		instanceID: uuid.NewString(),
		ttl:        ttl,
		renewEvery: ttl / 3,
	}
}

func (leadership *LeaderElector) Run(
	ctx context.Context,
	onLeader func(context.Context),
) error {
	ticker := time.NewTicker(leadership.renewEvery)
	defer ticker.Stop()

	var (
		isLeader  bool
		leaderCtx context.Context
		cancel    context.CancelFunc
	)

	for {
		select {
		case <-ctx.Done():
			if cancel != nil {
				cancel()
			}
			return nil

		case <-ticker.C:
			//if instance is not leader, try to acquire leadership
			if !isLeader {
				ok, err := leadership.TryAcquire(ctx)
				if err != nil {
					//if failed, retry on next tick
					continue
				}
				//if succeeded, update info
				if ok {
					isLeader = true

					//create child context, lives while instance is leader, cancelled when leadership lost
					leaderCtx, cancel = context.WithCancel(ctx)

					//goroutine for scheduler logic
					//this is non blocking so that the leader election loop keeps running
					go onLeader(leaderCtx)
				}
				continue
			}

			//if instance is leader, try to renew leadership
			ok, err := leadership.Renew(ctx)
			if err != nil || !ok {
				isLeader = false

				//stop all leader related stuff
				cancel()
				cancel = nil
			}
		}
	}
}

// Try acquiring leadership by setting key in redis (SET NX EX)
func (leadership *LeaderElector) TryAcquire(ctx context.Context) (bool, error) {
	ok, err := leadership.rdb.SetNX(ctx, leadership.key, leadership.instanceID, leadership.ttl).Result()
	return ok, err
}

// Renewing leadership
func (leadership *LeaderElector) Renew(ctx context.Context) (bool, error) {
	var renewScript = redis.NewScript(`
		if redis.call("GET", KEYS[1]) == ARGV[1] then
			return redis.call("PEXPIRE", KEYS[1], ARGV[2])
		else
			return 0
		end
		`)

	res, err := renewScript.Run(
		ctx,
		leadership.rdb,
		[]string{leadership.key},
		leadership.instanceID,
		leadership.ttl.Milliseconds(),
	).Int()

	return res == 1, err
}
