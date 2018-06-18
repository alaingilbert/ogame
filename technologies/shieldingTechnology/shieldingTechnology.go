package shieldingTechnology

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// ShieldingTechnology ...
type ShieldingTechnology struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *ShieldingTechnology {
	b := new(ShieldingTechnology)
	b.OGameID = 110
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 200, Crystal: 600}
	b.Requirements = map[ogame.ID]int{ogame.ResearchLab: 6, ogame.EnergyTechnology: 3}
	return b
}
