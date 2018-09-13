package ogame

import (
	"math"
)

type BaseLevelable struct {
	Base
	BaseCost       Resources
	IncreaseFactor float64
}

// GetPrice returns the price to build the given level
func (b BaseLevelable) GetPrice(level int) Resources {
	tmp := func(baseCost int, increaseFactor float64, level int) int {
		return int(float64(baseCost) * math.Pow(increaseFactor, float64(level-1)))
	}
	return Resources{
		Metal:     tmp(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:   tmp(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium: tmp(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:    tmp(b.BaseCost.Energy, b.IncreaseFactor, level),
	}
}
