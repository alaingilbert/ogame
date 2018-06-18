package ionCannon

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/defenses/baseDefense"
)

// IonCannon ...
type IonCannon struct {
	baseDefense.BaseDefense
}

// New ...
func New() *IonCannon {
	d := new(IonCannon)
	d.OGameID = 405
	d.Price = ogame.Resources{Metal: 2000, Crystal: 6000}
	d.StructuralIntegrity = 8000
	d.ShieldPower = 500
	d.WeaponPower = 150
	d.RapidfireFrom = map[ogame.ID]int{ogame.Bomber: 10, ogame.Deathstar: 100}
	d.Requirements = map[ogame.ID]int{ogame.Shipyard: 4, ogame.IonTechnology: 4}
	return d
}
