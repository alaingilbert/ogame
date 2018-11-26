package ogame

import "math"

type crystalMine struct {
	BaseBuilding
}

func newCrystalMine() *crystalMine {
	b := new(crystalMine)
	b.Name = "crystal mine"
	b.ID = CrystalMineID
	b.IncreaseFactor = 1.6
	b.BaseCost = Resources{Metal: 48, Crystal: 24}
	return b
}

// EnergyConsumption returns the building energy consumption
func (b *crystalMine) EnergyConsumption(level int) int {
	return int(math.Ceil(10 * float64(level) * math.Pow(1.1, float64(level))))
}

// Production returns the crystal production of the mine
func (b *crystalMine) Production(universeSpeed int, productionRatio, globalRatio float64, plasmaTech, level int) int {
	basicIncome := 15.0
	levelProduction := 20 * float64(universeSpeed) * (1 + float64(plasmaTech)*0.0066) * float64(level) * math.Pow(1.1, float64(level))
	production := int(levelProduction*productionRatio*globalRatio + (basicIncome * float64(universeSpeed)))
	return production
}
