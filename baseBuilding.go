package ogame

import "math"

// BaseBuilding ...
type BaseBuilding struct {
	ID             ID
	BaseCost       Resources
	IncreaseFactor float64
	Requirements   map[ID]int
}

// GetID ...
func (b BaseBuilding) GetID() ID {
	return b.ID
}

// GetBaseCost ...
func (b BaseBuilding) GetBaseCost() Resources {
	return b.BaseCost
}

// GetIncreaseFactor ...
func (b BaseBuilding) GetIncreaseFactor() float64 {
	return b.IncreaseFactor
}

// GetRequirements ...
func (b BaseBuilding) GetRequirements() map[ID]int {
	return b.Requirements
}

// GetPrice ...
func (b BaseBuilding) GetPrice(level int) Resources {
	return Resources{
		Metal:     buildingCost(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:   buildingCost(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium: buildingCost(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:    buildingCost(b.BaseCost.Energy, b.IncreaseFactor, level),
	}
}

// ConstructionTime ...
func (b BaseBuilding) ConstructionTime(level, universeSpeed int, facilities Facilities) int {
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
	return int(math.Floor(secs))
}

// GetLevel ...
func (b BaseBuilding) GetLevel(resourcesBuildings ResourcesBuildings, facilities Facilities, researches Researches) int {
	if b.ID.IsResourceBuilding() {
		return resourcesBuildings.ByID(b.ID)
	} else if b.ID.IsFacility() {
		return facilities.ByID(b.ID)
	}
	return 0
}

// IsAvailable ...
func (b BaseBuilding) IsAvailable(resourcesBuildings ResourcesBuildings, facilities Facilities, researches Researches, _ int) bool {
	for id, levelNeeded := range b.Requirements {
		if id.IsResourceBuilding() {
			if resourcesBuildings.ByID(id) < levelNeeded {
				return false
			}
		} else if id.IsFacility() {
			if facilities.ByID(id) < levelNeeded {
				return false
			}
		} else if id.IsTech() {
			if researches.ByID(id) < levelNeeded {
				return false
			}
		}
	}
	return true
}

func buildingCost(baseCost int, increaseFactor float64, level int) int {
	return int(float64(baseCost) * math.Pow(increaseFactor, float64(level-1)))
}
