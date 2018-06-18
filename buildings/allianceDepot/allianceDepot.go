package allianceDepot

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// AllianceDepot ...
type AllianceDepot struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *AllianceDepot {
	b := new(AllianceDepot)
	b.OGameID = 34
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 20000, Crystal: 40000}
	return b
}
