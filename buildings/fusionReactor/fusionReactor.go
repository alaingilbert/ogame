package fusionReactor

import (
	"math"

	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/buildings/baseBuilding"
)

// FusionReactor ...
type FusionReactor struct {
	baseBuilding.BaseBuilding
}

// New ...
func New() *FusionReactor {
	b := new(FusionReactor)
	b.OGameID = 12
	b.IncreaseFactor = 1.8
	b.BaseCost = ogame.Resources{Metal: 900, Crystal: 360, Deuterium: 180}
	b.Requirements = map[ogame.ID]int{ogame.DeuteriumSynthesizer: 5, ogame.EnergyTechnology: 3}
	return b
}

// Production ...
func (b *FusionReactor) Production(energyTechnology, lvl int) int {
	pct := 1.0
	lvlf := float64(lvl)
	energyTechnologyf := float64(energyTechnology)
	return int(math.Round(30 * lvlf * math.Pow(1.05+energyTechnologyf*0.01, lvlf) * pct))
}
