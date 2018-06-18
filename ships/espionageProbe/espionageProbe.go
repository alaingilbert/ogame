package espionageProbe

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// EspionageProbe ...
type EspionageProbe struct {
	baseShip.BaseShip
}

// New ...
func New() *EspionageProbe {
	s := new(EspionageProbe)
	s.OGameID = 210
	s.StructuralIntegrity = 1000
	s.ShieldPower = 0 //0.01
	s.WeaponPower = 0 //0.01
	s.CargoCapacity = 5
	s.BaseSpeed = 100000000
	s.FuelConsumption = 1
	s.RapidfireFrom = map[ogame.ID]int{ogame.Battlecruiser: 5, ogame.Destroyer: 5, ogame.Bomber: 5,
		ogame.Recycler: 5, ogame.ColonyShip: 5, ogame.Battleship: 5, ogame.Cruiser: 5,
		ogame.HeavyFighter: 5, ogame.LightFighter: 5, ogame.LargeCargo: 5, ogame.Deathstar: 1250,
		ogame.SmallCargo: 5}
	s.Price = ogame.Resources{Crystal: 1000}
	s.Requirements = map[ogame.ID]int{ogame.Shipyard: 3, ogame.CombustionDrive: 3, ogame.EspionageTechnology: 2}
	return s
}
