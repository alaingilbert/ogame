package naniteFactory

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// NaniteFactory ...
type NaniteFactory struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *NaniteFactory {
	b := new(NaniteFactory)
	b.OGameID = 15
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 1000000, Crystal: 500000, Deuterium: 100000}
	b.Requirements = map[ogame.ID]int{ogame.RoboticsFactory: 10, ogame.ComputerTechnology: 10}
	return b
}
