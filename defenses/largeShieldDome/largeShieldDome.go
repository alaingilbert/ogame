package largeShieldDome

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/defenses/baseDefense"
)

// LargeShieldDome ...
type LargeShieldDome struct {
	baseDefense.BaseDefense
}

// New ...
func New() *LargeShieldDome {
	d := new(LargeShieldDome)
	d.OGameID = 408
	d.Price = ogame.Resources{Metal: 50000, Crystal: 50000}
	d.StructuralIntegrity = 100000
	d.ShieldPower = 10000
	d.WeaponPower = 1
	d.Requirements = map[ogame.ID]int{ogame.ShieldingTechnology: 6, ogame.Shipyard: 6}
	return d
}
