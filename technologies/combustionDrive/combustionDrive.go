package combustionDrive

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/technologies/baseTechnology"
)

// CombustionDrive ...
type CombustionDrive struct {
	baseTechnology.BaseTechnology
}

// New ...
func New() *CombustionDrive {
	b := new(CombustionDrive)
	b.OGameID = 115
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 400, Deuterium: 600}
	b.Requirements = map[ogame.ID]int{ogame.EnergyTechnology: 1}
	return b
}
