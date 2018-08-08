package ogame

import "math"

// SolarPlant ...
type solarPlant struct {
	BaseBuilding
}

func newSolarPlant() *solarPlant {
	b := new(solarPlant)
	b.ID = SolarPlantID
	b.IncreaseFactor = 1.5
	b.BaseCost = Resources{Metal: 75, Crystal: 30}
	return b
}

// Production ...
func (b *solarPlant) Production(level int) int {
	return int(math.Floor(20 * float64(level) * math.Pow(1.1, float64(level))))
}
