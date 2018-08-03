package ogame

// LaserTechnology ...
type laserTechnology struct {
	BaseTechnology
}

// NewLaserTechnology ...
func NewLaserTechnology() *laserTechnology {
	b := new(laserTechnology)
	b.ID = LaserTechnologyID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 200, Crystal: 100}
	b.Requirements = map[ID]int{EnergyTechnologyID: 2}
	return b
}
