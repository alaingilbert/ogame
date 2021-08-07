package ogame

type reaper struct {
	BaseShip
}

func newReaper() *reaper {
	r := new(reaper)
	r.Name = "reaper"
	r.ID = ReaperID
	r.StructuralIntegrity = 140000
	r.ShieldPower = 700
	r.WeaponPower = 2800
	r.BaseCargoCapacity = 10000
	r.BaseSpeed = 7000
	r.FuelConsumption = 1100
	r.FuelCapacity = 10000
	r.RapidfireFrom = map[ID]int64{DeathstarID: 10, IonCannonID: 2}
	r.RapidfireAgainst = map[ID]int64{EspionageProbeID: 5, SolarSatelliteID: 5, CrawlerID: 5, BattleshipID: 7, BattlecruiserID: 7, BomberID: 4, DestroyerID: 3}
	r.Price = Resources{Metal: 85000, Crystal: 55000, Deuterium: 20000}
	r.Requirements = map[ID]int64{ShipyardID: 10, HyperspaceTechnologyID: 6, HyperspaceDriveID: 7, ShieldingTechnologyID: 6}
	return r
}
