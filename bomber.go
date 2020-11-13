package ogame

type bomber struct {
	BaseShip
}

func newBomber() *bomber {
	b := new(bomber)
	b.Name = "bomber"
	b.ID = BomberID
	b.StructuralIntegrity = 75000
	b.ShieldPower = 500
	b.WeaponPower = 1000
	b.BaseCargoCapacity = 500
	b.BaseSpeed = 4000
	b.FuelConsumption = 700
	b.FuelCapacity = 500
	b.RapidfireFrom = map[ID]int64{DeathstarID: 25, ReaperID: 4}
	b.RapidfireAgainst = map[ID]int64{EspionageProbeID: 5, SolarSatelliteID: 5,
		RocketLauncherID: 20, LightLaserID: 20, HeavyLaserID: 10, IonCannonID: 10, CrawlerID: 5}
	b.Price = Resources{Metal: 50000, Crystal: 25000, Deuterium: 15000}
	b.Requirements = map[ID]int64{ImpulseDriveID: 6, ShipyardID: 8, PlasmaTechnologyID: 5}
	return b
}
