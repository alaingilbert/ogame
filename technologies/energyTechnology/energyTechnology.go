package energyTechnology

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// EnergyTechnology ...
type EnergyTechnology struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *EnergyTechnology {
	b := new(EnergyTechnology)
	b.OGameID = 113
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Crystal: 800, Deuterium: 400}
	b.Requirements = map[ogame.ID]int{ogame.ResearchLab: 1}
	return b
}

// IsAvailable ...
func (t *EnergyTechnology) IsAvailable(_ ogame.ResourcesBuildings, facilities ogame.Facilities, _ ogame.Researches, _ int) bool {
	return facilities.ResearchLab >= 1
}
