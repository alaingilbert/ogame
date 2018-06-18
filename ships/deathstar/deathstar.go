package deathstar

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// Deathstar ...
type Deathstar struct {
	baseShip.BaseShip
}

// New ...
func New() *Deathstar {
	s := new(Deathstar)
	s.OGameID = 214
	s.StructuralIntegrity = 9000000
	s.ShieldPower = 50000
	s.WeaponPower = 200000
	s.CargoCapacity = 1000000
	s.BaseSpeed = 100
	s.FuelConsumption = 1
	s.RapidfireFrom = map[ogame.ID]int{}
	s.RapidfireAgainst = map[ogame.ID]int{ogame.SmallCargo: 250, ogame.LargeCargo: 250, ogame.LightFighter: 200,
		ogame.HeavyFighter: 100, ogame.Cruiser: 33, ogame.Battleship: 30, ogame.ColonyShip: 250,
		ogame.Recycler: 250, ogame.EspionageProbe: 1250, ogame.SolarSatellite: 1250, ogame.Bomber: 25,
		ogame.Destroyer: 5, ogame.RocketLauncher: 200, ogame.LightLaser: 200, ogame.HeavyLaser: 100,
		ogame.GaussCannon: 50, ogame.IonCannon: 100, ogame.Battlecruiser: 15}
	s.Price = ogame.Resources{Metal: 5000000, Crystal: 4000000, Deuterium: 1000000}
	s.Requirements = map[ogame.ID]int{ogame.Shipyard: 12, ogame.GravitonTechnology: 1, ogame.HyperspaceDrive: 7,
		ogame.HyperspaceTechnology: 6}
	return s
}
