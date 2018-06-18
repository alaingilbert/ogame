package largeCargo

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// LargeCargo ...
type LargeCargo struct {
	baseShip.BaseShip
}

// New ...
func New() *LargeCargo {
	s := new(LargeCargo)
	s.OGameID = 203
	s.StructuralIntegrity = 12000
	s.ShieldPower = 25
	s.WeaponPower = 5
	s.CargoCapacity = 25000
	s.BaseSpeed = 7500
	s.FuelConsumption = 50
	s.RapidfireFrom = map[ogame.ID]int{ogame.Battlecruiser: 3, ogame.Deathstar: 250}
	s.RapidfireAgainst = map[ogame.ID]int{ogame.EspionageProbe: 5, ogame.SolarSatellite: 5}
	s.Price = ogame.Resources{Metal: 6000, Crystal: 6000}
	s.Requirements = map[ogame.ID]int{ogame.Shipyard: 4, ogame.CombustionDrive: 6}
	return s
}
