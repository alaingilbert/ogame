package ogame

import "math"

type metalMine struct {
	BaseBuilding
}

func newMetalMine() *metalMine {
	b := new(metalMine)
	b.Name = "metal mine"
	b.ID = MetalMineID
	b.IncreaseFactor = 1.5
	b.BaseCost = Resources{Metal: 60, Crystal: 15}
	return b
}

// EnergyConsumption returns the building energy consumption
func (b *metalMine) EnergyConsumption(level int64) int64 {
	return int64(math.Ceil(10 * float64(level) * math.Pow(1.1, float64(level))))
}

// Production returns the metal production of the mine
func (b *metalMine) Production(universeSpeed int64, productionRatio, globalRatio float64, plasmaTech, level int64) int64 {
	basicIncome := 30.0 * float64(universeSpeed)
	levelProduction := 30.0 * (1.0 + (float64(plasmaTech) / 100.0)) * float64(universeSpeed) * float64(level) * math.Pow(1.1, float64(level))
	production := int64((levelProduction * productionRatio * globalRatio) + basicIncome)
	return production
}
