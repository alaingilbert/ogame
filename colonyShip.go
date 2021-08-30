package ogame

type colonyShip struct {
	BaseShip
}

func newColonyShip() *colonyShip {
	s := new(colonyShip)
	s.Name = "colony ship"
	s.ID = ColonyShipID
	s.StructuralIntegrity = 30000
	s.ShieldPower = 100
	s.WeaponPower = 50
	s.BaseCargoCapacity = 7500
	s.BaseSpeed = 2500
	s.FuelConsumption = 1000
	s.FuelCapacity = 7500
	s.RapidfireFrom = map[ID]int64{DeathstarID: 250}
	s.RapidfireAgainst = map[ID]int64{EspionageProbeID: 5, SolarSatelliteID: 5, CrawlerID: 5}
	s.Price = Resources{Metal: 10000, Crystal: 20000, Deuterium: 10000}
	s.Requirements = map[ID]int64{ShipyardID: 4, ImpulseDriveID: 3}
	return s
}
