package ogame

import "math"

type astrophysics struct {
	BaseTechnology
}

func newAstrophysics() *astrophysics {
	b := new(astrophysics)
	b.Name = "astrophysics"
	b.ID = AstrophysicsID
	b.IncreaseFactor = 1.75
	b.BaseCost = Resources{Metal: 4000, Crystal: 8000, Deuterium: 4000}
	b.Requirements = map[ID]int64{EspionageTechnologyID: 4, ImpulseDriveID: 3}
	return b
}

// GetPrice returns the price to build the given level
func (b astrophysics) GetPrice(level int64) Resources {
	tmp := func(baseCost int64, increaseFactor float64, level int64) int64 {
		return int64(math.Round(float64(baseCost)*math.Pow(increaseFactor, float64(level-1))/100) * 100)
	}
	return Resources{
		Metal:     tmp(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:   tmp(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium: tmp(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:    tmp(b.BaseCost.Energy, b.IncreaseFactor, level),
	}
}
