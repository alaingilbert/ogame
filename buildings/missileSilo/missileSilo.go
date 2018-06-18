package missileSilo

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// MissileSilo ...
type MissileSilo struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *MissileSilo {
	b := new(MissileSilo)
	b.OGameID = 44
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 20000, Crystal: 20000, Deuterium: 1000}
	b.Requirements = map[ogame.ID]int{ogame.Shipyard: 1}
	return b
}
