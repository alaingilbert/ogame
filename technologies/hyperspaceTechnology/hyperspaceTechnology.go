package hyperspaceTechnology

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// HyperspaceTechnology ...
type HyperspaceTechnology struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *HyperspaceTechnology {
	b := new(HyperspaceTechnology)
	b.OGameID = 114
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Crystal: 4000, Deuterium: 2000}
	b.Requirements = map[ogame.ID]int{ogame.ResearchLab: 7, ogame.ShieldingTechnology: 5, ogame.EnergyTechnology: 5}
	return b
}
