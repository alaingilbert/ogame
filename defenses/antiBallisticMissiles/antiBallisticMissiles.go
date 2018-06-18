package antiBallisticMissiles

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/defenses/baseDefense"
)

// AntiBallisticMissiles ...
type AntiBallisticMissiles struct {
	baseDefense.BaseDefense
}

// New ...
func New() *AntiBallisticMissiles {
	d := new(AntiBallisticMissiles)
	d.OGameID = 502
	d.Price = ogame.Resources{Metal: 8000, Crystal: 2000}
	d.StructuralIntegrity = 8000
	d.ShieldPower = 1
	d.WeaponPower = 1
	d.Requirements = map[ogame.ID]int{ogame.MissileSilo: 2}
	return d
}
