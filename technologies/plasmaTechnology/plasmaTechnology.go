package plasmaTechnology

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// PlasmaTechnology ...
type PlasmaTechnology struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *PlasmaTechnology {
	b := new(PlasmaTechnology)
	b.OGameID = 122
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 2000, Crystal: 4000, Deuterium: 1000}
	b.Requirements = map[ogame.ID]int{ogame.IonTechnology: 5, ogame.EnergyTechnology: 8, ogame.LaserTechnology: 10}
	return b
}
