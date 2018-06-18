package laserTechnology

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// LaserTechnology ...
type LaserTechnology struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *LaserTechnology {
	b := new(LaserTechnology)
	b.OGameID = 120
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 200, Crystal: 100}
	b.Requirements = map[ogame.ID]int{ogame.EnergyTechnology: 2}
	return b
}
