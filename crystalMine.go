package ogame

import "math"

// CrystalMine ...
type crystalMine struct {
	BaseBuilding
}

func newCrystalMine() *crystalMine {
	b := new(crystalMine)
	b.ID = CrystalMineID
	b.IncreaseFactor = 1.6
	b.BaseCost = Resources{Metal: 48, Crystal: 24}
	return b
}

// EnergyConsumption ...
func (b *crystalMine) EnergyConsumption(level int) int {
	return int(math.Ceil(10 * float64(level) * math.Pow(1.1, float64(level))))
}

// Production ...
func (b *crystalMine) Production(universeSpeed int, productionRatio float64, level int) int {
	basicIncome := 15.0
	levelProduction := 20 * float64(level) * math.Pow(1.1, float64(level)) * float64(universeSpeed)
	production := int(levelProduction*productionRatio + (basicIncome * float64(universeSpeed)))
	return production
}
