package heavyFighter

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// HeavyFighter ...
type HeavyFighter struct {
	baseShip.BaseShip
}

// New ...
func New() *HeavyFighter {
	s := new(HeavyFighter)
	s.OGameID = 205
	s.StructuralIntegrity = 10000
	s.ShieldPower = 25
	s.WeaponPower = 150
	s.CargoCapacity = 100
	s.BaseSpeed = 10000
	s.FuelConsumption = 75
	s.RapidfireFrom = map[ogame.ID]int{ogame.Battlecruiser: 4, ogame.Deathstar: 100}
	s.RapidfireAgainst = map[ogame.ID]int{ogame.EspionageProbe: 5, ogame.SolarSatellite: 5, ogame.SmallCargo: 3}
	s.Price = ogame.Resources{Metal: 6000, Crystal: 4000}
	s.Requirements = map[ogame.ID]int{ogame.Shipyard: 3, ogame.ImpulseDrive: 2, ogame.ArmourTechnology: 2}
	return s
}
