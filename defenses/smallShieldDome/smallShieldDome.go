package smallShieldDome

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/defenses/baseDefense"
)

// SmallShieldDome ...
type SmallShieldDome struct {
	baseDefense.BaseDefense
}

// New ...
func New() *SmallShieldDome {
	d := new(SmallShieldDome)
	d.OGameID = 407
	d.Price = ogame.Resources{Metal: 10000, Crystal: 10000}
	d.StructuralIntegrity = 20000
	d.ShieldPower = 2000
	d.WeaponPower = 1
	d.Requirements = map[ogame.ID]int{ogame.Shipyard: 1, ogame.ShieldingTechnology: 2}
	return d
}
