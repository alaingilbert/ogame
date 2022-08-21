package wrapper

import (
	"context"
	"time"
)

// ExponentialBackoff ...
type ExponentialBackoff struct {
	ctx context.Context
	val int
	max int
}

// NewExponentialBackoff ...
func NewExponentialBackoff(ctx context.Context, max int) *ExponentialBackoff {
	if max < 0 {
		max = 0
	}
	e := new(ExponentialBackoff)
	e.ctx = ctx
	e.max = max
	return e
}

// Wait ...
func (e *ExponentialBackoff) Wait() {
	if e.val == 0 {
		e.val = 1
	} else {
		select {
		case <-time.After(time.Duration(e.val) * time.Second):
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
