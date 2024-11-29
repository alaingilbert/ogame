package ogame

import (
	"math"
)

// BaseLevelable base struct for levelable (buildings, technologies)
type BaseLevelable struct {
	Base
	BaseCost       Resources
	IncreaseFactor float64
}

// GetPrice returns the price to build the given level
func (b BaseLevelable) GetPrice(level int64, lfBonuses LfBonuses) Resources {
	tmp := func(baseCost int64, increaseFactor float64, level int64) int64 {
		return int64(float64(baseCost) * math.Pow(increaseFactor, float64(level-1)))
	}
	price := Resources{
		Metal:     tmp(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:   tmp(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium: tmp(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:    tmp(b.BaseCost.Energy, b.IncreaseFactor, level),
	}
	bonus := lfBonuses.CostTimeBonuses[b.ID].Cost
	return price.SubPercent(bonus)
}
