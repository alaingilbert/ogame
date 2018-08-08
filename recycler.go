package ogame

// Recycler ...
type recycler struct {
	BaseShip
}

func newRecycler() *recycler {
	s := new(recycler)
	s.ID = RecyclerID
	s.StructuralIntegrity = 16000
	s.ShieldPower = 10
	s.WeaponPower = 1
	s.CargoCapacity = 20000
	s.BaseSpeed = 2000
	s.FuelConsumption = 300
	s.RapidfireFrom = map[ID]int{DeathstarID: 250}
	s.RapidfireAgainst = map[ID]int{EspionageProbeID: 5, SolarSatelliteID: 5}
	s.Price = Resources{Metal: 10000, Crystal: 6000, Deuterium: 2000}
	s.Requirements = map[ID]int{ShipyardID: 4, CombustionDriveID: 6, ShieldingTechnologyID: 2}
	return s
}
