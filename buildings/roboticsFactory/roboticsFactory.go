package roboticsFactory

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// RoboticsFactory ...
type RoboticsFactory struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *RoboticsFactory {
	b := new(RoboticsFactory)
	b.OGameID = 14
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 400, Crystal: 120, Deuterium: 200}
	return b
}
