package undergroundCrystalDen

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// IsAvailable ...
func IsAvailable() bool {
	return true
}

// UndergroundCrystalDen ...
type UndergroundCrystalDen struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *UndergroundCrystalDen {
	b := new(UndergroundCrystalDen)
	b.OGameID = 26
	b.IncreaseFactor = 2.3
	b.BaseCost = ogame.Resources{Metal: 2645, Crystal: 1322}
	return b
}
