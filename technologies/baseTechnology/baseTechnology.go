package baseTechnology

import (
	"math"

	"github.com/alaingilbert/ogame"
)

// BaseTechnology ...
type BaseTechnology struct {
	OGameID        ogame.ID
	BaseCost       ogame.Resources
	IncreaseFactor float64
	Requirements   map[ogame.ID]int
}

// GetOGameID ...
func (b BaseTechnology) GetOGameID() ogame.ID {
	return b.OGameID
}

// GetBaseCost ...
func (b BaseTechnology) GetBaseCost() ogame.Resources {
	return b.BaseCost
}

// GetIncreaseFactor ...
func (b BaseTechnology) GetIncreaseFactor() float64 {
	return b.IncreaseFactor
}

// GetRequirements ...
func (b BaseTechnology) GetRequirements() map[ogame.ID]int {
	return b.Requirements
}

// GetPrice ...
func (b BaseTechnology) GetPrice(level int) ogame.Resources {
	return ogame.Resources{
		Metal:     researchCost(b.BaseCost.Metal, b.IncreaseFactor, level),
		Crystal:   researchCost(b.BaseCost.Crystal, b.IncreaseFactor, level),
		Deuterium: researchCost(b.BaseCost.Deuterium, b.IncreaseFactor, level),
		Energy:    researchCost(b.BaseCost.Energy, b.IncreaseFactor, level),
	}
}

// ConstructionTime ...
func (b BaseTechnology) ConstructionTime(level, universeSpeed int, facilities ogame.Facilities) int {
	price := b.GetPrice(level)
	metalCost := float64(price.Metal)
	crystalCost := float64(price.Crystal)
	researchLabLvl := float64(facilities.ResearchLab)
	hours := (metalCost + crystalCost) / (1000 * (1 + researchLabLvl) * float64(universeSpeed))
	secs := hours * 3600
	return int(math.Floor(secs))
}

// GetLevel ...
func (b BaseTechnology) GetLevel(resourcesBuildings ogame.ResourcesBuildings, facilities ogame.Facilities, researches ogame.Researches) int {
	return researches.ByOGameID(b.OGameID)
}

// IsAvailable ...
func (b BaseTechnology) IsAvailable(resourcesBuildings ogame.ResourcesBuildings, facilities ogame.Facilities,
	researches ogame.Researches, _ int) bool {
	for ogameID, levelNeeded := range b.Requirements {
		if ogameID.IsResourceBuilding() {
			if resourcesBuildings.ByOGameID(ogameID) < levelNeeded {
				return false
			}
		} else if ogameID.IsFacility() {
			if facilities.ByOGameID(ogameID) < levelNeeded {
				return false
			}
		} else if ogameID.IsTech() {
			if researches.ByOGameID(ogameID) < levelNeeded {
				return false
			}
		}
	}
	return true
}

func researchCost(baseCost int, increaseFactor float64, level int) int {
	return int(float64(baseCost) * math.Pow(increaseFactor, float64(level-1)))
}
