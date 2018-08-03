package ogame

// HeavyLaser ...
type heavyLaser struct {
	BaseDefense
}

// NewHeavyLaser ...
func NewHeavyLaser() *heavyLaser {
	d := new(heavyLaser)
	d.ID = HeavyLaserID
	d.Price = Resources{Metal: 6000, Crystal: 2000}
	d.StructuralIntegrity = 8000
	d.ShieldPower = 100
	d.WeaponPower = 250
	d.RapidfireFrom = map[ID]int{BomberID: 10, DeathstarID: 100}
	d.Requirements = map[ID]int{ShipyardID: 4, EnergyTechnologyID: 3, LaserTechnologyID: 6}
	return d
}
