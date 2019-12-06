package ogame

type laserTechnology struct {
	BaseTechnology
}

func newLaserTechnology() *laserTechnology {
	b := new(laserTechnology)
	b.Name = "laser technology"
	b.ID = LaserTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 200, Crystal: 100}
	b.Requirements = map[ID]int64{EnergyTechnologyID: 2}
	return b
}
