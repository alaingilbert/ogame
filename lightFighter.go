package ogame

type lightFighter struct {
	BaseShip
}

func newLightFighter() *lightFighter {
	l := new(lightFighter)
	l.Name = "light fighter"
	l.ID = LightFighterID
	l.StructuralIntegrity = 4000
	l.ShieldPower = 10
	l.WeaponPower = 50
	l.BaseCargoCapacity = 50
	l.BaseSpeed = 12500
	l.FuelConsumption = 20
	l.RapidfireFrom = map[ID]int{CruiserID: 6, DeathstarID: 200}
	l.RapidfireAgainst = map[ID]int{EspionageProbeID: 5, SolarSatelliteID: 5}
	l.Price = Resources{Metal: 3000, Crystal: 1000}
	l.Requirements = map[ID]int{ShipyardID: 1, CombustionDriveID: 1}
	return l
}
