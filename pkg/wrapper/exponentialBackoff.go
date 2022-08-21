package wrapper

import (
	"context"
	"github.com/alaingilbert/clockwork"
	"time"
)

// ExponentialBackoff ...
type ExponentialBackoff struct {
	ctx   context.Context
	clock clockwork.Clock
	val   int
	max   int
}

// NewExponentialBackoff ...
func NewExponentialBackoff(ctx context.Context, max int) *ExponentialBackoff {
	if max < 0 {
		max = 0
	}
	e := new(ExponentialBackoff)
	e.ctx = ctx
	e.clock = clockwork.NewRealClock()
	e.max = max
	return e
}

func (e *ExponentialBackoff) SetClock(clock clockwork.Clock) {
	e.clock = clock
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
	if e.val == 0 {
		e.val = 1
	} else {
		select {
		case <-e.clock.After(time.Duration(e.val) * time.Second):
		case <-e.ctx.Done():
			return
		}
		e.val *= 2
		if e.max > 0 {
			if e.val > e.max {
				e.val = e.max
			}
		}
	}
}

// Reset ...
func (e *ExponentialBackoff) Reset() {
	e.val = 0
}
