package ogame

import "math"

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
func (b *deuteriumSynthesizer) EnergyConsumption(level int64) int64 {
	return int64(math.Ceil(20 * float64(level) * math.Pow(1.1, float64(level))))
}

// Production returns the deuterium production of the mine
func (b *deuteriumSynthesizer) Production(universeSpeed, avgTemp int64, productionRatio, globalRatio float64, plasmaTech, level int64) int64 {
	return int64(math.Round(10 * (1 + float64(plasmaTech)*0.0033) * float64(level) * math.Pow(1.1, float64(level)) * (-0.004*float64(avgTemp) + 1.36) * float64(universeSpeed) * productionRatio * globalRatio))
}
