package metalMine

import (
	"math"

	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// MetalMine ...
type MetalMine struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *MetalMine {
	b := new(MetalMine)
	b.OGameID = 1
	b.IncreaseFactor = 1.5
	b.BaseCost = ogame.Resources{Metal: 60, Crystal: 15}
	return b
}

// EnergyConsumption ...
func (b *MetalMine) EnergyConsumption(level int) int {
	return int(math.Ceil(10 * float64(level) * math.Pow(1.1, float64(level))))
}

// Production ...
func (b *MetalMine) Production(universeSpeed int, productionRatio float64, level int) int {
	basicIncome := 30.0 * float64(universeSpeed)
	levelProduction := basicIncome * float64(level) * math.Pow(1.1, float64(level))
	production := int(levelProduction*productionRatio + basicIncome)
	return production
}
