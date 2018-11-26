package ogame

import "math"

type fusionReactor struct {
	BaseBuilding
}

func newFusionReactor() *fusionReactor {
	b := new(fusionReactor)
	b.Name = "fusion reactor"
	b.ID = FusionReactorID
	b.IncreaseFactor = 1.8
	b.BaseCost = Resources{Metal: 900, Crystal: 360, Deuterium: 180}
	b.Requirements = map[ID]int{DeuteriumSynthesizerID: 5, EnergyTechnologyID: 3}
	return b
}

// Production returns the energy production of the reactor
func (b *fusionReactor) Production(energyTechnology, lvl int) int {
	pct := 1.0
	lvlf := float64(lvl)
	energyTechnologyf := float64(energyTechnology)
	return int(math.Round(30 * lvlf * math.Pow(1.05+energyTechnologyf*0.01, lvlf) * pct))
}

// GetFuelConsumption returns the deuterium consumed by the fusion reactor
func (b fusionReactor) GetFuelConsumption(universeSpeed int, ratio float64, lvl int) int {
	return int(math.Abs(math.Floor(-10 * float64(universeSpeed) * float64(lvl) * math.Pow(1.1, float64(lvl)) * ratio)))
}
