package ogame

import (
	"math"
)

type spaceDock struct {
	BaseBuilding
}

func newSpaceDock() *spaceDock {
	b := new(spaceDock)
	b.Name = "space dock"
	b.ID = SpaceDockID
	b.IncreaseFactor = 5
	b.BaseCost = Resources{Metal: 200, Deuterium: 50, Energy: 50}
	b.Requirements = map[ID]int64{ShipyardID: 2}
	return b
}

// GetPrice returns the price to build the given level
func (b spaceDock) GetPrice(level int64, lfBonuses LfBonuses) Resources {
	tmp := func(baseCost int64, increaseFactor float64, level int64) int64 {
		return int64(float64(baseCost) * math.Pow(increaseFactor, float64(level-1)))
	}
	price := Resources{
		Metal:     tmp(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:   tmp(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium: tmp(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:    tmp(b.BaseCost.Energy, 2.5, level),
	}
	bonus := lfBonuses.CostTimeBonuses[b.ID].Cost
	return price.SubPercent(bonus)
}
