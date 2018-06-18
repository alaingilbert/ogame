package cruiser

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// Cruiser ...
type Cruiser struct {
	baseShip.BaseShip
}

// New ...
func New() *Cruiser {
	s := new(Cruiser)
	s.OGameID = 206
	s.StructuralIntegrity = 27000
	s.ShieldPower = 50
	s.WeaponPower = 400
	s.CargoCapacity = 800
	s.BaseSpeed = 15000
	s.FuelConsumption = 300
	s.RapidfireFrom = map[ogame.ID]int{ogame.Battlecruiser: 4, ogame.Deathstar: 33}
	s.RapidfireAgainst = map[ogame.ID]int{ogame.EspionageProbe: 5, ogame.SolarSatellite: 5,
		ogame.LightFighter: 6, ogame.RocketLauncher: 10}
	s.Price = ogame.Resources{Metal: 20000, Crystal: 7000, Deuterium: 2000}
	s.Requirements = map[ogame.ID]int{ogame.Shipyard: 5, ogame.ImpulseDrive: 4, ogame.IonTechnology: 2}
	return s
}
