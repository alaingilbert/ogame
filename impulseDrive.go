package ogame

// ImpulseDrive ...
type impulseDrive struct {
	BaseTechnology
}

func newImpulseDrive() *impulseDrive {
	b := new(impulseDrive)
	b.Name = "impulse drive"
	b.ID = ImpulseDriveID
	b.IncreaseFactor = 2.0
	b.BaseCost = Resources{Metal: 2000, Crystal: 4000, Deuterium: 600}
	b.Requirements = map[ID]int{ResearchLabID: 2, EnergyTechnologyID: 1}
	return b
}
