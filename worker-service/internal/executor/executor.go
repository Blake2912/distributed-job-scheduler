package executor

import (
	"context"
)

type Executor interface {
	Execute(ctx context.Context, payload string) error
}
