package plasmaTurret

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/defenses/baseDefense"
)

// PlasmaTurret ...
type PlasmaTurret struct {
	baseDefense.BaseDefense
}

// New ...
func New() *PlasmaTurret {
	d := new(PlasmaTurret)
	d.OGameID = 406
	d.Price = ogame.Resources{Metal: 50000, Crystal: 50000, Deuterium: 30000}
	d.StructuralIntegrity = 100000
	d.ShieldPower = 300
	d.WeaponPower = 3000
	d.Requirements = map[ogame.ID]int{ogame.Shipyard: 8, ogame.PlasmaTechnology: 7}
	return d
}
