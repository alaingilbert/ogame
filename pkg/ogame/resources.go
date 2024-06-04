package ogame

import (
	"fmt"
	"github.com/alaingilbert/ogame/pkg/utils"
	humanize "github.com/dustin/go-humanize"
	"math"
)

// ResourcesDetails ...
type ResourcesDetails struct {
	Metal struct {
		Available         int64
		StorageCapacity   int64
		CurrentProduction int64
		// DenCapacity       int
	}
	Crystal struct {
		Available         int64
		StorageCapacity   int64
		CurrentProduction int64
		// DenCapacity       int
	}
	Deuterium struct {
		Available         int64
		StorageCapacity   int64
		CurrentProduction int64
		// DenCapacity       int
	}
	Food struct {
		Available           int64
		StorageCapacity     int64
		Overproduction      int64
		ConsumedIn          int64
		TimeTillFoodRunsOut int64
	}
	Population struct {
		Available   int64
		T2Lifeforms int64
		T3Lifeforms int64
		LivingSpace int64
		Satisfied   int64
		Hungry      float64
		GrowthRate  float64
		BunkerSpace int64
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
		Food:       r.Food.Available,
		Population: r.Population.Available,
		Energy:     r.Energy.Available,
		Darkmatter: r.Darkmatter.Available,
	}
}

// Resources represent ogame resources
type Resources struct {
	Metal      int64
	Crystal    int64
	Deuterium  int64
	Energy     int64
	Darkmatter int64
	Population int64
	Food       int64
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
		Metal:     utils.MaxInt(r.Metal-v.Metal, 0),
		Crystal:   utils.MaxInt(r.Crystal-v.Crystal, 0),
		Deuterium: utils.MaxInt(r.Deuterium-v.Deuterium, 0),
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

// Div finds how many price a res can afford
func (r Resources) Div(price Resources) int64 {
	nb := int64(math.MaxInt)
	if price.Metal > 0 {
		nb = r.Metal / price.Metal
	}
	if price.Crystal > 0 {
		nb = utils.MinInt(r.Crystal/price.Crystal, nb)
	}
	if price.Deuterium > 0 {
		nb = utils.MinInt(r.Deuterium/price.Deuterium, nb)
	}
	return nb
}

// CanAfford alias to Gte
func (r Resources) CanAfford(cost Resources) bool {
	return r.Gte(cost)
}

// CanAfford2 alias to Gte
func (r Resources) CanAfford2(cost Resources, population bool) bool {
	return r.Gte2(cost, population)
}

// Gte greater than or equal
func (r Resources) Gte(val Resources) bool {
	return r.Gte2(val, false)
}

// Gte2 greater than or equal
func (r Resources) Gte2(val Resources, population bool) bool {
	base := r.Metal >= val.Metal &&
		r.Crystal >= val.Crystal &&
		r.Deuterium >= val.Deuterium
	if population {
		base = base && r.Population >= val.Population
	}
	return base
}

// Lte less than or equal
func (r Resources) Lte(val Resources) bool {
	return r.Metal <= val.Metal &&
		r.Crystal <= val.Crystal &&
		r.Deuterium <= val.Deuterium
}

// FitsIn get the number of ships required to transport the resource
func (r Resources) FitsIn(ship Ship, techs Researches, bonus LfBonuses, characterClass CharacterClass, multiplier float64, probeRaids bool) int64 {
	cargo := ship.GetCargoCapacity(techs, bonus, characterClass, multiplier, probeRaids)
	if cargo == 0 {
		return 0
	}
	return int64(math.Ceil(float64(r.Total()) / float64(cargo)))
}

// SubPercent subtract the percentage from the initial values
func (r Resources) SubPercent(pct float64) Resources {
	return Resources{
		Metal:     r.Metal - int64(float64(r.Metal)*pct),
		Crystal:   r.Crystal - int64(float64(r.Crystal)*pct),
		Deuterium: r.Deuterium - int64(float64(r.Deuterium)*pct),
		Energy:    r.Energy - int64(float64(r.Energy)*pct),
	}
}
