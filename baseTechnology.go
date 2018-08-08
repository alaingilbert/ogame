package ogame

import (
	"math"
)

// BaseTechnology ...
type BaseTechnology struct {
	ID             ID
	BaseCost       Resources
	IncreaseFactor float64
	Requirements   map[ID]int
}

// GetOGameID ...
func (b BaseTechnology) GetOGameID() ID {
	return b.ID
}

// GetBaseCost ...
func (b BaseTechnology) GetBaseCost() Resources {
	return b.BaseCost
}

// GetIncreaseFactor ...
func (b BaseTechnology) GetIncreaseFactor() float64 {
	return b.IncreaseFactor
}

// GetRequirements ...
func (b BaseTechnology) GetRequirements() map[ID]int {
	return b.Requirements
}

// GetPrice ...
func (b BaseTechnology) GetPrice(level int) Resources {
	return Resources{
		Metal:     researchCost(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:   researchCost(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium: researchCost(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:    researchCost(b.BaseCost.Energy, b.IncreaseFactor, level),
	}
}

// ConstructionTime ...
func (b BaseTechnology) ConstructionTime(level, universeSpeed int, facilities Facilities) int {
	price := b.GetPrice(level)
	metalCost := float64(price.Metal)
	crystalCost := float64(price.Crystal)
	researchLabLvl := float64(facilities.ResearchLab)
	hours := (metalCost + crystalCost) / (1000 * (1 + researchLabLvl) * float64(universeSpeed))
	secs := hours * 3600
	return int(math.Floor(secs))
}

// GetLevel ...
func (b BaseTechnology) GetLevel(resourcesBuildings ResourcesBuildings, facilities Facilities, researches Researches) int {
	return researches.ByID(b.ID)
}

// IsAvailable ...
func (b BaseTechnology) IsAvailable(resourcesBuildings ResourcesBuildings, facilities Facilities,
	researches Researches, _ int) bool {
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

func researchCost(baseCost int, increaseFactor float64, level int) int {
	return int(float64(baseCost) * math.Pow(increaseFactor, float64(level-1)))
}
