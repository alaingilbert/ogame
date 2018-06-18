package heavyLaser

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/defenses/baseDefense"
)

// HeavyLaser ...
type HeavyLaser struct {
	baseDefense.BaseDefense
}

// New ...
func New() *HeavyLaser {
	d := new(HeavyLaser)
	d.OGameID = 403
	d.Price = ogame.Resources{Metal: 6000, Crystal: 2000}
	d.StructuralIntegrity = 8000
	d.ShieldPower = 100
	d.WeaponPower = 250
	d.RapidfireFrom = map[ogame.ID]int{ogame.Bomber: 10, ogame.Deathstar: 100}
	d.Requirements = map[ogame.ID]int{ogame.Shipyard: 4, ogame.EnergyTechnology: 3, ogame.LaserTechnology: 6}
	return d
}
