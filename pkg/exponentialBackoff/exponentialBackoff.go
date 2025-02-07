package exponentialBackoff

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
	max   uint32
}

// New ...
func New(ctx context.Context, max int) *ExponentialBackoff {
	return NewWithClock(ctx, clockwork.NewRealClock(), max)
}

// NewWithClock ...
func NewWithClock(ctx context.Context, clock clockwork.Clock, maxSleep int) *ExponentialBackoff {
	maxSleep = max(maxSleep, 0)
	e := new(ExponentialBackoff)
	e.ctx = ctx
	e.clock = clock
	e.max = uint32(max(maxSleep, 0))
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
		if e.ctx.Err() != nil {
			return
		}
	}
}

func (e *ExponentialBackoff) Iterator() func(func() bool) {
	return func(yield func() bool) {
		e.LoopForever(yield)
	}
}

// Wait ...
func (e *ExponentialBackoff) Wait() {
	currVal := atomic.LoadUint32(&e.val)
	newVal := uint32(1)
	if currVal > 0 {
		newVal = currVal * 2
		if e.max > 0 {
			newVal = min(newVal, e.max)
		}
	}
	atomic.StoreUint32(&e.val, newVal)
	select {
	case <-e.clock.After(time.Duration(newVal) * time.Second):
	case <-e.ctx.Done():
	}
}

// Reset ...
func (e *ExponentialBackoff) Reset() {
	atomic.StoreUint32(&e.val, 0)
}
