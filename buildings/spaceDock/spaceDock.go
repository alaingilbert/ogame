package spaceDock

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// SpaceDock ...
type SpaceDock struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *SpaceDock {
	b := new(SpaceDock)
	b.OGameID = 36
	b.IncreaseFactor = 5
	b.BaseCost = ogame.Resources{Metal: 200, Crystal: 50, Energy: 50}
	b.Requirements = map[ogame.ID]int{ogame.Shipyard: 2}
	return b
}
