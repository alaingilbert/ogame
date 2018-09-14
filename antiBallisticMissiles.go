package ogame

type antiBallisticMissiles struct {
	BaseDefense
}

func newAntiBallisticMissiles() *antiBallisticMissiles {
	d := new(antiBallisticMissiles)
	d.Name = "anti ballistic missiles"
	d.ID = AntiBallisticMissilesID
	d.Price = Resources{Metal: 8000, Crystal: 2000}
	d.StructuralIntegrity = 8000
	d.ShieldPower = 1
	d.WeaponPower = 1
	d.Requirements = map[ID]int{MissileSiloID: 2}
	return d
}
