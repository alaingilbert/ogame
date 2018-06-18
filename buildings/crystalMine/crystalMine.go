package crystalMine

import (
	"math"

	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// CrystalMine ...
type CrystalMine struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *CrystalMine {
	b := new(CrystalMine)
	b.OGameID = 2
	b.IncreaseFactor = 1.6
	b.BaseCost = ogame.Resources{Metal: 48, Crystal: 24}
	return b
}

// EnergyConsumption ...
func (b *CrystalMine) EnergyConsumption(level int) int {
	return int(math.Ceil(10 * float64(level) * math.Pow(1.1, float64(level))))
}

// Production ...
func (b *CrystalMine) Production(universeSpeed int, productionRatio float64, level int) int {
	basicIncome := 15.0
	levelProduction := 20 * float64(level) * math.Pow(1.1, float64(level)) * float64(universeSpeed)
	production := int(levelProduction*productionRatio + (basicIncome * float64(universeSpeed)))
	return production
}
