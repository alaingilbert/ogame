package ogame

// PlasmaTurret ...
type plasmaTurret struct {
	BaseDefense
}

func newPlasmaTurret() *plasmaTurret {
	d := new(plasmaTurret)
	d.Name = "plasma turret"
	d.ID = PlasmaTurretID
	d.Price = Resources{Metal: 50000, Crystal: 50000, Deuterium: 30000}
	d.StructuralIntegrity = 100000
	d.ShieldPower = 300
	d.WeaponPower = 3000
	d.Requirements = map[ID]int{ShipyardID: 8, PlasmaTechnologyID: 7}
	return d
}
