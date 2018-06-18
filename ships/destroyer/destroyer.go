package destroyer

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// Destroyer ...
type Destroyer struct {
	baseShip.BaseShip
}

// New ...
func New() *Destroyer {
	s := new(Destroyer)
	s.OGameID = 213
	s.StructuralIntegrity = 110000
	s.ShieldPower = 500
	s.WeaponPower = 2000
	s.CargoCapacity = 2000
	s.BaseSpeed = 5000
	s.FuelConsumption = 1000
	s.RapidfireFrom = map[ogame.ID]int{ogame.Deathstar: 5}
	s.RapidfireAgainst = map[ogame.ID]int{ogame.EspionageProbe: 5, ogame.SolarSatellite: 5,
		ogame.LightLaser: 10, ogame.Battlecruiser: 2}
	s.Price = ogame.Resources{Metal: 60000, Crystal: 50000, Deuterium: 15000}
	s.Requirements = map[ogame.ID]int{ogame.Shipyard: 9, ogame.HyperspaceDrive: 6, ogame.HyperspaceTechnology: 5}
	return s
}
