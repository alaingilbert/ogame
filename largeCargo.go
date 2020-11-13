package ogame

type largeCargo struct {
	BaseShip
}

func newLargeCargo() *largeCargo {
	s := new(largeCargo)
	s.Name = "large cargo"
	s.ID = LargeCargoID
	s.StructuralIntegrity = 12000
	s.ShieldPower = 25
	s.WeaponPower = 5
	s.BaseCargoCapacity = 25000
	s.BaseSpeed = 7500
	s.FuelConsumption = 50
	s.FuelCapacity = 25000
	s.RapidfireFrom = map[ID]int64{BattlecruiserID: 3, DeathstarID: 250}
	s.RapidfireAgainst = map[ID]int64{EspionageProbeID: 5, SolarSatelliteID: 5, CrawlerID: 5}
	s.Price = Resources{Metal: 6000, Crystal: 6000}
	s.Requirements = map[ID]int64{ShipyardID: 4, CombustionDriveID: 6}
	return s
}
