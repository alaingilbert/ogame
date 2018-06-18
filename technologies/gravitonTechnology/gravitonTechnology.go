package gravitonTechnology

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// GravitonTechnology ...
type GravitonTechnology struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *GravitonTechnology {
	b := new(GravitonTechnology)
	b.OGameID = 199
	b.IncreaseFactor = 3.0
	b.BaseCost = ogame.Resources{Energy: 300000}
	b.Requirements = map[ogame.ID]int{ogame.ResearchLab: 12}
	return b
}

// IsAvailable ...
func (b GravitonTechnology) IsAvailable(resourcesBuildings ogame.ResourcesBuildings, facilities ogame.Facilities,
	researches ogame.Researches, energy int) bool {
	if energy < 300000 {
		return false
	}
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
