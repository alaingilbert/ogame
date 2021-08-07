package ogame

type lightFighter struct {
	BaseShip
}

func newLightFighter() *lightFighter {
	l := new(lightFighter)
	l.Name = "light fighter"
	l.ID = LightFighterID
	l.StructuralIntegrity = 4000
	l.ShieldPower = 10
	l.WeaponPower = 50
	l.BaseCargoCapacity = 50
	l.BaseSpeed = 12500
	l.FuelConsumption = 20
	l.FuelCapacity = 50
	l.RapidfireFrom = map[ID]int64{CruiserID: 6, DeathstarID: 200, PathfinderID: 3}
	l.RapidfireAgainst = map[ID]int64{EspionageProbeID: 5, SolarSatelliteID: 5, CrawlerID: 5}
	l.Price = Resources{Metal: 3000, Crystal: 1000}
	l.Requirements = map[ID]int64{ShipyardID: 1, CombustionDriveID: 1}
	return l
}
