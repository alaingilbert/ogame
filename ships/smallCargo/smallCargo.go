package smallCargo

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// SmallCargo ...
type SmallCargo struct {
	baseShip.BaseShip
}

// New ...
func New() *SmallCargo {
	s := new(SmallCargo)
	s.OGameID = 202
	s.StructuralIntegrity = 4000
	s.ShieldPower = 10
	s.WeaponPower = 5
	s.CargoCapacity = 5000
	s.BaseSpeed = 5000
	s.FuelConsumption = 10
	s.RapidfireFrom = map[ogame.ID]int{ogame.Battlecruiser: 3, ogame.HeavyFighter: 3, ogame.Deathstar: 250}
	s.RapidfireAgainst = map[ogame.ID]int{ogame.EspionageProbe: 5, ogame.SolarSatellite: 5}
	s.Price = ogame.Resources{Metal: 2000, Crystal: 2000}
	s.Requirements = map[ogame.ID]int{ogame.Shipyard: 2, ogame.CombustionDrive: 2}
	return s
}
