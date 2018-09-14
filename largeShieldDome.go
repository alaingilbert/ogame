package ogame

type largeShieldDome struct {
	BaseDefense
}

func newLargeShieldDome() *largeShieldDome {
	d := new(largeShieldDome)
	d.Name = "large shield dome"
	d.ID = LargeShieldDomeID
	d.Price = Resources{Metal: 50000, Crystal: 50000}
	d.StructuralIntegrity = 100000
	d.ShieldPower = 10000
	d.WeaponPower = 1
	d.Requirements = map[ID]int{ShieldingTechnologyID: 6, ShipyardID: 6}
	return d
}
