package seabedDeuteriumDen

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// IsAvailable ...
func IsAvailable() bool {
	return true
}

// SeabedDeuteriumDen ...
type SeabedDeuteriumDen struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *SeabedDeuteriumDen {
	b := new(SeabedDeuteriumDen)
	b.OGameID = 27
	b.IncreaseFactor = 2.3
	b.BaseCost = ogame.Resources{Metal: 2645, Crystal: 2645}
	return b
}
