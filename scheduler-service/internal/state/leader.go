package state

import (
	"sync/atomic"
)

var isLeader atomic.Bool

func SetLeader(val bool) {
	isLeader.Store(val)
}

func IsLeader() bool {
	return isLeader.Load()
}
