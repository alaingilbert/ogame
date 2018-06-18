package ionTechnology

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// IonTechnology ...
type IonTechnology struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *IonTechnology {
	b := new(IonTechnology)
	b.OGameID = 121
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 1000, Crystal: 300, Deuterium: 100}
	b.Requirements = map[ogame.ID]int{ogame.ResearchLab: 4, ogame.EnergyTechnology: 4, ogame.LaserTechnology: 5}
	return b
}
