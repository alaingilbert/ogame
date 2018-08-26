package ogame

// HeavyFighter ...
type heavyFighter struct {
	BaseShip
}

func newHeavyFighter() *heavyFighter {
	s := new(heavyFighter)
	s.Name = "heavy fighter"
	s.ID = HeavyFighterID
	s.StructuralIntegrity = 10000
	s.ShieldPower = 25
	s.WeaponPower = 150
	s.CargoCapacity = 100
	s.BaseSpeed = 10000
	s.FuelConsumption = 75
	s.RapidfireFrom = map[ID]int{BattlecruiserID: 4, DeathstarID: 100}
	s.RapidfireAgainst = map[ID]int{EspionageProbeID: 5, SolarSatelliteID: 5, SmallCargoID: 3}
	s.Price = Resources{Metal: 6000, Crystal: 4000}
	s.Requirements = map[ID]int{ShipyardID: 3, ImpulseDriveID: 2, ArmourTechnologyID: 2}
	return s
}
