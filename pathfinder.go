package ogame

type pathfinder struct {
	BaseShip
}

func newPathfinder() *pathfinder {
	p := new(pathfinder)
	p.Name = "pathfinder"
	p.ID = PathfinderID
	p.StructuralIntegrity = 23000
	p.ShieldPower = 100
	p.WeaponPower = 200
	p.BaseCargoCapacity = 10000
	p.BaseSpeed = 12000
	p.FuelConsumption = 300
	p.FuelCapacity = 10000
	p.RapidfireFrom = map[ID]int64{BattleshipID: 5, DeathstarID: 30}
	p.RapidfireAgainst = map[ID]int64{EspionageProbeID: 5, SolarSatelliteID: 5, CrawlerID: 5, CruiserID: 3, LightFighterID: 3, HeavyFighterID: 2}
	p.Price = Resources{Metal: 8000, Crystal: 15000, Deuterium: 8000}
	p.Requirements = map[ID]int64{ShipyardID: 5, HyperspaceDriveID: 2}
	return p
}
