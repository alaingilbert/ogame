package ogame

// Bomber ...
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
	b.CargoCapacity = 500
	b.BaseSpeed = 4000
	b.FuelConsumption = 1000
	b.RapidfireFrom = map[ID]int{DeathstarID: 25}
	b.RapidfireAgainst = map[ID]int{EspionageProbeID: 5, SolarSatelliteID: 5,
		RocketLauncherID: 20, LightLaserID: 20, HeavyLaserID: 10, IonCannonID: 10}
	b.Price = Resources{Metal: 50000, Crystal: 25000, Deuterium: 15000}
	b.Requirements = map[ID]int{ImpulseDriveID: 6, ShipyardID: 8, PlasmaTechnologyID: 5}
	return b
}
