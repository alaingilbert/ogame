package lightFighter

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// LightFighter ...
type LightFighter struct {
	baseShip.BaseShip
}

// New ...
func New() *LightFighter {
	l := new(LightFighter)
	l.OGameID = 204
	l.StructuralIntegrity = 4000
	l.ShieldPower = 10
	l.WeaponPower = 50
	l.CargoCapacity = 50
	l.BaseSpeed = 12500
	l.FuelConsumption = 20
	l.RapidfireFrom = map[ogame.ID]int{ogame.Cruiser: 6, ogame.Deathstar: 200}
	l.RapidfireAgainst = map[ogame.ID]int{ogame.EspionageProbe: 5, ogame.SolarSatellite: 5}
	l.Price = ogame.Resources{Metal: 3000, Crystal: 1000}
	l.Requirements = map[ogame.ID]int{ogame.Shipyard: 1, ogame.CombustionDrive: 1}
	return l
}
