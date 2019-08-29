package ogame

import (
	"fmt"

	"github.com/dustin/go-humanize"
	"github.com/google/gxui/math"
)

// ResourcesDetails ...
type ResourcesDetails struct {
	Metal struct {
		Available         int
		StorageCapacity   int
		CurrentProduction int
		// DenCapacity       int
	}
	Crystal struct {
		Available         int
		StorageCapacity   int
		CurrentProduction int
		// DenCapacity       int
	}
	Deuterium struct {
		Available         int
		StorageCapacity   int
		CurrentProduction int
		// DenCapacity       int
	}
	Energy struct {
		Available         int
		CurrentProduction int
		Consumption       int
	}
	Darkmatter struct {
		Available int
		Purchased int
		Found     int
	}
}

// Available returns the resources available
func (r ResourcesDetails) Available() Resources {
	return Resources{
		Metal:      r.Metal.Available,
		Crystal:    r.Crystal.Available,
		Deuterium:  r.Deuterium.Available,
		Energy:     r.Energy.Available,
		Darkmatter: r.Darkmatter.Available,
	}
}

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

// Sub subtract v from r
func (r Resources) Sub(v Resources) Resources {
	return Resources{
		Metal:     math.Max(r.Metal-v.Metal, 0),
		Crystal:   math.Max(r.Crystal-v.Crystal, 0),
		Deuterium: math.Max(r.Deuterium-v.Deuterium, 0),
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

// Div finds how many price a res can afford
func (r Resources) Div(price Resources) int {
	nb := math.MaxInt
	if price.Metal > 0 {
		nb = r.Metal / price.Metal
	}
	if price.Crystal > 0 {
		nb = math.Min(r.Crystal/price.Crystal, nb)
	}
	if price.Deuterium > 0 {
		nb = math.Min(r.Deuterium/price.Deuterium, nb)
	}
	return int(nb)
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
