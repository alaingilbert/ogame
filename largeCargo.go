package ogame

// LargeCargo ...
type largeCargo struct {
	BaseShip
}

// NewLargeCargo ...
func NewLargeCargo() *largeCargo {
	s := new(largeCargo)
	s.ID = LargeCargoID
	s.StructuralIntegrity = 12000
	s.ShieldPower = 25
	s.WeaponPower = 5
	s.CargoCapacity = 25000
	s.BaseSpeed = 7500
	s.FuelConsumption = 50
	s.RapidfireFrom = map[ID]int{BattlecruiserID: 3, DeathstarID: 250}
	s.RapidfireAgainst = map[ID]int{EspionageProbeID: 5, SolarSatelliteID: 5}
	s.Price = Resources{Metal: 6000, Crystal: 6000}
	s.Requirements = map[ID]int{ShipyardID: 4, CombustionDriveID: 6}
	return s
}
