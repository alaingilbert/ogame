package ogame

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
	s.CargoCapacity = 0
	s.BaseSpeed = 0
	s.FuelConsumption = 0
	s.RapidfireFrom = map[ID]int{DeathstarID: 1250}
	s.Price = Resources{Crystal: 2000, Deuterium: 500}
	s.Requirements = map[ID]int{ShipyardID: 1}
	return s
}

// Production gets the energy production of nbr solar satellite
func (s *solarSatellite) Production(temperatureMax, nbr int) int {
	return int((float64(temperatureMax)+140)/6) * nbr
}

// GetLevel only useful so the solar satellite can implement Building interface
func (s *solarSatellite) GetLevel(ResourcesBuildings, Facilities, Researches) int {
	return 0
}
