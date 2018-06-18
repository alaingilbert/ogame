package shipyard

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// Shipyard ...
type Shipyard struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *Shipyard {
	b := new(Shipyard)
	b.OGameID = 21
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 400, Crystal: 200, Deuterium: 100}
	b.Requirements = map[ogame.ID]int{ogame.RoboticsFactory: 2}
	return b
}
