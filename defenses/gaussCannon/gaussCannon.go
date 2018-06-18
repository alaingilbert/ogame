package gaussCannon

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/defenses/baseDefense"
)

// GaussCannon ...
type GaussCannon struct {
	baseDefense.BaseDefense
}

// New ...
func New() *GaussCannon {
	d := new(GaussCannon)
	d.OGameID = 404
	d.Price = ogame.Resources{Metal: 20000, Crystal: 15000, Deuterium: 2000}
	d.StructuralIntegrity = 35000
	d.ShieldPower = 200
	d.WeaponPower = 1100
	d.RapidfireFrom = map[ogame.ID]int{ogame.Deathstar: 50}
	d.Requirements = map[ogame.ID]int{ogame.Shipyard: 6, ogame.WeaponsTechnology: 3, ogame.EnergyTechnology: 6,
		ogame.ShieldingTechnology: 1}
	return d
}
