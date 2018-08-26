package ogame

// IonCannon ...
type ionCannon struct {
	BaseDefense
}

func newIonCannon() *ionCannon {
	d := new(ionCannon)
	d.Name = "ion cannon"
	d.ID = IonCannonID
	d.Price = Resources{Metal: 2000, Crystal: 6000}
	d.StructuralIntegrity = 8000
	d.ShieldPower = 500
	d.WeaponPower = 150
	d.RapidfireFrom = map[ID]int{BomberID: 10, DeathstarID: 100}
	d.Requirements = map[ID]int{ShipyardID: 4, IonTechnologyID: 4}
	return d
}
