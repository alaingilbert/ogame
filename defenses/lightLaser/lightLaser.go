package lightLaser

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/defenses/baseDefense"
)

// LightLaser ...
type LightLaser struct {
	baseDefense.BaseDefense
}

// New ...
func New() *LightLaser {
	d := new(LightLaser)
	d.OGameID = 402
	d.Price = ogame.Resources{Metal: 1500, Crystal: 500}
	d.StructuralIntegrity = 2000
	d.ShieldPower = 25
	d.WeaponPower = 100
	d.RapidfireFrom = map[ogame.ID]int{ogame.Destroyer: 10, ogame.Bomber: 20, ogame.Deathstar: 200}
	d.Requirements = map[ogame.ID]int{ogame.Shipyard: 2, ogame.LaserTechnology: 3}
	return d
}
