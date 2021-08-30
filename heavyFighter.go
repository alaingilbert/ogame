package ogame

type heavyFighter struct {
	BaseShip
}

func newHeavyFighter() *heavyFighter {
	s := new(heavyFighter)
	s.Name = "heavy fighter"
	s.ID = HeavyFighterID
	s.StructuralIntegrity = 10000
	s.ShieldPower = 25
	s.WeaponPower = 150
	s.BaseCargoCapacity = 100
	s.BaseSpeed = 10000
	s.FuelConsumption = 75
	s.FuelCapacity = 100
	s.RapidfireFrom = map[ID]int64{BattlecruiserID: 4, DeathstarID: 100, PathfinderID: 2}
	s.RapidfireAgainst = map[ID]int64{EspionageProbeID: 5, SolarSatelliteID: 5, SmallCargoID: 3, CrawlerID: 5}
	s.Price = Resources{Metal: 6000, Crystal: 4000}
	s.Requirements = map[ID]int64{ShipyardID: 3, ImpulseDriveID: 2, ArmourTechnologyID: 2}
	return s
}
