package solarPlant

import (
	"math"

	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// IsAvailable ...
func IsAvailable() bool {
	return true
}

// SolarPlant ...
type SolarPlant struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *SolarPlant {
	b := new(SolarPlant)
	b.OGameID = 4
	b.IncreaseFactor = 1.5
	b.BaseCost = ogame.Resources{Metal: 75, Crystal: 30}
	return b
}

// Production ...
func (b *SolarPlant) Production(level int) int {
	return int(math.Floor(20 * float64(level) * math.Pow(1.1, float64(level))))
}
