package ogame

type recycler struct {
	BaseShip
}

func newRecycler() *recycler {
	s := new(recycler)
	s.Name = "recycler"
	s.ID = RecyclerID
	s.StructuralIntegrity = 16000
	s.ShieldPower = 10
	s.WeaponPower = 1
	s.BaseCargoCapacity = 20000
	s.BaseSpeed = 2000
	s.FuelConsumption = 300
	s.FuelCapacity = 20000
	s.RapidfireFrom = map[ID]int64{DeathstarID: 250}
	s.RapidfireAgainst = map[ID]int64{EspionageProbeID: 5, SolarSatelliteID: 5, CrawlerID: 5}
	s.Price = Resources{Metal: 10000, Crystal: 6000, Deuterium: 2000}
	s.Requirements = map[ID]int64{ShipyardID: 4, CombustionDriveID: 6, ShieldingTechnologyID: 2}
	return s
}
