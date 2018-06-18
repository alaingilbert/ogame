package recycler

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// Recycler ...
type Recycler struct {
	baseShip.BaseShip
}

// New ...
func New() *Recycler {
	s := new(Recycler)
	s.OGameID = 209
	s.StructuralIntegrity = 16000
	s.ShieldPower = 10
	s.WeaponPower = 1
	s.CargoCapacity = 20000
	s.BaseSpeed = 2000
	s.FuelConsumption = 300
	s.RapidfireFrom = map[ogame.ID]int{ogame.Deathstar: 250}
	s.RapidfireAgainst = map[ogame.ID]int{ogame.EspionageProbe: 5, ogame.SolarSatellite: 5}
	s.Price = ogame.Resources{Metal: 10000, Crystal: 6000, Deuterium: 2000}
	s.Requirements = map[ogame.ID]int{ogame.Shipyard: 4, ogame.CombustionDrive: 6, ogame.ShieldingTechnology: 2}
	return s
}
