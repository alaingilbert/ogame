package ogame

// GaussCannon ...
type gaussCannon struct {
	BaseDefense
}

func newGaussCannon() *gaussCannon {
	d := new(gaussCannon)
	d.Name = "gauss cannon"
	d.ID = GaussCannonID
	d.Price = Resources{Metal: 20000, Crystal: 15000, Deuterium: 2000}
	d.StructuralIntegrity = 35000
	d.ShieldPower = 200
	d.WeaponPower = 1100
	d.RapidfireFrom = map[ID]int{DeathstarID: 50}
	d.Requirements = map[ID]int{ShipyardID: 6, WeaponsTechnologyID: 3, EnergyTechnologyID: 6, ShieldingTechnologyID: 1}
	return d
}
