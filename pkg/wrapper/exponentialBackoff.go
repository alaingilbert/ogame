package wrapper

import (
	"context"
	"github.com/alaingilbert/clockwork"
	"sync/atomic"
	"time"
)

// ExponentialBackoff ...
type ExponentialBackoff struct {
	ctx   context.Context
	clock clockwork.Clock
	val   uint32 // atomic
	max   int
}

// NewExponentialBackoff ...
func NewExponentialBackoff(ctx context.Context, clock clockwork.Clock, max int) *ExponentialBackoff {
	if max < 0 {
		max = 0
	}
	e := new(ExponentialBackoff)
	e.ctx = ctx
	e.clock = clock
	e.max = max
	return e
}

// LoopForever execute the callback with exponential backoff
// The callback return true if we should continue retrying
// or false if we should stop and exit.
func (e *ExponentialBackoff) LoopForever(clb func() bool) {
	for {
		keepLooping := clb()
		if !keepLooping {
			return
		}
		e.Wait()
		select {
		case <-e.ctx.Done():
			return
		default:
		}
	}
}

// Wait ...
func (e *ExponentialBackoff) Wait() {
	currVal := atomic.LoadUint32(&e.val)
	if currVal == 0 {
		atomic.StoreUint32(&e.val, 1)
		return
	}

	newVal := currVal * 2
	if e.max > 0 && newVal > uint32(e.max) {
		newVal = uint32(e.max)
	}
	atomic.StoreUint32(&e.val, newVal)
	select {
	case <-e.clock.After(time.Duration(currVal) * time.Second):
	case <-e.ctx.Done():
		return
	}
}

// Reset ...
func (e *ExponentialBackoff) Reset() {
	atomic.StoreUint32(&e.val, 0)
}
