package ogame

import (
	"math"
	"time"
)

// BaseBuilding base struct for buildings
type BaseBuilding struct {
	BaseLevelable
}

// ConstructionTime returns the duration it takes to build given level
func (b BaseBuilding) ConstructionTime(level, universeSpeed int, facilities Facilities) time.Duration {
	price := b.GetPrice(level)
	metalCost := float64(price.Metal)
	crystalCost := float64(price.Crystal)
	roboticLvl := float64(facilities.RoboticsFactory)
	naniteLvl := float64(facilities.NaniteFactory)
	hours := (metalCost + crystalCost) / (2500 * (1 + roboticLvl) * float64(universeSpeed) * math.Pow(2, naniteLvl))
	secs := hours * 3600
	if (level - 1) < 5 {
		secs = secs * (2 / (7 - (float64(level) - 1)))
	}
	secs = math.Max(1, secs)
	return time.Duration(int(math.Floor(secs))) * time.Second
}

// GetLevel returns current level of a building
func (b BaseBuilding) GetLevel(resourcesBuildings ResourcesBuildings, facilities Facilities, researches Researches) int {
	if b.ID.IsResourceBuilding() {
		return resourcesBuildings.ByID(b.ID)
	} else if b.ID.IsFacility() {
		return facilities.ByID(b.ID)
	}
	return 0
}
