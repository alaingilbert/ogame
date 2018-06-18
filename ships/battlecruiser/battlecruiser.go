package battlecruiser

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// Battlecruiser ...
type Battlecruiser struct {
	baseShip.BaseShip
}

// New ...
func New() *Battlecruiser {
	b := new(Battlecruiser)
	b.OGameID = 215
	b.StructuralIntegrity = 70000
	b.ShieldPower = 400
	b.WeaponPower = 700
	b.CargoCapacity = 750
	b.BaseSpeed = 1000
	b.FuelConsumption = 250
	b.RapidfireFrom = map[ogame.ID]int{ogame.Destroyer: 2, ogame.Deathstar: 15}
	b.RapidfireAgainst = map[ogame.ID]int{
		ogame.EspionageProbe: 5, ogame.SolarSatellite: 5, ogame.SmallCargo: 3, ogame.LargeCargo: 3,
		ogame.HeavyFighter: 4, ogame.Cruiser: 4, ogame.Battleship: 7,
	}
	b.Price = ogame.Resources{Metal: 30000, Crystal: 40000, Deuterium: 15000}
	b.Requirements = map[ogame.ID]int{ogame.LaserTechnology: 12, ogame.HyperspaceTechnology: 5,
		ogame.HyperspaceDrive: 5, ogame.Shipyard: 8}
	return b
}
