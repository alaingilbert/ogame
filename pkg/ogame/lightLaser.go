package ogame

type lightLaser struct {
	BaseDefense
}

func newLightLaser() *lightLaser {
	d := new(lightLaser)
	d.Name = "light laser"
	d.ID = LightLaserID
	d.Price = Resources{Metal: 1500, Crystal: 500}
	d.StructuralIntegrity = 2000
	d.ShieldPower = 25
	d.WeaponPower = 100
	d.RapidfireFrom = map[ID]int64{DestroyerID: 10, BomberID: 20, DeathstarID: 200}
	d.Requirements = map[ID]int64{ShipyardID: 2, LaserTechnologyID: 3}
	return d
}
