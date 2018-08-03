package ogame

// SolarSatellite ...
type solarSatellite struct {
	BaseShip
}

// NewSolarSatellite ...
func NewSolarSatellite() *solarSatellite {
	s := new(solarSatellite)
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

// Production ...
func (s *solarSatellite) Production(temperatureMax, nbr int) int {
	return int((float64(temperatureMax)+140)/6) * nbr
}

// GetIncreaseFactor ...
func (s *solarSatellite) GetIncreaseFactor() float64 {
	return 0
}

// GetBaseCost ...
func (s *solarSatellite) GetBaseCost() Resources {
	return s.Price
}

// GetLevel ...
func (s *solarSatellite) GetLevel(ResourcesBuildings, Facilities, Researches) int {
	return 0
}
