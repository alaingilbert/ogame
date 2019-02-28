package ogame

type battleship struct {
	BaseShip
}

func newBattleship() *battleship {
	b := new(battleship)
	b.Name = "battleship"
	b.ID = BattleshipID
	b.StructuralIntegrity = 60000
	b.ShieldPower = 200
	b.WeaponPower = 1000
	b.BaseCargoCapacity = 1500
	b.BaseSpeed = 10000
	b.FuelConsumption = 500
	b.RapidfireFrom = map[ID]int{BattlecruiserID: 7, DeathstarID: 30}
	b.RapidfireAgainst = map[ID]int{EspionageProbeID: 5, SolarSatelliteID: 5}
	b.Price = Resources{Metal: 45000, Crystal: 15000}
	b.Requirements = map[ID]int{ShipyardID: 7, HyperspaceDriveID: 4}
	return b
}
