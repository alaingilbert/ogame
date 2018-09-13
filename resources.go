package ogame

import (
	"fmt"

	"github.com/dustin/go-humanize"
)

// Resources represent ogame resources
type Resources struct {
	Metal      int
	Crystal    int
	Deuterium  int
	Energy     int
	Darkmatter int
}

func (r Resources) String() string {
	return fmt.Sprintf("[%s|%s|%s]",
		humanize.Comma(int64(r.Metal)), humanize.Comma(int64(r.Crystal)), humanize.Comma(int64(r.Deuterium)))
}

// Total returns the sum of resources
func (r Resources) Total() int {
	return r.Deuterium + r.Crystal + r.Metal
}

// Value returns normalized total value of all resources
func (r Resources) Value() int {
	return r.Deuterium*3 + r.Crystal*2 + r.Metal
}

// Sub substract v from r
func (r Resources) Sub(v Resources) Resources {
	return Resources{
		Metal:     r.Metal - v.Metal,
		Crystal:   r.Crystal - v.Crystal,
		Deuterium: r.Deuterium - v.Deuterium,
	}
}

// Add adds two resources together
func (r Resources) Add(v Resources) Resources {
	return Resources{
		Metal:     r.Metal + v.Metal,
		Crystal:   r.Crystal + v.Crystal,
		Deuterium: r.Deuterium + v.Deuterium,
	}
}

// Mul multiply resources with scalar.
func (r Resources) Mul(scalar int) Resources {
	return Resources{
		Metal:     r.Metal * scalar,
		Crystal:   r.Crystal * scalar,
		Deuterium: r.Deuterium * scalar,
	}
}

// CanAfford alias to Gte
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
