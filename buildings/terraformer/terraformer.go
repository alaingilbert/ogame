package terraformer

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// Terraformer ...
type Terraformer struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *Terraformer {
	b := new(Terraformer)
	b.OGameID = 33
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 0, Crystal: 50000, Deuterium: 100000}
	b.Requirements = map[ogame.ID]int{ogame.NaniteFactory: 1, ogame.EnergyTechnology: 12}
	return b
}
