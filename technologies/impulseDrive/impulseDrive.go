package impulseDrive

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// ImpulseDrive ...
type ImpulseDrive struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *ImpulseDrive {
	b := new(ImpulseDrive)
	b.OGameID = 117
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 2000, Crystal: 4000, Deuterium: 600}
	b.Requirements = map[ogame.ID]int{ogame.ResearchLab: 2, ogame.EnergyTechnology: 1}
	return b
}
