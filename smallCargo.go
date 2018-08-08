package ogame

// SmallCargo ...
type smallCargo struct {
	BaseShip
}

func newSmallCargo() *smallCargo {
	s := new(smallCargo)
	s.ID = SmallCargoID
	s.StructuralIntegrity = 4000
	s.ShieldPower = 10
	s.WeaponPower = 5
	s.CargoCapacity = 5000
	s.BaseSpeed = 5000
	s.FuelConsumption = 10
	s.RapidfireFrom = map[ID]int{BattlecruiserID: 3, HeavyFighterID: 3, DeathstarID: 250}
	s.RapidfireAgainst = map[ID]int{EspionageProbeID: 5, SolarSatelliteID: 5}
	s.Price = Resources{Metal: 2000, Crystal: 2000}
	s.Requirements = map[ID]int{ShipyardID: 2, CombustionDriveID: 2}
	return s
}
