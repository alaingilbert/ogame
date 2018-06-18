package colonyShip

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// ColonyShip ...
type ColonyShip struct {
	baseShip.BaseShip
}

// New ...
func New() *ColonyShip {
	s := new(ColonyShip)
	s.OGameID = 208
	s.StructuralIntegrity = 30000
	s.ShieldPower = 100
	s.WeaponPower = 50
	s.CargoCapacity = 7500
	s.BaseSpeed = 2500
	s.FuelConsumption = 1000
	s.RapidfireFrom = map[ogame.ID]int{ogame.Deathstar: 250}
	s.RapidfireAgainst = map[ogame.ID]int{ogame.EspionageProbe: 5, ogame.SolarSatellite: 5}
	s.Price = ogame.Resources{Metal: 10000, Crystal: 20000, Deuterium: 10000}
	s.Requirements = map[ogame.ID]int{ogame.Shipyard: 4, ogame.ImpulseDrive: 3}
	return s
}
