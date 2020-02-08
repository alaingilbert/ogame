package ogame

import "math"

type solarSatellite struct {
	BaseShip
}

func newSolarSatellite() *solarSatellite {
	s := new(solarSatellite)
	s.Name = "solar satellite"
	s.ID = SolarSatelliteID
	s.StructuralIntegrity = 2000
	s.ShieldPower = 1
	s.WeaponPower = 1
	s.BaseCargoCapacity = 0
	s.BaseSpeed = 0
	s.FuelConsumption = 0
	s.RapidfireFrom = map[ID]int64{DeathstarID: 1250}
	s.Price = Resources{Crystal: 2000, Deuterium: 500}
	s.Requirements = map[ID]int64{ShipyardID: 1}
	return s
}

// Production gets the energy production of nbr solar satellite
func (s *solarSatellite) Production(temp Temperature, nbr int64, isCollector bool) int64 {
	energy := int64(float64(temp.Mean()+160)/6) * nbr
	if isCollector {
		energy += int64(math.Round(0.1 * float64(energy)))
	}
	return energy
}

// GetLevel only useful so the solar satellite can implement Building interface
func (s *solarSatellite) GetLevel(LazyResourcesBuildings, LazyFacilities, LazyResearches) int64 {
	return 0
}
