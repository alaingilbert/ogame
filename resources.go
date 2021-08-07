package ogame

import (
	"fmt"
	stdmath "math"

	humanize "github.com/dustin/go-humanize"
	"github.com/google/gxui/math"
)

// ResourcesDetails ...
type ResourcesDetails struct {
	Metal struct {
		Available         int64
		StorageCapacity   int64
		CurrentProduction int64
		BaseProduction    float64 // Production per Second
		// DenCapacity       int
	}
	Crystal struct {
		Available         int64
		StorageCapacity   int64
		CurrentProduction int64
		BaseProduction    float64 // Production per Second
		// DenCapacity       int
	}
	Deuterium struct {
		Available         int64
		StorageCapacity   int64
		CurrentProduction int64
		BaseProduction    float64 // Production per Second
		// DenCapacity       int
	}
	Energy struct {
		Available         int64
		CurrentProduction int64
		Consumption       int64
	}
	Darkmatter struct {
		Available int64
		Purchased int64
		Found     int64
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

// Production returns the resources currently Produced
func (r ResourcesDetails) Production() Resources {
	return Resources{
		Metal:      r.Metal.CurrentProduction,
		Crystal:    r.Crystal.CurrentProduction,
		Deuterium:  r.Deuterium.CurrentProduction,
		Energy:     r.Energy.CurrentProduction,
		Darkmatter: 0,
	}
}

// Storage returns the resources that can be stored
func (r ResourcesDetails) Storage() Resources {
	return Resources{
		Metal:      r.Metal.StorageCapacity,
		Crystal:    r.Crystal.StorageCapacity,
		Deuterium:  r.Deuterium.StorageCapacity,
		Energy:     0,
		Darkmatter: 0,
	}
}

// Resources represent ogame resources
type Resources struct {
	Metal      int64
	Crystal    int64
	Deuterium  int64
	Energy     int64
	Darkmatter int64
}

func (r Resources) String() string {
	return fmt.Sprintf("[%s|%s|%s]",
		humanize.Comma(r.Metal), humanize.Comma(r.Crystal), humanize.Comma(r.Deuterium))
}

// Total returns the sum of resources
func (r Resources) Total() int64 {
	return r.Deuterium + r.Crystal + r.Metal
}

// Value returns normalized total value of all resources
func (r Resources) Value() int64 {
	return r.Deuterium*3 + r.Crystal*2 + r.Metal
}

// Sub subtract v from r
func (r Resources) Sub(v Resources) Resources {
	return Resources{
		Metal:     max64(r.Metal-v.Metal, 0),
		Crystal:   max64(r.Crystal-v.Crystal, 0),
		Deuterium: max64(r.Deuterium-v.Deuterium, 0),
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
func (r Resources) Mul(scalar int64) Resources {
	return Resources{
		Metal:     r.Metal * scalar,
		Crystal:   r.Crystal * scalar,
		Deuterium: r.Deuterium * scalar,
	}
}

func min64(values ...int64) int64 {
	m := int64(math.MaxInt)
	for _, v := range values {
		if v < m {
			m = v
		}
	}
	return m
}

func max64(values ...int64) int64 {
	m := int64(math.MinInt)
	for _, v := range values {
		if v > m {
			m = v
		}
	}
	return m
}

// Div finds how many price a res can afford
func (r Resources) Div(price Resources) int64 {
	nb := int64(math.MaxInt)
	if price.Metal > 0 {
		nb = r.Metal / price.Metal
	}
	if price.Crystal > 0 {
		nb = min64(r.Crystal/price.Crystal, nb)
	}
	if price.Deuterium > 0 {
		nb = min64(r.Deuterium/price.Deuterium, nb)
	}
	return nb
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

// FitsIn get the number of ships required to transport the resource
func (r Resources) FitsIn(ship Ship, techs Researches, probeRaids, isCollector, isPioneers bool) int64 {
	cargo := ship.GetCargoCapacity(techs, probeRaids, isCollector, isPioneers)
	if cargo == 0 {
		return 0
	}
	return int64(stdmath.Ceil(float64(r.Total()) / float64(cargo)))
}

// SubReal subtract v from r
func (r Resources) SubReal(v Resources) Resources {
	return Resources{
		Metal:     r.Metal - v.Metal,
		Crystal:   r.Crystal - v.Crystal,
		Deuterium: r.Deuterium - v.Deuterium,
	}
}
