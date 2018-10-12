package ogame

import "time"

// ExponentialBackoff ...
type ExponentialBackoff struct {
	val int
	max int
}

// NewExponentialBackoff ...
func NewExponentialBackoff(max int) *ExponentialBackoff {
	if max < 0 {
		max = 0
	}
	e := new(ExponentialBackoff)
	e.max = max
	return e
}

// Wait ...
func (e *ExponentialBackoff) Wait() {
	if e.val == 0 {
		e.val = 1
	} else {
		time.Sleep(time.Duration(e.val) * time.Second)
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
