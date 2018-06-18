package deuteriumSynthesizer

import (
	"math"

	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// DeuteriumSynthesizer ...
type DeuteriumSynthesizer struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *DeuteriumSynthesizer {
	b := new(DeuteriumSynthesizer)
	b.OGameID = 3
	b.IncreaseFactor = 1.5
	b.BaseCost = ogame.Resources{Metal: 225, Crystal: 75}
	return b
}

// EnergyConsumption ...
func (b *DeuteriumSynthesizer) EnergyConsumption(level int) int {
	return int(math.Ceil(20 * float64(level) * math.Pow(1.1, float64(level))))
}

// Production ...
func (b *DeuteriumSynthesizer) Production(universeSpeed, maxTemp int, productionRatio float64, level int) int {
	return int(math.Round(10 * float64(level) * math.Pow(1.1, float64(level)) * (1.44 - 0.004*float64(maxTemp)) * float64(universeSpeed) * productionRatio))
}
