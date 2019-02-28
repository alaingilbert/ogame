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
	s.RapidfireFrom = map[ID]int{BattlecruiserID: 4, DeathstarID: 33}
	s.RapidfireAgainst = map[ID]int{EspionageProbeID: 5, SolarSatelliteID: 5,
		LightFighterID: 6, RocketLauncherID: 10}
	s.Price = Resources{Metal: 20000, Crystal: 7000, Deuterium: 2000}
	s.Requirements = map[ID]int{ShipyardID: 5, ImpulseDriveID: 4, IonTechnologyID: 2}
	return s
}
