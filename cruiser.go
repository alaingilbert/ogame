package ogame

type cruiser struct {
	BaseShip
}

func newCruiser() *cruiser {
	s := new(cruiser)
	s.Name = "cruiser"
	s.ID = CruiserID
	s.StructuralIntegrity = 27000
	s.ShieldPower = 50
	s.WeaponPower = 400
	s.BaseCargoCapacity = 800
	s.BaseSpeed = 15000
	s.FuelConsumption = 300
	s.FuelCapacity = 800
	s.RapidfireFrom = map[ID]int64{BattlecruiserID: 4, DeathstarID: 33, PathfinderID: 3}
	s.RapidfireAgainst = map[ID]int64{EspionageProbeID: 5, SolarSatelliteID: 5,
		LightFighterID: 6, RocketLauncherID: 10, CrawlerID: 5}
	s.Price = Resources{Metal: 20000, Crystal: 7000, Deuterium: 2000}
	s.Requirements = map[ID]int64{ShipyardID: 5, ImpulseDriveID: 4, IonTechnologyID: 2}
	return s
}
