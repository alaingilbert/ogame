package researchLab

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// ResearchLab ...
type ResearchLab struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *ResearchLab {
	b := new(ResearchLab)
	b.OGameID = 31
	b.IncreaseFactor = 2.0
	b.BaseCost = ogame.Resources{Metal: 200, Crystal: 400, Deuterium: 200}
	return b
}
