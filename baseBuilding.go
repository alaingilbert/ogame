package ogame

import (
	"math"
	"time"
)

// BaseBuilding base struct for buildings
type BaseBuilding struct {
	BaseLevelable
}

// DeconstructionPrice returns the price to tear down to given level
func (b BaseBuilding) DeconstructionPrice(level int64, techs Researches) Resources {
	tmp := func(baseCost int64, increaseFactor float64, level int64) int64 {
		return int64(math.Floor(float64(baseCost)*math.Pow(increaseFactor, float64(level-1))) * (1 - 0.04*float64(techs.IonTechnology)))
	}
	return Resources{
		Metal:     tmp(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:   tmp(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium: tmp(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:    tmp(b.BaseCost.Energy, b.IncreaseFactor, level),
	}
}

// ConstructionTime returns the duration it takes to build given level. Deconstruction time is the same function.
func (b BaseBuilding) ConstructionTime(level, universeSpeed int64, facilities Facilities, hasTechnocrat, isDiscoverer bool) time.Duration {
	price := b.GetPrice(level)
	metalCost := float64(price.Metal)
	crystalCost := float64(price.Crystal)
	roboticLvl := float64(facilities.RoboticsFactory)
	naniteLvl := float64(facilities.NaniteFactory)
	hours := (metalCost + crystalCost) / (2500 * (1 + roboticLvl) * float64(universeSpeed) * math.Pow(2, naniteLvl))
	secs := hours * 3600
	if b.ID != NaniteFactoryID && (level-1) < 5 {
		secs = secs * (2 / (7 - (float64(level) - 1)))
	}
	secs = math.Max(1, secs)
	return time.Duration(int64(math.Floor(secs))) * time.Second
}

// GetLevel returns current level of a building
func (b BaseBuilding) GetLevel(lazyResourcesBuildings LazyResourcesBuildings, lazyFacilities LazyFacilities, _ LazyResearches) int64 {
	if b.ID.IsResourceBuilding() {
		return lazyResourcesBuildings().ByID(b.ID)
	} else if b.ID.IsFacility() {
		return lazyFacilities().ByID(b.ID)
	}
	return 0
}
