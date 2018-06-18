package ogame

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

// Resources ...
type Resources struct {
	Metal      int
	Crystal    int
	Deuterium  int
	Energy     int
	Darkmatter int
}

// String ...
func (r Resources) String() string {
	return fmt.Sprintf("[%s|%s|%s]",
		humanize.Comma(int64(r.Metal)), humanize.Comma(int64(r.Crystal)), humanize.Comma(int64(r.Deuterium)))
}

// Total ...
func (r Resources) Total() int {
	return r.Deuterium + r.Crystal + r.Metal
}

// Value returns normalized total value of all resources
func (r Resources) Value() int {
	return r.Deuterium*3 + r.Crystal*2 + r.Metal
}

// Sub ...
func (r Resources) Sub(v Resources) Resources {
	return Resources{
		Metal:     r.Metal - v.Metal,
		Crystal:   r.Crystal - v.Crystal,
		Deuterium: r.Deuterium - v.Deuterium,
	}
}

// Add ...
func (r Resources) Add(v Resources) Resources {
	return Resources{
		Metal:     r.Metal + v.Metal,
		Crystal:   r.Crystal + v.Crystal,
		Deuterium: r.Deuterium + v.Deuterium,
	}
}

// CanAfford ...
func (r Resources) CanAfford(cost Resources) bool {
	return r.Gte(cost)
}

// Gte greater than or equal
func (r Resources) Gte(val Resources) bool {
	return r.Metal >= val.Metal &&
		r.Crystal >= val.Crystal &&
		r.Deuterium >= val.Deuterium
}

// Lte less than or equal
func (r Resources) Lte(val Resources) bool {
	return r.Metal <= val.Metal &&
		r.Crystal <= val.Crystal &&
		r.Deuterium <= val.Deuterium
}
