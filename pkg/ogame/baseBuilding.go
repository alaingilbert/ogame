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
func (b BaseBuilding) DeconstructionPrice(level int64, techs IResearches) Resources {
	tmp := func(baseCost int64, increaseFactor float64, level int64) int64 {
		return int64(math.Floor(float64(baseCost)*math.Pow(increaseFactor, float64(level-1))) * (1 - 0.04*float64(techs.GetIonTechnology())))
	}
	return Resources{
		Metal:      tmp(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:    tmp(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium:  tmp(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:     tmp(b.BaseCost.Energy, b.IncreaseFactor, level),
		Population: tmp(b.BaseCost.Population, b.IncreaseFactor, level),
	}
}

func (b BaseBuilding) BuildingConstructionTime(level, universeSpeed int64, acc BuildingAccelerators, lfBonuses LfBonuses) time.Duration {
	price := b.GetPrice(level, lfBonuses)
	metalCost := float64(price.Metal)
	crystalCost := float64(price.Crystal)
	roboticLvl := float64(acc.GetRoboticsFactory())
	naniteLvl := float64(acc.GetNaniteFactory())
	hours := (metalCost + crystalCost) / (2500 * (1 + roboticLvl) * float64(universeSpeed) * math.Pow(2, naniteLvl))
	secs := hours * 3600
	if b.ID != NaniteFactoryID && (level-1) < 5 {
		secs = secs * (2 / (7 - (float64(level) - 1)))
	}
	secs = math.Max(1, secs)
	dur := time.Duration(int64(math.Floor(secs))) * time.Second
	bonus := lfBonuses.CostTimeBonuses[b.ID].Duration
	return time.Duration(float64(dur) - float64(dur)*bonus)
}

// ConstructionTime returns the duration it takes to build given level. Deconstruction time is the same function.
func (b BaseBuilding) ConstructionTime(level, universeSpeed int64, facilities BuildAccelerators, lfBonuses LfBonuses, _ CharacterClass, _ bool) time.Duration {
	return b.BuildingConstructionTime(level, universeSpeed, facilities, lfBonuses)
}

// GetLevel returns current level of a building
func (b BaseBuilding) GetLevel(resourcesBuildings IResourcesBuildings, facilities IFacilities, _ IResearches) int64 {
	if b.ID.IsResourceBuilding() {
		return resourcesBuildings.ByID(b.ID)
	} else if b.ID.IsFacility() {
		return facilities.ByID(b.ID)
	}
	return 0
}
