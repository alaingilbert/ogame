package ogame

import "math"

// MetalMine ...
type metalMine struct {
	BaseBuilding
}

// NeNewMetalMinew ...
func NewMetalMine() *metalMine {
	b := new(metalMine)
	b.ID = MetalMineID
	b.IncreaseFactor = 1.5
	b.BaseCost = Resources{Metal: 60, Crystal: 15}
	return b
}

// EnergyConsumption ...
func (b *metalMine) EnergyConsumption(level int) int {
	return int(math.Ceil(10 * float64(level) * math.Pow(1.1, float64(level))))
}

// Production ...
func (b *metalMine) Production(universeSpeed int, productionRatio float64, level int) int {
	basicIncome := 30.0 * float64(universeSpeed)
	levelProduction := basicIncome * float64(level) * math.Pow(1.1, float64(level))
	production := int(levelProduction*productionRatio + basicIncome)
	return production
}
