package ogame

type ionCannon struct {
	BaseDefense
}

func newIonCannon() *ionCannon {
	d := new(ionCannon)
	d.Name = "ion cannon"
	d.ID = IonCannonID
	d.Price = Resources{Metal: 5000, Crystal: 3000}
	d.StructuralIntegrity = 8000
	d.ShieldPower = 500
	d.WeaponPower = 150
	d.RapidfireFrom = map[ID]int64{BomberID: 10, DeathstarID: 100}
	d.RapidfireAgainst = map[ID]int64{ReaperID: 2}
	d.Requirements = map[ID]int64{ShipyardID: 4, IonTechnologyID: 4}
	return d
}
