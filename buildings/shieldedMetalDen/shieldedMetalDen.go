package shieldedMetalDen

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// IsAvailable ...
func IsAvailable() bool {
	return true
}

// ShieldedMetalDen ...
type ShieldedMetalDen struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *ShieldedMetalDen {
	b := new(ShieldedMetalDen)
	b.OGameID = 25
	b.IncreaseFactor = 2.3
	b.BaseCost = ogame.Resources{Metal: 2645}
	return b
}
