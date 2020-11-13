package ogame

type deathstar struct {
	BaseShip
}

func newDeathstar() *deathstar {
	s := new(deathstar)
	s.Name = "deathstar"
	s.ID = DeathstarID
	s.StructuralIntegrity = 9000000
	s.ShieldPower = 50000
	s.WeaponPower = 200000
	s.BaseCargoCapacity = 1000000
	s.BaseSpeed = 100
	s.FuelConsumption = 1
	s.FuelCapacity = 1000000
	s.RapidfireFrom = map[ID]int64{}
	s.RapidfireAgainst = map[ID]int64{SmallCargoID: 250, LargeCargoID: 250, LightFighterID: 200,
		HeavyFighterID: 100, CruiserID: 33, BattleshipID: 30, ColonyShipID: 250,
		RecyclerID: 250, EspionageProbeID: 1250, SolarSatelliteID: 1250, BomberID: 25,
		DestroyerID: 5, RocketLauncherID: 200, LightLaserID: 200, HeavyLaserID: 100,
		GaussCannonID: 50, IonCannonID: 100, BattlecruiserID: 15, CrawlerID: 1250, PathfinderID: 30, ReaperID: 10}
	s.Price = Resources{Metal: 5000000, Crystal: 4000000, Deuterium: 1000000}
	s.Requirements = map[ID]int64{ShipyardID: 12, GravitonTechnologyID: 1, HyperspaceDriveID: 7,
		HyperspaceTechnology.ID: 6}
	return s
}
