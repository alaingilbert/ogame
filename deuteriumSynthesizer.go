package ogame

import "math"

// DeuteriumSynthesizer ...
type deuteriumSynthesizer struct {
	BaseBuilding
}

func newDeuteriumSynthesizer() *deuteriumSynthesizer {
	b := new(deuteriumSynthesizer)
	b.Name = "deuterium synthesizer"
	b.ID = DeuteriumSynthesizerID
	b.IncreaseFactor = 1.5
	b.BaseCost = Resources{Metal: 225, Crystal: 75}
	return b
}

// EnergyConsumption returns the building energy consumption
func (b *deuteriumSynthesizer) EnergyConsumption(level int) int {
	return int(math.Ceil(20 * float64(level) * math.Pow(1.1, float64(level))))
}

// Production returns the deuterium production of the mine
func (b *deuteriumSynthesizer) Production(universeSpeed, maxTemp int, productionRatio float64, level int) int {
	return int(math.Round(10 * float64(level) * math.Pow(1.1, float64(level)) * (1.44 - 0.004*float64(maxTemp)) * float64(universeSpeed) * productionRatio))
}
