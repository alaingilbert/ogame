package ogame

// Battlecruiser ...
type battlecruiser struct {
	BaseShip
}

func newBattlecruiser() *battlecruiser {
	b := new(battlecruiser)
	b.ID = BattlecruiserID
	b.StructuralIntegrity = 70000
	b.ShieldPower = 400
	b.WeaponPower = 700
	b.CargoCapacity = 750
	b.BaseSpeed = 1000
	b.FuelConsumption = 250
	b.RapidfireFrom = map[ID]int{DestroyerID: 2, DeathstarID: 15}
	b.RapidfireAgainst = map[ID]int{
		EspionageProbeID: 5, SolarSatelliteID: 5, SmallCargoID: 3, LargeCargoID: 3,
		HeavyFighterID: 4, CruiserID: 4, BattleshipID: 7,
	}
	b.Price = Resources{Metal: 30000, Crystal: 40000, Deuterium: 15000}
	b.Requirements = map[ID]int{LaserTechnologyID: 12, HyperspaceTechnologyID: 5,
		HyperspaceDriveID: 5, ShipyardID: 8}
	return b
}
