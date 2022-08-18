package ogame

type heavyLaser struct {
	BaseDefense
}

func newHeavyLaser() *heavyLaser {
	d := new(heavyLaser)
	d.Name = "heavy laser"
	d.ID = HeavyLaserID
	d.Price = Resources{Metal: 6000, Crystal: 2000}
	d.StructuralIntegrity = 8000
	d.ShieldPower = 100
	d.WeaponPower = 250
	d.RapidfireFrom = map[ID]int64{BomberID: 10, DeathstarID: 100}
	d.Requirements = map[ID]int64{ShipyardID: 4, EnergyTechnologyID: 3, LaserTechnologyID: 6}
	return d
}
