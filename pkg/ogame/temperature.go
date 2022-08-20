package ogame

import "math"

// Temperature planet temperature values
type Temperature struct {
	Min int64
	Max int64
}

// Mean returns the planet mean temperature
func (t Temperature) Mean() int64 {
	return int64(math.Round(float64(t.Min+t.Max) / 2))
}
