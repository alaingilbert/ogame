package ogame

// EspionageProbe ...
type espionageProbe struct {
	BaseShip
}

func newEspionageProbe() *espionageProbe {
	s := new(espionageProbe)
	s.ID = EspionageProbeID
	s.StructuralIntegrity = 1000
	s.ShieldPower = 0 //0.01
	s.WeaponPower = 0 //0.01
	s.CargoCapacity = 5
	s.BaseSpeed = 100000000
	s.FuelConsumption = 1
	s.RapidfireFrom = map[ID]int{BattlecruiserID: 5, DestroyerID: 5, BomberID: 5,
		RecyclerID: 5, ColonyShipID: 5, BattleshipID: 5, CruiserID: 5,
		HeavyFighterID: 5, LightFighterID: 5, LargeCargoID: 5, DeathstarID: 1250,
		SmallCargoID: 5}
	s.Price = Resources{Crystal: 1000}
	s.Requirements = map[ID]int{ShipyardID: 3, CombustionDriveID: 3, EspionageTechnologyID: 2}
	return s
}
