package ogame

type interplanetaryMissiles struct {
	BaseDefense
}

func newInterplanetaryMissiles() *interplanetaryMissiles {
	d := new(interplanetaryMissiles)
	d.Name = "interplanetary missiles"
	d.ID = InterplanetaryMissilesID
	d.Price = Resources{Metal: 12500, Crystal: 2500, Deuterium: 10000}
	d.StructuralIntegrity = 15000
	d.ShieldPower = 1
	d.WeaponPower = 12000
	d.Requirements = map[ID]int64{MissileSiloID: 4, ImpulseDriveID: 1}
	return d
}
