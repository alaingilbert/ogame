package ogame

// Destroyer ...
type destroyer struct {
	BaseShip
}

func newDestroyer() *destroyer {
	s := new(destroyer)
	s.ID = DestroyerID
	s.StructuralIntegrity = 110000
	s.ShieldPower = 500
	s.WeaponPower = 2000
	s.CargoCapacity = 2000
	s.BaseSpeed = 5000
	s.FuelConsumption = 1000
	s.RapidfireFrom = map[ID]int{DeathstarID: 5}
	s.RapidfireAgainst = map[ID]int{EspionageProbeID: 5, SolarSatelliteID: 5,
		LightLaserID: 10, BattlecruiserID: 2}
	s.Price = Resources{Metal: 60000, Crystal: 50000, Deuterium: 15000}
	s.Requirements = map[ID]int{ShipyardID: 9, HyperspaceDriveID: 6, HyperspaceTechnologyID: 5}
	return s
}
