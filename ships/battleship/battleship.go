package battleship

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// Battleship ...
type Battleship struct {
	baseShip.BaseShip
}

// New ...
func New() *Battleship {
	b := new(Battleship)
	b.OGameID = 207
	b.StructuralIntegrity = 60000
	b.ShieldPower = 200
	b.WeaponPower = 1000
	b.CargoCapacity = 1500
	b.BaseSpeed = 10000
	b.FuelConsumption = 500
	b.RapidfireFrom = map[ogame.ID]int{ogame.Battlecruiser: 7, ogame.Deathstar: 30}
	b.RapidfireAgainst = map[ogame.ID]int{ogame.EspionageProbe: 5, ogame.SolarSatellite: 5}
	b.Price = ogame.Resources{Metal: 45000, Crystal: 15000}
	b.Requirements = map[ogame.ID]int{ogame.Shipyard: 7, ogame.HyperspaceDrive: 4}
	return b
}
