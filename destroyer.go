package ogame

type destroyer struct {
	BaseShip
}

func newDestroyer() *destroyer {
	s := new(destroyer)
	s.Name = "destroyer"
	s.ID = DestroyerID
	s.StructuralIntegrity = 110000
	s.ShieldPower = 500
	s.WeaponPower = 2000
	s.BaseCargoCapacity = 2000
	s.BaseSpeed = 5000
	s.FuelConsumption = 1000
	s.FuelCapacity = 2000
	s.RapidfireFrom = map[ID]int64{DeathstarID: 5, ReaperID: 3}
	s.RapidfireAgainst = map[ID]int64{EspionageProbeID: 5, SolarSatelliteID: 5,
		LightLaserID: 10, BattlecruiserID: 2, CrawlerID: 5}
	s.Price = Resources{Metal: 60000, Crystal: 50000, Deuterium: 15000}
	s.Requirements = map[ID]int64{ShipyardID: 9, HyperspaceDriveID: 6, HyperspaceTechnologyID: 5}
	return s
}
