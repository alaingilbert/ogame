package ogame

import "math"

type solarPlant struct {
	BaseBuilding
}

func newSolarPlant() *solarPlant {
	b := new(solarPlant)
	b.Name = "solar plant"
	b.ID = SolarPlantID
	b.IncreaseFactor = 1.5
	b.BaseCost = Resources{Metal: 75, Crystal: 30}
	return b
}

// Production returns the energy produced by the solar plant at provided level
func (b *solarPlant) Production(level int64) int64 {
	return int64(math.Floor(20 * float64(level) * math.Pow(1.1, float64(level))))
}
