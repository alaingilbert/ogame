package interplanetaryMissiles

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/defenses/baseDefense"
)

// InterplanetaryMissiles ...
type InterplanetaryMissiles struct {
	baseDefense.BaseDefense
}

// New ...
func New() *InterplanetaryMissiles {
	d := new(InterplanetaryMissiles)
	d.OGameID = 503
	d.Price = ogame.Resources{Metal: 12500, Crystal: 2500, Deuterium: 10000}
	d.StructuralIntegrity = 15000
	d.ShieldPower = 1
	d.WeaponPower = 12000
	d.Requirements = map[ogame.ID]int{ogame.MissileSilo: 4, ogame.ImpulseDrive: 1}
	return d
}
