package solarSatellite

import (
	"github.com/alaingilbert/ogame"
	"github.com/alaingilbert/ogame/ships/baseShip"
)

// SolarSatellite ...
type SolarSatellite struct {
	baseShip.BaseShip
}

// New ...
func New() *SolarSatellite {
	s := new(SolarSatellite)
	s.OGameID = 212
	s.StructuralIntegrity = 2000
	s.ShieldPower = 1
	s.WeaponPower = 1
	s.CargoCapacity = 0
	s.BaseSpeed = 0
	s.FuelConsumption = 0
	s.RapidfireFrom = map[ogame.ID]int{ogame.Deathstar: 1250}
	s.Price = ogame.Resources{Crystal: 2000, Deuterium: 500}
	s.Requirements = map[ogame.ID]int{ogame.Shipyard: 1}
	return s
}

// Production ...
func (s *SolarSatellite) Production(temperatureMax, nbr int) int {
	return int((float64(temperatureMax)+140)/6) * nbr
}

// IsAvailable ...
func IsAvailable(shipyard int) bool {
	return shipyard >= 1
}
